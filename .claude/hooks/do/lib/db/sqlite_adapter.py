"""SQLite Database Adapter.

Implements DatabaseAdapter for SQLite with:
- WAL mode for better concurrent access
- Foreign key enforcement
- FTS5 full-text search support
- Automatic schema migrations
"""

import sqlite3
import time
from pathlib import Path
from typing import Any, Dict, List, Optional

from .base import DatabaseAdapter

# SQLite-specific migrations
SQLITE_MIGRATIONS = {
    1: """
-- Migration 001: Core tables

-- Session tracking
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    project_path TEXT,
    user_name TEXT,
    started_at INTEGER,
    ended_at INTEGER,
    archived INTEGER DEFAULT 0
);

-- Observations (work tracking)
CREATE TABLE IF NOT EXISTS observations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT,
    type TEXT,
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
    status TEXT DEFAULT 'draft',
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
CREATE INDEX IF NOT EXISTS idx_sessions_user ON sessions(user_name);
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


class SQLiteAdapter(DatabaseAdapter):
    """SQLite database adapter with WAL mode and FTS5 support.

    Features:
    - WAL mode for better concurrent access
    - Foreign key enforcement
    - FTS5 for full-text search
    - Row factory for dict-like access
    - Automatic migrations

    Usage:
        adapter = SQLiteAdapter(path=".do/memory.db")
        adapter.connect()
        adapter.run_migrations()

        # Or with context manager
        with SQLiteAdapter() as db:
            db.run_migrations()
            users = db.fetchall("SELECT * FROM sessions")
    """

    def __init__(self, path: Optional[str] = None):
        """Initialize SQLite adapter.

        Args:
            path: Path to database file. Defaults to .do/memory.db
                  relative to current working directory.
        """
        if path is None:
            path = str(Path.cwd() / ".do" / "memory.db")

        self.db_path = Path(path)
        self._connection: Optional[sqlite3.Connection] = None
        self._last_insert_id: int = 0

    def _ensure_directory(self) -> None:
        """Ensure the database directory exists."""
        self.db_path.parent.mkdir(parents=True, exist_ok=True)

    def connect(self) -> None:
        """Establish SQLite connection with optimized settings."""
        if self._connection is not None:
            return

        self._ensure_directory()
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

    def close(self) -> None:
        """Close the database connection."""
        if self._connection is not None:
            self._connection.close()
            self._connection = None

    def _get_connection(self) -> sqlite3.Connection:
        """Get connection, connecting if needed."""
        if self._connection is None:
            self.connect()
        return self._connection  # type: ignore

    def execute(self, sql: str, params: Optional[tuple] = None) -> sqlite3.Cursor:
        """Execute a SQL statement.

        Args:
            sql: SQL statement with ? placeholders
            params: Optional parameter tuple

        Returns:
            SQLite cursor
        """
        conn = self._get_connection()
        cursor = conn.cursor()

        if params:
            cursor.execute(sql, params)
        else:
            cursor.execute(sql)

        self._last_insert_id = cursor.lastrowid or 0
        conn.commit()
        return cursor

    def executemany(self, sql: str, params_list: List[tuple]) -> None:
        """Execute SQL with multiple parameter sets.

        Args:
            sql: SQL statement with ? placeholders
            params_list: List of parameter tuples
        """
        conn = self._get_connection()
        cursor = conn.cursor()
        cursor.executemany(sql, params_list)
        conn.commit()

    def fetchone(self, sql: str, params: Optional[tuple] = None) -> Optional[Dict[str, Any]]:
        """Execute query and fetch single row.

        Args:
            sql: SELECT statement
            params: Optional parameter tuple

        Returns:
            Dict with column names as keys, or None
        """
        conn = self._get_connection()
        cursor = conn.cursor()

        if params:
            cursor.execute(sql, params)
        else:
            cursor.execute(sql)

        row = cursor.fetchone()
        if row is None:
            return None

        return dict(row)

    def fetchall(self, sql: str, params: Optional[tuple] = None) -> List[Dict[str, Any]]:
        """Execute query and fetch all rows.

        Args:
            sql: SELECT statement
            params: Optional parameter tuple

        Returns:
            List of dicts with column names as keys
        """
        conn = self._get_connection()
        cursor = conn.cursor()

        if params:
            cursor.execute(sql, params)
        else:
            cursor.execute(sql)

        return [dict(row) for row in cursor.fetchall()]

    def get_current_version(self) -> int:
        """Get current migration version.

        Returns:
            Current version number (0 if no migrations)
        """
        conn = self._get_connection()
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
        return result[0] if result and result[0] is not None else 0

    def run_migrations(self) -> None:
        """Run all pending migrations."""
        current_version = self.get_current_version()
        conn = self._get_connection()

        for version in sorted(SQLITE_MIGRATIONS.keys()):
            if version > current_version:
                self._apply_migration(conn, version, SQLITE_MIGRATIONS[version])

    def _apply_migration(
        self, conn: sqlite3.Connection, version: int, sql: str
    ) -> None:
        """Apply a single migration.

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

    def begin_transaction(self) -> None:
        """Begin a database transaction."""
        conn = self._get_connection()
        conn.execute("BEGIN TRANSACTION")

    def commit(self) -> None:
        """Commit current transaction."""
        conn = self._get_connection()
        conn.commit()

    def rollback(self) -> None:
        """Rollback current transaction."""
        conn = self._get_connection()
        conn.rollback()

    @property
    def placeholder(self) -> str:
        """SQLite uses ? for parameter placeholders."""
        return "?"

    @property
    def last_insert_id(self) -> int:
        """Get the last inserted row ID."""
        return self._last_insert_id
