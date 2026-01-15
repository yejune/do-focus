#!/usr/bin/env python3
"""Memory Store for Do Framework.

Database-agnostic CRUD module for sessions, observations, summaries, and plans.
Supports SQLite (local) and MySQL (team) through the adapter pattern.

Observation types: decision, bugfix, feature, refactor, docs, delegation
Plan statuses: draft, approved, completed
Timestamps: Unix epoch (int)
"""

import time
from typing import Any, Dict, List, Optional

from .db import DatabaseAdapter, get_db_adapter


class MemoryStore:
    """Session, observation, summary, and plan CRUD operations.

    Provides high-level access to the Do memory system with:
    - Session lifecycle management (create, end, archive)
    - Observation tracking with full-text search
    - Session summaries with full-text search
    - Plan management with status tracking
    - Progressive context disclosure for restoration

    Supports both SQLite (local work) and MySQL (team collaboration)
    through the DatabaseAdapter abstraction.

    Usage:
        # Auto-detect database from config/environment
        store = MemoryStore()
        store.create_session("session-123", "/path/to/project")
        store.add_observation("session-123", "feature", "Added login page")
        context = store.get_context_summary("/path/to/project", level=2)
        store.close()

        # Explicit adapter
        from .db import SQLiteAdapter, MySQLAdapter
        store = MemoryStore(adapter=SQLiteAdapter(path=".do/memory.db"))
        store = MemoryStore(adapter=MySQLAdapter(host="team-db.example.com", ...))
    """

    def __init__(
        self,
        db_path: Optional[str] = None,
        adapter: Optional[DatabaseAdapter] = None,
        user_name: Optional[str] = None,
    ):
        """Initialize MemoryStore with database adapter.

        Args:
            db_path: Optional path to SQLite database file (for backwards compatibility).
                    Ignored if adapter is provided.
            adapter: Optional pre-configured DatabaseAdapter instance.
                    If not provided, auto-detects from config/environment.
            user_name: Optional user name for team collaboration tracking.
                      Can also be set via DO_USER_NAME environment variable.
        """
        if adapter is not None:
            self._db = adapter
        elif db_path is not None:
            # Backwards compatibility: explicit SQLite path
            from .db import SQLiteAdapter
            self._db = SQLiteAdapter(path=db_path)
            self._db.connect()
            self._db.run_migrations()
        else:
            # Auto-detect from config/environment
            self._db = get_db_adapter()

        import os
        self._user_name = user_name or os.environ.get("DO_USER_NAME", "")

    # =========================================================================
    # Session Management
    # =========================================================================

    def create_session(
        self, session_id: str, project_path: str, user_name: Optional[str] = None
    ) -> None:
        """Create a new session.

        Args:
            session_id: Unique session identifier (e.g., UUID)
            project_path: Absolute path to the project root
            user_name: Optional user name (defaults to instance user_name)
        """
        user = user_name or self._user_name
        ph = self._db.placeholder

        self._db.execute(
            f"""
            INSERT INTO sessions (id, project_path, user_name, started_at, archived)
            VALUES ({ph}, {ph}, {ph}, {ph}, 0)
            """,
            (session_id, project_path, user, int(time.time())),
        )

    def end_session(self, session_id: str) -> None:
        """End a session by setting ended_at timestamp.

        Args:
            session_id: Session identifier to end
        """
        ph = self._db.placeholder
        self._db.execute(
            f"""
            UPDATE sessions
            SET ended_at = {ph}
            WHERE id = {ph} AND ended_at IS NULL
            """,
            (int(time.time()), session_id),
        )

    def get_session(self, session_id: str) -> Optional[Dict[str, Any]]:
        """Get a session by ID.

        Args:
            session_id: Session identifier

        Returns:
            Session dict with id, project_path, user_name, started_at, ended_at, archived
            or None if not found
        """
        ph = self._db.placeholder
        row = self._db.fetchone(
            f"""
            SELECT id, project_path, user_name, started_at, ended_at, archived
            FROM sessions
            WHERE id = {ph}
            """,
            (session_id,),
        )

        if row is None:
            return None

        return {
            "id": row["id"],
            "project_path": row["project_path"],
            "user_name": row.get("user_name", ""),
            "started_at": row["started_at"],
            "ended_at": row["ended_at"],
            "archived": bool(row["archived"]),
        }

    def get_recent_sessions(
        self, limit: int = 10, user_name: Optional[str] = None
    ) -> List[Dict[str, Any]]:
        """Get recent sessions ordered by start time.

        Args:
            limit: Maximum number of sessions to return (default: 10)
            user_name: Optional filter by user name (for team collaboration)

        Returns:
            List of session dicts
        """
        ph = self._db.placeholder

        if user_name:
            rows = self._db.fetchall(
                f"""
                SELECT id, project_path, user_name, started_at, ended_at, archived
                FROM sessions
                WHERE archived = 0 AND user_name = {ph}
                ORDER BY started_at DESC
                LIMIT {ph}
                """,
                (user_name, limit),
            )
        else:
            rows = self._db.fetchall(
                f"""
                SELECT id, project_path, user_name, started_at, ended_at, archived
                FROM sessions
                WHERE archived = 0
                ORDER BY started_at DESC
                LIMIT {ph}
                """,
                (limit,),
            )

        return [
            {
                "id": row["id"],
                "project_path": row["project_path"],
                "user_name": row.get("user_name", ""),
                "started_at": row["started_at"],
                "ended_at": row["ended_at"],
                "archived": bool(row["archived"]),
            }
            for row in rows
        ]

    def get_team_sessions(
        self, project_path: str, limit: int = 20
    ) -> List[Dict[str, Any]]:
        """Get team sessions for a project (all users).

        Args:
            project_path: Project path to filter
            limit: Maximum number of sessions to return (default: 20)

        Returns:
            List of session dicts from all team members
        """
        ph = self._db.placeholder
        rows = self._db.fetchall(
            f"""
            SELECT id, project_path, user_name, started_at, ended_at, archived
            FROM sessions
            WHERE project_path = {ph} AND archived = 0
            ORDER BY started_at DESC
            LIMIT {ph}
            """,
            (project_path, limit),
        )

        return [
            {
                "id": row["id"],
                "project_path": row["project_path"],
                "user_name": row.get("user_name", ""),
                "started_at": row["started_at"],
                "ended_at": row["ended_at"],
                "archived": bool(row["archived"]),
            }
            for row in rows
        ]

    def archive_session(self, session_id: str) -> None:
        """Archive a session (soft delete).

        Args:
            session_id: Session identifier to archive
        """
        ph = self._db.placeholder
        now = int(time.time())

        # End session if not already ended, then archive
        self._db.execute(
            f"""
            UPDATE sessions
            SET ended_at = COALESCE(ended_at, {ph}), archived = 1
            WHERE id = {ph}
            """,
            (now, session_id),
        )

    # =========================================================================
    # Observations
    # =========================================================================

    def add_observation(
        self,
        session_id: str,
        obs_type: str,
        content: str,
        file_path: Optional[str] = None,
        agent_name: Optional[str] = None,
    ) -> int:
        """Add an observation to a session.

        Args:
            session_id: Session identifier
            obs_type: Type of observation (decision, bugfix, feature,
                     refactor, docs, delegation)
            content: Observation content/description
            file_path: Optional related file path
            agent_name: Optional agent name for delegation tracking.
                       For delegation type, this tracks which agent was invoked.
                       Stored in file_path column (reused for agent tracking).

        Returns:
            ID of the created observation
        """
        ph = self._db.placeholder
        # For delegation type, agent_name takes precedence over file_path
        path_or_agent = agent_name if obs_type == "delegation" else file_path
        self._db.execute(
            f"""
            INSERT INTO observations (session_id, type, content, file_path, created_at)
            VALUES ({ph}, {ph}, {ph}, {ph}, {ph})
            """,
            (session_id, obs_type, content, path_or_agent, int(time.time())),
        )
        return self._db.last_insert_id

    def get_observations(
        self,
        session_id: Optional[str] = None,
        obs_type: Optional[str] = None,
    ) -> List[Dict[str, Any]]:
        """Get observations with optional filters.

        Args:
            session_id: Optional filter by session ID
            obs_type: Optional filter by observation type

        Returns:
            List of observation dicts
        """
        ph = self._db.placeholder
        sql = """
            SELECT id, session_id, type, content, file_path, created_at
            FROM observations
            WHERE 1=1
        """
        params: List[Any] = []

        if session_id is not None:
            sql += f" AND session_id = {ph}"
            params.append(session_id)

        if obs_type is not None:
            sql += f" AND type = {ph}"
            params.append(obs_type)

        sql += " ORDER BY created_at DESC"

        return self._db.fetchall(sql, tuple(params) if params else None)

    def search_observations(self, query: str) -> List[Dict[str, Any]]:
        """Search observations using full-text search.

        Uses FTS5 for SQLite, FULLTEXT for MySQL.

        Args:
            query: Search query string

        Returns:
            List of matching observation dicts ranked by relevance
        """
        from .db import MySQLAdapter

        if isinstance(self._db, MySQLAdapter):
            return self._db.fulltext_search_observations(query)

        # SQLite FTS5 search
        safe_query = query.replace('"', '""')
        return self._db.fetchall(
            """
            SELECT o.id, o.session_id, o.type, o.content, o.file_path, o.created_at
            FROM observations_fts
            JOIN observations o ON observations_fts.rowid = o.id
            WHERE observations_fts MATCH ?
            ORDER BY rank
            LIMIT 50
            """,
            (f'"{safe_query}"',),
        )

    def get_team_observations(
        self, project_path: str, limit: int = 50
    ) -> List[Dict[str, Any]]:
        """Get team observations for a project (all users).

        Args:
            project_path: Project path to filter
            limit: Maximum number of observations to return (default: 50)

        Returns:
            List of observation dicts from all team members, including user info
        """
        ph = self._db.placeholder
        return self._db.fetchall(
            f"""
            SELECT o.id, o.session_id, o.type, o.content, o.file_path, o.created_at,
                   s.user_name
            FROM observations o
            JOIN sessions s ON o.session_id = s.id
            WHERE s.project_path = {ph} AND s.archived = 0
            ORDER BY o.created_at DESC
            LIMIT {ph}
            """,
            (project_path, limit),
        )

    # =========================================================================
    # Session Summaries
    # =========================================================================

    def add_summary(
        self,
        session_id: str,
        request: str,
        investigation: str,
        result: str,
    ) -> int:
        """Add a session summary.

        Args:
            session_id: Session identifier
            request: What was requested/asked
            investigation: What was investigated/done
            result: The outcome/result

        Returns:
            ID of the created summary
        """
        ph = self._db.placeholder
        self._db.execute(
            f"""
            INSERT INTO session_summaries (session_id, request, investigation, result, created_at)
            VALUES ({ph}, {ph}, {ph}, {ph}, {ph})
            """,
            (session_id, request, investigation, result, int(time.time())),
        )
        return self._db.last_insert_id

    def get_summaries(
        self,
        session_id: Optional[str] = None,
    ) -> List[Dict[str, Any]]:
        """Get session summaries with optional filter.

        Args:
            session_id: Optional filter by session ID

        Returns:
            List of summary dicts
        """
        ph = self._db.placeholder

        if session_id is not None:
            return self._db.fetchall(
                f"""
                SELECT id, session_id, request, investigation, result, created_at
                FROM session_summaries
                WHERE session_id = {ph}
                ORDER BY created_at DESC
                """,
                (session_id,),
            )

        return self._db.fetchall(
            """
            SELECT id, session_id, request, investigation, result, created_at
            FROM session_summaries
            ORDER BY created_at DESC
            """
        )

    def search_summaries(self, query: str) -> List[Dict[str, Any]]:
        """Search summaries using full-text search.

        Uses FTS5 for SQLite, FULLTEXT for MySQL.

        Args:
            query: Search query string

        Returns:
            List of matching summary dicts ranked by relevance
        """
        from .db import MySQLAdapter

        if isinstance(self._db, MySQLAdapter):
            return self._db.fulltext_search_summaries(query)

        # SQLite FTS5 search
        safe_query = query.replace('"', '""')
        return self._db.fetchall(
            """
            SELECT s.id, s.session_id, s.request, s.investigation, s.result, s.created_at
            FROM summaries_fts
            JOIN session_summaries s ON summaries_fts.rowid = s.id
            WHERE summaries_fts MATCH ?
            ORDER BY rank
            LIMIT 50
            """,
            (f'"{safe_query}"',),
        )

    def get_team_summaries(
        self, project_path: str, days: int = 7
    ) -> List[Dict[str, Any]]:
        """Get team summaries for a project within a time window.

        Args:
            project_path: Project path to filter
            days: Number of days to look back (default: 7)

        Returns:
            List of summary dicts from all team members, including user info
        """
        ph = self._db.placeholder
        cutoff = int(time.time()) - (days * 24 * 60 * 60)

        return self._db.fetchall(
            f"""
            SELECT ss.id, ss.session_id, ss.request, ss.investigation, ss.result,
                   ss.created_at, s.user_name
            FROM session_summaries ss
            JOIN sessions s ON ss.session_id = s.id
            WHERE s.project_path = {ph}
              AND s.archived = 0
              AND ss.created_at >= {ph}
            ORDER BY ss.created_at DESC
            """,
            (project_path, cutoff),
        )

    # =========================================================================
    # Plans
    # =========================================================================

    def create_plan(
        self,
        session_id: str,
        title: str,
        file_path: str,
        content: str,
    ) -> int:
        """Create a new plan.

        Args:
            session_id: Session identifier
            title: Plan title
            file_path: Path to plan file
            content: Plan content/markdown

        Returns:
            ID of the created plan
        """
        ph = self._db.placeholder
        now = int(time.time())

        self._db.execute(
            f"""
            INSERT INTO plans (session_id, title, file_path, content, status, created_at, updated_at)
            VALUES ({ph}, {ph}, {ph}, {ph}, 'draft', {ph}, {ph})
            """,
            (session_id, title, file_path, content, now, now),
        )
        return self._db.last_insert_id

    def update_plan_status(self, plan_id: int, status: str) -> None:
        """Update a plan's status.

        Args:
            plan_id: Plan identifier
            status: New status (draft, approved, completed)
        """
        if status not in ("draft", "approved", "completed"):
            raise ValueError(f"Invalid status: {status}")

        ph = self._db.placeholder
        self._db.execute(
            f"""
            UPDATE plans
            SET status = {ph}, updated_at = {ph}
            WHERE id = {ph}
            """,
            (status, int(time.time()), plan_id),
        )

    def get_plans(self, status: Optional[str] = None) -> List[Dict[str, Any]]:
        """Get plans with optional status filter.

        Args:
            status: Optional filter by status (draft, approved, completed)

        Returns:
            List of plan dicts
        """
        ph = self._db.placeholder

        if status is not None:
            return self._db.fetchall(
                f"""
                SELECT id, session_id, title, file_path, content, status, created_at, updated_at
                FROM plans
                WHERE status = {ph}
                ORDER BY updated_at DESC
                """,
                (status,),
            )

        return self._db.fetchall(
            """
            SELECT id, session_id, title, file_path, content, status, created_at, updated_at
            FROM plans
            ORDER BY updated_at DESC
            """
        )

    def get_plan(self, plan_id: int) -> Optional[Dict[str, Any]]:
        """Get a plan by ID.

        Args:
            plan_id: Plan identifier

        Returns:
            Plan dict or None if not found
        """
        ph = self._db.placeholder
        return self._db.fetchone(
            f"""
            SELECT id, session_id, title, file_path, content, status, created_at, updated_at
            FROM plans
            WHERE id = {ph}
            """,
            (plan_id,),
        )

    # =========================================================================
    # Context Restoration
    # =========================================================================

    def get_context_summary(
        self, project_path: str, user_name: Optional[str] = None, level: int = 1
    ) -> str:
        """Get progressive context summary for restoration.

        Progressive disclosure levels:
        - Level 1: Recent 3 sessions summary (~50 tokens)
        - Level 2: + Related observations (~200 tokens)
        - Level 3: + Team context (~500 tokens)

        Args:
            project_path: Project path to filter sessions
            user_name: Optional user name to filter sessions (for personal context)
            level: Disclosure level (1, 2, or 3)

        Returns:
            Formatted context summary string
        """
        ph = self._db.placeholder
        limit = 3 if level < 3 else 10

        # Get recent sessions for this project
        if user_name:
            sessions = self._db.fetchall(
                f"""
                SELECT id, user_name, started_at, ended_at
                FROM sessions
                WHERE project_path = {ph} AND archived = 0 AND user_name = {ph}
                ORDER BY started_at DESC
                LIMIT {ph}
                """,
                (project_path, user_name, limit),
            )
        else:
            sessions = self._db.fetchall(
                f"""
                SELECT id, user_name, started_at, ended_at
                FROM sessions
                WHERE project_path = {ph} AND archived = 0
                ORDER BY started_at DESC
                LIMIT {ph}
                """,
                (project_path, limit),
            )

        if not sessions:
            return "No previous sessions found for this project."

        lines: List[str] = []
        lines.append("## Recent Sessions")

        for session in sessions:
            session_id = session["id"]
            session_user = session.get("user_name", "")
            started = session["started_at"]
            ended = session["ended_at"]

            # Format timestamps
            started_str = time.strftime("%Y-%m-%d %H:%M", time.localtime(started))
            ended_str = (
                time.strftime("%H:%M", time.localtime(ended)) if ended else "ongoing"
            )

            user_info = f" by {session_user}" if session_user else ""
            lines.append(f"\n### Session: {session_id[:8]}...{user_info}")
            lines.append(f"Time: {started_str} - {ended_str}")

            # Level 1: Just session info with summary count
            if level == 1:
                summary_count = self._db.fetchone(
                    f"SELECT COUNT(*) as cnt FROM session_summaries WHERE session_id = {ph}",
                    (session_id,),
                )
                obs_count = self._db.fetchone(
                    f"SELECT COUNT(*) as cnt FROM observations WHERE session_id = {ph}",
                    (session_id,),
                )
                lines.append(
                    f"Summaries: {summary_count['cnt'] if summary_count else 0}, "
                    f"Observations: {obs_count['cnt'] if obs_count else 0}"
                )

            # Level 2+: Include observations
            if level >= 2:
                obs_limit = 5 if level == 2 else 20
                observations = self._db.fetchall(
                    f"""
                    SELECT type, content, file_path
                    FROM observations
                    WHERE session_id = {ph}
                    ORDER BY created_at DESC
                    LIMIT {ph}
                    """,
                    (session_id, obs_limit),
                )

                if observations:
                    lines.append("\nObservations:")
                    for obs in observations:
                        file_info = f" ({obs['file_path']})" if obs["file_path"] else ""
                        content = obs["content"]
                        if level == 2 and len(content) > 100:
                            content = content[:100] + "..."
                        lines.append(f"- [{obs['type']}] {content}{file_info}")

            # Level 3: Include summaries
            if level >= 3:
                summaries = self._db.fetchall(
                    f"""
                    SELECT request, investigation, result
                    FROM session_summaries
                    WHERE session_id = {ph}
                    ORDER BY created_at DESC
                    LIMIT 5
                    """,
                    (session_id,),
                )

                if summaries:
                    lines.append("\nSummaries:")
                    for summary in summaries:
                        lines.append(f"- Request: {summary['request']}")
                        lines.append(f"  Investigation: {summary['investigation']}")
                        lines.append(f"  Result: {summary['result']}")

        # Level 3: Include team context if user_name is provided
        if level >= 3 and user_name:
            team_context = self.get_team_context(project_path, exclude_user=user_name)
            if team_context:
                lines.append("\n" + team_context)

        return "\n".join(lines)

    def get_team_context(
        self, project_path: str, exclude_user: Optional[str] = None
    ) -> str:
        """Get context from other team members' recent work.

        Provides visibility into what teammates have been working on,
        useful for coordination and avoiding conflicts.

        Args:
            project_path: Project path to filter
            exclude_user: User name to exclude (typically the current user)

        Returns:
            Formatted team context string
        """
        ph = self._db.placeholder
        cutoff = int(time.time()) - (7 * 24 * 60 * 60)  # Last 7 days

        # Get recent team observations
        if exclude_user:
            team_obs = self._db.fetchall(
                f"""
                SELECT o.type, o.content, o.file_path, o.created_at, s.user_name
                FROM observations o
                JOIN sessions s ON o.session_id = s.id
                WHERE s.project_path = {ph}
                  AND s.archived = 0
                  AND s.user_name IS NOT NULL
                  AND s.user_name != ''
                  AND s.user_name != {ph}
                  AND o.created_at >= {ph}
                ORDER BY o.created_at DESC
                LIMIT 20
                """,
                (project_path, exclude_user, cutoff),
            )
        else:
            team_obs = self._db.fetchall(
                f"""
                SELECT o.type, o.content, o.file_path, o.created_at, s.user_name
                FROM observations o
                JOIN sessions s ON o.session_id = s.id
                WHERE s.project_path = {ph}
                  AND s.archived = 0
                  AND s.user_name IS NOT NULL
                  AND s.user_name != ''
                  AND o.created_at >= {ph}
                ORDER BY o.created_at DESC
                LIMIT 20
                """,
                (project_path, cutoff),
            )

        if not team_obs:
            return ""

        lines: List[str] = []
        lines.append("## Team Activity (Last 7 Days)")

        # Group by user for cleaner output
        by_user: Dict[str, List[Dict[str, Any]]] = {}
        for obs in team_obs:
            user = obs.get("user_name", "unknown")
            if user not in by_user:
                by_user[user] = []
            by_user[user].append(obs)

        for user, observations in by_user.items():
            lines.append(f"\n### {user}")
            for obs in observations[:5]:  # Max 5 per user
                file_info = f" ({obs['file_path']})" if obs.get("file_path") else ""
                content = obs["content"]
                if len(content) > 80:
                    content = content[:80] + "..."
                lines.append(f"- [{obs['type']}] {content}{file_info}")

        return "\n".join(lines)

    # =========================================================================
    # Team Collaboration
    # =========================================================================

    def get_team_activity(
        self, project_path: str, days: int = 7
    ) -> List[Dict[str, Any]]:
        """Get team activity for a project (MySQL team mode).

        Args:
            project_path: Project path to filter
            days: Number of days to look back (default: 7)

        Returns:
            List of activity dicts with user_name, session_count, observation_count
        """
        ph = self._db.placeholder
        cutoff = int(time.time()) - (days * 24 * 60 * 60)

        return self._db.fetchall(
            f"""
            SELECT
                s.user_name,
                COUNT(DISTINCT s.id) as session_count,
                COUNT(o.id) as observation_count
            FROM sessions s
            LEFT JOIN observations o ON s.id = o.session_id
            WHERE s.project_path = {ph}
              AND s.started_at >= {ph}
              AND s.user_name IS NOT NULL
              AND s.user_name != ''
            GROUP BY s.user_name
            ORDER BY session_count DESC
            """,
            (project_path, cutoff),
        )

    # =========================================================================
    # Cleanup
    # =========================================================================

    def close(self) -> None:
        """Close the database connection."""
        self._db.close()

    def __enter__(self) -> "MemoryStore":
        """Context manager entry."""
        return self

    def __exit__(self, exc_type: Any, exc_val: Any, exc_tb: Any) -> None:
        """Context manager exit."""
        self.close()


# Convenience functions for quick access
def get_memory_store(
    db_path: Optional[str] = None,
    adapter: Optional[DatabaseAdapter] = None,
) -> MemoryStore:
    """Get a MemoryStore instance.

    Args:
        db_path: Optional custom SQLite database path (backwards compatible).
        adapter: Optional pre-configured DatabaseAdapter.

    Returns:
        Initialized MemoryStore.
    """
    return MemoryStore(db_path=db_path, adapter=adapter)
