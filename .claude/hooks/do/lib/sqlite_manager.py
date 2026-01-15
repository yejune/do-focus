"""
SQLite Manager for Do Framework Memory System.

DEPRECATED: Use db.SQLiteAdapter instead.

This module is kept for backwards compatibility.
New code should use:
    from .db import SQLiteAdapter, get_db_adapter

Provides SQLite initialization, migration system, and connection management.
Database path: .do/memory.db (project root relative)

Features:
- WAL mode for better concurrent access
- Foreign key enforcement
- Automatic migration system
- FTS5 for full-text search
"""

import sqlite3
import time
import warnings
from pathlib import Path
from typing import Optional

# Deprecation warning for direct usage
warnings.warn(
    "sqlite_manager is deprecated. Use 'from .db import SQLiteAdapter' instead.",
    DeprecationWarning,
    stacklevel=2,
)

# Migration definitions
MIGRATIONS = {
    1: """
-- Migration 001: Core tables

-- Session tracking
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    project_path TEXT,
    started_at INTEGER,
    ended_at INTEGER,
    archived INTEGER DEFAULT 0
);

-- Observations (work tracking)
CREATE TABLE IF NOT EXISTS observations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT,
    type TEXT,  -- decision, bugfix, feature, refactor, docs, delegation
    content TEXT,
    file_path TEXT,
    created_at INTEGER,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);

-- Session summaries (daily reports)
CREATE TABLE IF NOT EXISTS session_summaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT,
    request TEXT,
    investigation TEXT,
    result TEXT,
    created_at INTEGER,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);

-- Plan storage
CREATE TABLE IF NOT EXISTS plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT,
    title TEXT,
    file_path TEXT,
    content TEXT,
    status TEXT DEFAULT 'draft',  -- draft, approved, completed
    created_at INTEGER,
    updated_at INTEGER,
    FOREIGN KEY (session_id) REFERENCES sessions(id)
);

-- Migration tracking
CREATE TABLE IF NOT EXISTS migrations (
    version INTEGER PRIMARY KEY,
    applied_at INTEGER
);

-- Indexes for common queries
CREATE INDEX IF NOT EXISTS idx_sessions_project ON sessions(project_path);
CREATE INDEX IF NOT EXISTS idx_sessions_archived ON sessions(archived);
CREATE INDEX IF NOT EXISTS idx_observations_session ON observations(session_id);
CREATE INDEX IF NOT EXISTS idx_observations_type ON observations(type);
CREATE INDEX IF NOT EXISTS idx_observations_created ON observations(created_at);
CREATE INDEX IF NOT EXISTS idx_summaries_session ON session_summaries(session_id);
CREATE INDEX IF NOT EXISTS idx_plans_session ON plans(session_id);
CREATE INDEX IF NOT EXISTS idx_plans_status ON plans(status);
""",
    2: """
-- Migration 002: FTS5 full-text search

-- FTS5 virtual table for observations
CREATE VIRTUAL TABLE IF NOT EXISTS observations_fts USING fts5(
    content,
    type,
    content='observations',
    content_rowid='id'
);

-- FTS5 virtual table for session summaries
CREATE VIRTUAL TABLE IF NOT EXISTS summaries_fts USING fts5(
    request,
    investigation,
    result,
    content='session_summaries',
    content_rowid='id'
);

-- Trigger: observations insert -> FTS update
CREATE TRIGGER IF NOT EXISTS observations_ai AFTER INSERT ON observations BEGIN
    INSERT INTO observations_fts(rowid, content, type)
    VALUES (new.id, new.content, new.type);
END;

-- Trigger: observations delete -> FTS update
CREATE TRIGGER IF NOT EXISTS observations_ad AFTER DELETE ON observations BEGIN
    INSERT INTO observations_fts(observations_fts, rowid, content, type)
    VALUES ('delete', old.id, old.content, old.type);
END;

-- Trigger: observations update -> FTS update
CREATE TRIGGER IF NOT EXISTS observations_au AFTER UPDATE ON observations BEGIN
    INSERT INTO observations_fts(observations_fts, rowid, content, type)
    VALUES ('delete', old.id, old.content, old.type);
    INSERT INTO observations_fts(rowid, content, type)
    VALUES (new.id, new.content, new.type);
END;

-- Trigger: session_summaries insert -> FTS update
CREATE TRIGGER IF NOT EXISTS summaries_ai AFTER INSERT ON session_summaries BEGIN
    INSERT INTO summaries_fts(rowid, request, investigation, result)
    VALUES (new.id, new.request, new.investigation, new.result);
END;

-- Trigger: session_summaries delete -> FTS update
CREATE TRIGGER IF NOT EXISTS summaries_ad AFTER DELETE ON session_summaries BEGIN
    INSERT INTO summaries_fts(summaries_fts, rowid, request, investigation, result)
    VALUES ('delete', old.id, old.request, old.investigation, old.result);
END;

-- Trigger: session_summaries update -> FTS update
CREATE TRIGGER IF NOT EXISTS summaries_au AFTER UPDATE ON session_summaries BEGIN
    INSERT INTO summaries_fts(summaries_fts, rowid, request, investigation, result)
    VALUES ('delete', old.id, old.request, old.investigation, old.result);
    INSERT INTO summaries_fts(rowid, request, investigation, result)
    VALUES (new.id, new.request, new.investigation, new.result);
END;
""",
}


