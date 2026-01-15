"""Abstract Database Adapter Interface.

Defines the contract for all database adapters (SQLite, MySQL, etc.).
All adapters must implement these methods to ensure compatibility.
"""

from abc import ABC, abstractmethod
from typing import Any, Dict, List, Optional


class DatabaseAdapter(ABC):
    """Abstract base class for database adapters.

    Provides a unified interface for different database backends.
    Implementations must handle:
    - Connection management
    - Query execution with parameterized queries
    - Transaction management
    - Schema migrations

    Usage:
        class MyAdapter(DatabaseAdapter):
            def connect(self) -> None:
                # Implementation
                ...

        with MyAdapter() as db:
            db.execute("INSERT INTO users (name) VALUES (?)", ("John",))
            users = db.fetchall("SELECT * FROM users")
    """

    @abstractmethod
    def connect(self) -> None:
        """Establish database connection.

        Should configure:
        - Connection pooling if applicable
        - Character encoding (UTF-8)
        - Timeout settings
        """
        ...

    @abstractmethod
    def close(self) -> None:
        """Close database connection and release resources."""
        ...

    @abstractmethod
    def execute(self, sql: str, params: Optional[tuple] = None) -> Any:
        """Execute a SQL statement.

        Args:
            sql: SQL statement with parameter placeholders (? for SQLite, %s for MySQL)
            params: Optional tuple of parameter values

        Returns:
            Cursor or execution result depending on implementation

        Raises:
            DatabaseError: If execution fails
        """
        ...

    @abstractmethod
    def executemany(self, sql: str, params_list: List[tuple]) -> None:
        """Execute a SQL statement with multiple parameter sets.

        Args:
            sql: SQL statement with parameter placeholders
            params_list: List of parameter tuples

        Raises:
            DatabaseError: If execution fails
        """
        ...

    @abstractmethod
    def fetchone(self, sql: str, params: Optional[tuple] = None) -> Optional[Dict[str, Any]]:
        """Execute query and fetch single row.

        Args:
            sql: SELECT statement
            params: Optional parameter tuple

        Returns:
            Dict with column names as keys, or None if no row found
        """
        ...

    @abstractmethod
    def fetchall(self, sql: str, params: Optional[tuple] = None) -> List[Dict[str, Any]]:
        """Execute query and fetch all rows.

        Args:
            sql: SELECT statement
            params: Optional parameter tuple

        Returns:
            List of dicts with column names as keys
        """
        ...

    @abstractmethod
    def run_migrations(self) -> None:
        """Run all pending database migrations.

        Should:
        - Check current schema version
        - Apply migrations in order
        - Record applied migrations in migrations table
        """
        ...

    @abstractmethod
    def get_current_version(self) -> int:
        """Get current migration version.

        Returns:
            Current version number (0 if no migrations applied)
        """
        ...

    @abstractmethod
    def begin_transaction(self) -> None:
        """Begin a database transaction.

        For explicit transaction control. Note that some operations
        may auto-commit; use transactions for atomic operations.
        """
        ...

    @abstractmethod
    def commit(self) -> None:
        """Commit current transaction."""
        ...

    @abstractmethod
    def rollback(self) -> None:
        """Rollback current transaction."""
        ...

    @property
    @abstractmethod
    def placeholder(self) -> str:
        """Get parameter placeholder for this database.

        Returns:
            '?' for SQLite, '%s' for MySQL
        """
        ...

    @property
    @abstractmethod
    def last_insert_id(self) -> int:
        """Get the last inserted row ID.

        Returns:
            ID of the last inserted row
        """
        ...

    def __enter__(self) -> "DatabaseAdapter":
        """Context manager entry - connect to database."""
        self.connect()
        return self

    def __exit__(self, exc_type: Any, exc_val: Any, exc_tb: Any) -> None:
        """Context manager exit - close database connection."""
        self.close()

    def convert_sql(self, sql: str) -> str:
        """Convert SQL placeholder syntax if needed.

        Converts '?' placeholders to the database-specific format.
        Override in subclasses that use different placeholder syntax.

        Args:
            sql: SQL with '?' placeholders

        Returns:
            SQL with database-specific placeholders
        """
        return sql
