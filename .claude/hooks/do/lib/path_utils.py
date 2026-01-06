"""Path utility functions for Do hooks"""

from pathlib import Path


def find_project_root() -> Path:
    """
    Find project root by locating .do directory.

    Starts from current file location and traverses upward
    until .do directory is found.

    Returns:
        Path: Project root directory containing .do/

    Fallback:
        Returns Path.cwd() if .do not found
    """
    current = Path(__file__).resolve().parent

    # Traverse upward to find .do directory
    while current != current.parent:
        if (current / ".do").is_dir():
            return current
        current = current.parent

    # Fallback to current working directory
    return Path.cwd()
