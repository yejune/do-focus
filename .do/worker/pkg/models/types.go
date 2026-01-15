// Package models defines shared types for the Do Worker Service.
package models

import "time"

// Session represents a Claude session.
type Session struct {
	ID        string     `json:"id" db:"id"`
	UserName  string     `json:"user_name" db:"user_name"`
	ProjectID string     `json:"project_id,omitempty" db:"project_id"`
	StartedAt time.Time  `json:"started_at" db:"started_at"`
	EndedAt   *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	Summary   string     `json:"summary,omitempty" db:"summary"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// Observation represents a recorded observation during a session.
type Observation struct {
	ID         int64     `json:"id" db:"id"`
	SessionID  string    `json:"session_id" db:"session_id"`
	AgentName  string    `json:"agent_name" db:"agent_name"`
	Type       string    `json:"type" db:"type"` // decision, pattern, learning, insight
	Content    string    `json:"content" db:"content"`
	Importance int       `json:"importance" db:"importance"` // 1-5
	Tags       string    `json:"tags,omitempty" db:"tags"`   // JSON array
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Summary represents a session or period summary.
type Summary struct {
	ID        int64     `json:"id" db:"id"`
	SessionID string    `json:"session_id,omitempty" db:"session_id"`
	Type      string    `json:"type" db:"type"` // session, daily, weekly
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Plan represents a development plan.
type Plan struct {
	ID        int64     `json:"id" db:"id"`
	SessionID string    `json:"session_id,omitempty" db:"session_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Status    string    `json:"status" db:"status"` // draft, active, completed
	FilePath  string    `json:"file_path,omitempty" db:"file_path"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TeamContext represents context from team members.
type TeamContext struct {
	UserName     string    `json:"user_name" db:"user_name"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
	Summary      string    `json:"summary" db:"summary"`
	ActivePlan   string    `json:"active_plan,omitempty" db:"active_plan"`
}

// Project represents a registered project with aggregated session data.
type Project struct {
	ID           string    `json:"id"`
	Path         string    `json:"path"`
	SessionCount int       `json:"session_count"`
	LastActivity time.Time `json:"last_activity"`
}

// ContextInjectResponse is the response for context injection.
type ContextInjectResponse struct {
	Session      *Session       `json:"session,omitempty"`
	Observations []Observation  `json:"observations,omitempty"`
	ActivePlan   *Plan          `json:"active_plan,omitempty"`
	TeamContext  []TeamContext  `json:"team_context,omitempty"`
	Markdown     string         `json:"markdown"`
}

// CreateSessionRequest is the request to create a new session.
type CreateSessionRequest struct {
	ID        string `json:"id" binding:"required"`
	UserName  string `json:"user_name" binding:"required"`
	ProjectID string `json:"project_id"`
}

// EndSessionRequest is the request to end a session.
type EndSessionRequest struct {
	Summary string `json:"summary,omitempty"`
}

// CreateObservationRequest is the request to create an observation.
type CreateObservationRequest struct {
	SessionID  string   `json:"session_id" binding:"required"`
	AgentName  string   `json:"agent_name"`
	Type       string   `json:"type" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Importance int      `json:"importance"`
	Tags       []string `json:"tags,omitempty"`
}

// CreateSummaryRequest is the request to create a summary.
type CreateSummaryRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Type      string `json:"type" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// CreatePlanRequest is the request to create a plan.
type CreatePlanRequest struct {
	SessionID string `json:"session_id,omitempty"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	FilePath  string `json:"file_path,omitempty"`
}

// HealthResponse is the health check response.
type HealthResponse struct {
	Status   string `json:"status"`
	DBType   string `json:"db_type"`
	DBStatus string `json:"db_status"`
	Version  string `json:"version"`
}

// ErrorResponse is a standard error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}
