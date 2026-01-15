"""Database Adapter Package for Do Framework.

Provides abstract interface and implementations for SQLite and MySQL databases.
Supports local (SQLite) for individual work and remote (MySQL) for team collaboration.

Usage:
    from .claude.hooks.do.lib.db import get_db_adapter

    # Auto-detect from environment/config
    db = get_db_adapter()

    # Or specify explicitly
    from .claude.hooks.do.lib.db import SQLiteAdapter, MySQLAdapter
    db = SQLiteAdapter(path=".do/memory.db")
    db = MySQLAdapter(host="localhost", database="do_memory", user="user", password="pass")

Configuration:
    Environment variables:
        DO_DB_TYPE: sqlite | mysql (default: sqlite)
        DO_DB_PATH: SQLite database path (default: .do/memory.db)
        DO_DB_HOST: MySQL host
        DO_DB_PORT: MySQL port (default: 3306)
        DO_DB_NAME: MySQL database name
        DO_DB_USER: MySQL username
        DO_DB_PASSWORD: MySQL password

    Config file (.do/db_config.json):
        {
            "type": "sqlite",
            "path": ".do/memory.db"
        }
        or
        {
            "type": "mysql",
            "host": "localhost",
            "port": 3306,
            "database": "do_memory",
            "user": "team_user",
            "password": "secret"
        }
"""

from .base import DatabaseAdapter
from .sqlite_adapter import SQLiteAdapter
from .mysql_adapter import MySQLAdapter
from .factory import get_db_adapter, get_db_config

__all__ = [
    "DatabaseAdapter",
    "SQLiteAdapter",
    "MySQLAdapter",
    "get_db_adapter",
    "get_db_config",
]
