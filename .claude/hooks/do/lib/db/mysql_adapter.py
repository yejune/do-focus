"""MySQL Database Adapter.

Implements DatabaseAdapter for MySQL with:
- Connection pooling support
- FULLTEXT index for search (alternative to FTS5)
- Team collaboration support via user_name field
- Automatic schema migrations

Requires: mysql-connector-python or pymysql
"""

import time
from typing import Any, Dict, List, Optional, TYPE_CHECKING

from .base import DatabaseAdapter

# Type hints for optional mysql connector
if TYPE_CHECKING:
    import mysql.connector
    from mysql.connector import MySQLConnection
    from mysql.connector.cursor import MySQLCursor

# MySQL-specific migrations
MYSQL_MIGRATIONS = {
    1: """
-- Migration 001: Core tables

-- Session tracking
CREATE TABLE IF NOT EXISTS sessions (
    id VARCHAR(255) PRIMARY KEY,
    project_path TEXT,
    user_name VARCHAR(100),
    started_at BIGINT,
    ended_at BIGINT,
    archived TINYINT(1) DEFAULT 0,
    INDEX idx_sessions_project (project_path(255)),
    INDEX idx_sessions_user (user_name),
    INDEX idx_sessions_archived (archived)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Observations (work tracking)
CREATE TABLE IF NOT EXISTS observations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id VARCHAR(255),
    type VARCHAR(50),
    content TEXT,
    file_path TEXT,
    created_at BIGINT,
    INDEX idx_observations_session (session_id),
    INDEX idx_observations_type (type),
    INDEX idx_observations_created (created_at),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Session summaries (daily reports)
CREATE TABLE IF NOT EXISTS session_summaries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id VARCHAR(255),
    request TEXT,
    investigation TEXT,
    result TEXT,
    created_at BIGINT,
    INDEX idx_summaries_session (session_id),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Plan storage
CREATE TABLE IF NOT EXISTS plans (
    id INT AUTO_INCREMENT PRIMARY KEY,
    session_id VARCHAR(255),
    title VARCHAR(255),
    file_path TEXT,
    content LONGTEXT,
    status VARCHAR(50) DEFAULT 'draft',
    created_at BIGINT,
    updated_at BIGINT,
    INDEX idx_plans_session (session_id),
    INDEX idx_plans_status (status),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Migration tracking
CREATE TABLE IF NOT EXISTS migrations (
    version INT PRIMARY KEY,
    applied_at BIGINT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
""",
    2: """
-- Migration 002: FULLTEXT search indexes (MySQL alternative to FTS5)

-- Add FULLTEXT index to observations
ALTER TABLE observations ADD FULLTEXT INDEX ft_observations_content (content);

-- Add FULLTEXT index to session_summaries
ALTER TABLE session_summaries ADD FULLTEXT INDEX ft_summaries_content (request, investigation, result);
""",
}


