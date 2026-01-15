"""Database Adapter Factory.

Provides factory function to create database adapters based on configuration.
Supports environment variables and config file for database selection.
"""

import json
import os
from pathlib import Path
from typing import Any, Dict, Optional

from .base import DatabaseAdapter
from .sqlite_adapter import SQLiteAdapter
from .mysql_adapter import MySQLAdapter


def get_default_db_path() -> str:
    """Get global DB path.

    Returns the default database path, prioritizing environment variable
    if set, otherwise using the global path (~/.do/memory.db).

    Returns:
        str: Path to the database file
    """
    # Environment variable takes priority
    if os.environ.get("DO_DB_PATH"):
        return os.environ["DO_DB_PATH"]

    # Global path (~/.do/memory.db)
    home = Path.home()
    global_dir = home / ".do"
    global_dir.mkdir(parents=True, exist_ok=True)
    return str(global_dir / "memory.db")


def get_db_config(config_path: Optional[str] = None) -> Dict[str, Any]:
    """Load database configuration.

    Configuration priority:
    1. Environment variables (DO_DB_*)
    2. Config file (.do/db_config.json or specified path)
    3. Default values (sqlite with global path)

    Args:
        config_path: Optional path to config file

    Returns:
        Dict with database configuration

    Example config files:
        SQLite:
        {
            "type": "sqlite",
            "path": "~/.do/memory.db"
        }

        MySQL:
        {
            "type": "mysql",
            "host": "localhost",
            "port": 3306,
            "database": "do_memory",
            "user": "team_user",
            "password": "secret"
        }
    """
    config: Dict[str, Any] = {
        "type": "sqlite",
        "path": get_default_db_path(),
    }

    # 1. Check environment variables first (highest priority)
    env_type = os.environ.get("DO_DB_TYPE", "").lower()
    if env_type:
        config["type"] = env_type

    # For sqlite, path is already set by get_default_db_path()
    # which handles DO_DB_PATH environment variable

    elif config["type"] == "mysql":
        config["host"] = os.environ.get("DO_DB_HOST", "localhost")
        config["port"] = int(os.environ.get("DO_DB_PORT", "3306"))
        config["database"] = os.environ.get("DO_DB_NAME", "do_memory")
        config["user"] = os.environ.get("DO_DB_USER", "root")
        config["password"] = os.environ.get("DO_DB_PASSWORD", "")
        config["user_name"] = os.environ.get("DO_USER_NAME", "")

    # 2. Check config file (if no env vars set the type)
    if not env_type:
        if config_path is None:
            config_path = str(Path.cwd() / ".do" / "db_config.json")

        config_file = Path(config_path)
        if config_file.exists():
            try:
                with open(config_file, "r", encoding="utf-8") as f:
                    file_config = json.load(f)
                    config.update(file_config)
            except (json.JSONDecodeError, IOError) as e:
                # Log warning but continue with defaults
                import sys
                print(f"Warning: Failed to load db config: {e}", file=sys.stderr)

    return config


def get_db_adapter(
    config_path: Optional[str] = None,
    auto_migrate: bool = True,
) -> DatabaseAdapter:
    """Get database adapter based on configuration.

    Factory function that creates the appropriate adapter based on
    environment variables or config file.

    Args:
        config_path: Optional path to config file
        auto_migrate: If True, run migrations after connecting (default: True)

    Returns:
        DatabaseAdapter instance (SQLiteAdapter or MySQLAdapter)

    Raises:
        ValueError: If unknown database type specified

    Usage:
        # Auto-detect from environment/config
        db = get_db_adapter()
        db.connect()

        # With context manager
        with get_db_adapter() as db:
            users = db.fetchall("SELECT * FROM sessions")

        # Force specific type via environment
        os.environ["DO_DB_TYPE"] = "mysql"
        os.environ["DO_DB_HOST"] = "team-db.example.com"
        db = get_db_adapter()
    """
    config = get_db_config(config_path)
    db_type = config.get("type", "sqlite").lower()

    adapter: DatabaseAdapter

    if db_type == "sqlite":
        adapter = SQLiteAdapter(path=config.get("path"))

    elif db_type == "mysql":
        adapter = MySQLAdapter(
            host=config.get("host", "localhost"),
            port=config.get("port", 3306),
            database=config.get("database", "do_memory"),
            user=config.get("user", "root"),
            password=config.get("password", ""),
        )

    else:
        raise ValueError(
            f"Unknown database type: {db_type}. "
            f"Supported types: sqlite, mysql"
        )

    # Connect and optionally run migrations
    adapter.connect()
    if auto_migrate:
        adapter.run_migrations()

    return adapter


def create_sqlite_adapter(
    path: Optional[str] = None,
    auto_migrate: bool = True,
) -> SQLiteAdapter:
    """Create SQLite adapter directly.

    Convenience function for explicit SQLite usage.

    Args:
        path: Path to database file
        auto_migrate: If True, run migrations after connecting

    Returns:
        SQLiteAdapter instance
    """
    adapter = SQLiteAdapter(path=path)
    adapter.connect()
    if auto_migrate:
        adapter.run_migrations()
    return adapter


def create_mysql_adapter(
    host: str = "localhost",
    port: int = 3306,
    database: str = "do_memory",
    user: str = "root",
    password: str = "",
    auto_migrate: bool = True,
) -> MySQLAdapter:
    """Create MySQL adapter directly.

    Convenience function for explicit MySQL usage.

    Args:
        host: MySQL server host
        port: MySQL server port
        database: Database name
        user: MySQL username
        password: MySQL password
        auto_migrate: If True, run migrations after connecting

    Returns:
        MySQLAdapter instance
    """
    adapter = MySQLAdapter(
        host=host,
        port=port,
        database=database,
        user=user,
        password=password,
    )
    adapter.connect()
    if auto_migrate:
        adapter.run_migrations()
    return adapter