class SQLiteManager:
    """
    SQLite database manager for Do Framework memory system.

    Handles:
    - Database initialization with WAL mode
    - Foreign key enforcement
    - Automatic schema migrations
    - Connection pooling (single connection for now)

    Usage:
        manager = SQLiteManager()  # Uses default .do/memory.db
        manager.run_migrations()
        conn = manager.get_connection()
        # ... use connection ...
        manager.close()
    """

    def __init__(self, db_path: Optional[str] = None):
        """
        Initialize SQLite manager.

        Args:
            db_path: Path to database file. Defaults to .do/memory.db
                    relative to current working directory.
        """
        if db_path is None:
            # Default to .do/memory.db in project root
            db_path = str(Path.cwd() / ".do" / "memory.db")

        self.db_path = Path(db_path)
        self._connection: Optional[sqlite3.Connection] = None
        self._ensure_directory()

    def _ensure_directory(self) -> None:
        """Ensure the database directory exists."""
        self.db_path.parent.mkdir(parents=True, exist_ok=True)

    def get_connection(self) -> sqlite3.Connection:
        """
        Get or create a database connection.

        Returns:
            sqlite3.Connection with WAL mode and foreign keys enabled.
        """
        if self._connection is None:
            self._connection = sqlite3.connect(
                str(self.db_path),
                check_same_thread=False,
                timeout=30.0,
            )
            # Enable WAL mode for better concurrent access
            self._connection.execute("PRAGMA journal_mode=WAL")
            # Enable foreign key constraints
            self._connection.execute("PRAGMA foreign_keys=ON")
            # Optimize for performance
            self._connection.execute("PRAGMA synchronous=NORMAL")
            self._connection.execute("PRAGMA cache_size=-64000")  # 64MB cache
            # Return rows as dictionaries
            self._connection.row_factory = sqlite3.Row

        return self._connection

    def get_current_version(self) -> int:
        """
        Get the current migration version.

        Returns:
            Current version number, or 0 if no migrations applied.
        """
        conn = self.get_connection()
        cursor = conn.cursor()

        # Check if migrations table exists
        cursor.execute(
            """
            SELECT name FROM sqlite_master
            WHERE type='table' AND name='migrations'
            """
        )
        if cursor.fetchone() is None:
            return 0

        # Get highest applied version
        cursor.execute("SELECT MAX(version) FROM migrations")
        result = cursor.fetchone()
        return result[0] if result[0] is not None else 0

    def run_migrations(self) -> None:
        """
        Run all pending migrations.

        Applies migrations in order, tracking each in the migrations table.
        """
        conn = self.get_connection()
        current_version = self.get_current_version()

        for version in sorted(MIGRATIONS.keys()):
            if version > current_version:
                self._apply_migration(conn, version, MIGRATIONS[version])

    def _apply_migration(
        self, conn: sqlite3.Connection, version: int, sql: str
    ) -> None:
        """
        Apply a single migration.

        Args:
            conn: Database connection
            version: Migration version number
            sql: SQL statements to execute
        """
        cursor = conn.cursor()

        try:
            # Execute migration SQL
            cursor.executescript(sql)

            # Record migration
            cursor.execute(
                "INSERT INTO migrations (version, applied_at) VALUES (?, ?)",
                (version, int(time.time())),
            )

            conn.commit()
        except sqlite3.Error as e:
            conn.rollback()
            raise RuntimeError(f"Migration {version} failed: {e}") from e

    def close(self) -> None:
        """Close the database connection."""
        if self._connection is not None:
            self._connection.close()
            self._connection = None

    def __enter__(self) -> "SQLiteManager":
        """Context manager entry."""
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        """Context manager exit."""
        self.close()


# Convenience function for quick access
def get_db(db_path: Optional[str] = None) -> SQLiteManager:
    """
    Get a SQLiteManager instance with migrations applied.

    Args:
        db_path: Optional custom database path.

    Returns:
        Initialized SQLiteManager with all migrations applied.
    """
    manager = SQLiteManager(db_path)
    manager.run_migrations()
    return manager