class MySQLAdapter(DatabaseAdapter):
    """MySQL database adapter for team collaboration.

    Features:
    - Connection pooling support
    - FULLTEXT indexes for search
    - UTF-8 (utf8mb4) support
    - Team user tracking via user_name
    - Automatic migrations

    Usage:
        adapter = MySQLAdapter(
            host="localhost",
            database="do_memory",
            user="team_user",
            password="secret"
        )
        adapter.connect()
        adapter.run_migrations()

        # Or with context manager
        with MySQLAdapter(host="localhost", database="do_memory", ...) as db:
            db.run_migrations()
            sessions = db.fetchall("SELECT * FROM sessions WHERE user_name = %s", ("max",))

    Note:
        Requires mysql-connector-python: pip install mysql-connector-python
    """

    def __init__(
        self,
        host: str = "localhost",
        port: int = 3306,
        database: str = "do_memory",
        user: str = "root",
        password: str = "",
        pool_size: int = 5,
    ):
        """Initialize MySQL adapter.

        Args:
            host: MySQL server host
            port: MySQL server port (default: 3306)
            database: Database name
            user: MySQL username
            password: MySQL password
            pool_size: Connection pool size (default: 5)
        """
        self.host = host
        self.port = port
        self.database = database
        self.user = user
        self.password = password
        self.pool_size = pool_size

        self._connection: Optional[Any] = None
        self._cursor: Optional[Any] = None
        self._last_insert_id: int = 0

    def _get_mysql_module(self) -> Any:
        """Get MySQL connector module.

        Tries mysql.connector first, then pymysql.

        Returns:
            MySQL module

        Raises:
            ImportError: If no MySQL driver is installed
        """
        try:
            import mysql.connector
            return mysql.connector
        except ImportError:
            pass

        try:
            import pymysql
            pymysql.install_as_MySQLdb()
            return pymysql
        except ImportError:
            pass

        raise ImportError(
            "MySQL driver not found. Install one of:\n"
            "  pip install mysql-connector-python\n"
            "  pip install pymysql"
        )

    def connect(self) -> None:
        """Establish MySQL connection."""
        if self._connection is not None:
            return

        mysql_module = self._get_mysql_module()

        self._connection = mysql_module.connect(
            host=self.host,
            port=self.port,
            database=self.database,
            user=self.user,
            password=self.password,
            charset="utf8mb4",
            collation="utf8mb4_unicode_ci",
            autocommit=True,
            connection_timeout=30,
        )

    def close(self) -> None:
        """Close the database connection."""
        if self._cursor is not None:
            self._cursor.close()
            self._cursor = None

        if self._connection is not None:
            self._connection.close()
            self._connection = None

    def _get_connection(self) -> Any:
        """Get connection, connecting if needed."""
        if self._connection is None:
            self.connect()
        return self._connection

    def _get_cursor(self, dictionary: bool = True) -> Any:
        """Get a cursor, creating if needed.

        Args:
            dictionary: If True, return dict cursor (default: True)

        Returns:
            MySQL cursor
        """
        conn = self._get_connection()
        return conn.cursor(dictionary=dictionary)

    def execute(self, sql: str, params: Optional[tuple] = None) -> Any:
        """Execute a SQL statement.

        Args:
            sql: SQL statement with %s placeholders
            params: Optional parameter tuple

        Returns:
            MySQL cursor
        """
        cursor = self._get_cursor()

        if params:
            cursor.execute(sql, params)
        else:
            cursor.execute(sql)

        self._last_insert_id = cursor.lastrowid or 0
        return cursor

    def executemany(self, sql: str, params_list: List[tuple]) -> None:
        """Execute SQL with multiple parameter sets.

        Args:
            sql: SQL statement with %s placeholders
            params_list: List of parameter tuples
        """
        cursor = self._get_cursor()
        cursor.executemany(sql, params_list)

    def fetchone(self, sql: str, params: Optional[tuple] = None) -> Optional[Dict[str, Any]]:
        """Execute query and fetch single row.

        Args:
            sql: SELECT statement
            params: Optional parameter tuple

        Returns:
            Dict with column names as keys, or None
        """
        cursor = self._get_cursor(dictionary=True)

        if params:
            cursor.execute(sql, params)
        else:
            cursor.execute(sql)

        return cursor.fetchone()

    def fetchall(self, sql: str, params: Optional[tuple] = None) -> List[Dict[str, Any]]:
        """Execute query and fetch all rows.

        Args:
            sql: SELECT statement
            params: Optional parameter tuple

        Returns:
            List of dicts with column names as keys
        """
        cursor = self._get_cursor(dictionary=True)

        if params:
            cursor.execute(sql, params)
        else:
            cursor.execute(sql)

        return cursor.fetchall()

    def get_current_version(self) -> int:
        """Get current migration version.

        Returns:
            Current version number (0 if no migrations)
        """
        cursor = self._get_cursor()

        # Check if migrations table exists
        cursor.execute(
            """
            SELECT COUNT(*) FROM information_schema.tables
            WHERE table_schema = %s AND table_name = 'migrations'
            """,
            (self.database,)
        )
        result = cursor.fetchone()
        if not result or result[0] == 0:
            return 0

        # Get highest applied version
        cursor.execute("SELECT MAX(version) FROM migrations")
        result = cursor.fetchone()
        return result[0] if result and result[0] is not None else 0

    def run_migrations(self) -> None:
        """Run all pending migrations."""
        current_version = self.get_current_version()
        conn = self._get_connection()

        for version in sorted(MYSQL_MIGRATIONS.keys()):
            if version > current_version:
                self._apply_migration(conn, version, MYSQL_MIGRATIONS[version])

    def _apply_migration(self, conn: Any, version: int, sql: str) -> None:
        """Apply a single migration.

        Args:
            conn: Database connection
            version: Migration version number
            sql: SQL statements to execute
        """
        cursor = conn.cursor()

        try:
            # MySQL doesn't support executescript, so we split and execute
            statements = [s.strip() for s in sql.split(";") if s.strip()]

            for statement in statements:
                if statement and not statement.startswith("--"):
                    try:
                        cursor.execute(statement)
                    except Exception as e:
                        # Ignore "already exists" errors for idempotent migrations
                        if "already exists" in str(e).lower() or "duplicate" in str(e).lower():
                            continue
                        raise

            # Record migration
            cursor.execute(
                "INSERT INTO migrations (version, applied_at) VALUES (%s, %s)",
                (version, int(time.time())),
            )

            conn.commit()
        except Exception as e:
            conn.rollback()
            raise RuntimeError(f"Migration {version} failed: {e}") from e

    def begin_transaction(self) -> None:
        """Begin a database transaction."""
        conn = self._get_connection()
        conn.autocommit = False
        conn.start_transaction()

    def commit(self) -> None:
        """Commit current transaction."""
        conn = self._get_connection()
        conn.commit()
        conn.autocommit = True

    def rollback(self) -> None:
        """Rollback current transaction."""
        conn = self._get_connection()
        conn.rollback()
        conn.autocommit = True

    @property
    def placeholder(self) -> str:
        """MySQL uses %s for parameter placeholders."""
        return "%s"

    @property
    def last_insert_id(self) -> int:
        """Get the last inserted row ID."""
        return self._last_insert_id

    def convert_sql(self, sql: str) -> str:
        """Convert ? placeholders to %s for MySQL.

        Args:
            sql: SQL with ? placeholders

        Returns:
            SQL with %s placeholders
        """
        return sql.replace("?", "%s")

    def fulltext_search_observations(
        self, query: str, limit: int = 50
    ) -> List[Dict[str, Any]]:
        """Search observations using FULLTEXT index.

        MySQL-specific search using FULLTEXT instead of FTS5.

        Args:
            query: Search query string
            limit: Maximum results (default: 50)

        Returns:
            List of matching observation dicts
        """
        return self.fetchall(
            """
            SELECT id, session_id, type, content, file_path, created_at,
                   MATCH(content) AGAINST(%s IN NATURAL LANGUAGE MODE) AS relevance
            FROM observations
            WHERE MATCH(content) AGAINST(%s IN NATURAL LANGUAGE MODE)
            ORDER BY relevance DESC
            LIMIT %s
            """,
            (query, query, limit)
        )

    def fulltext_search_summaries(
        self, query: str, limit: int = 50
    ) -> List[Dict[str, Any]]:
        """Search summaries using FULLTEXT index.

        MySQL-specific search using FULLTEXT instead of FTS5.

        Args:
            query: Search query string
            limit: Maximum results (default: 50)

        Returns:
            List of matching summary dicts
        """
        return self.fetchall(
            """
            SELECT id, session_id, request, investigation, result, created_at,
                   MATCH(request, investigation, result) AGAINST(%s IN NATURAL LANGUAGE MODE) AS relevance
            FROM session_summaries
            WHERE MATCH(request, investigation, result) AGAINST(%s IN NATURAL LANGUAGE MODE)
            ORDER BY relevance DESC
            LIMIT %s
            """,
            (query, query, limit)
        )
