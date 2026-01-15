package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/do-focus/worker/pkg/models"
	"github.com/gin-gonic/gin"
)

const version = "0.1.0"

// setupRoutes configures all API routes.
func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", s.handleHealth)

	// API routes
	api := s.router.Group("/api")
	{
		// Context injection for SessionStart hook
		api.GET("/context/inject", s.handleContextInject)

		// Session management
		api.GET("/sessions", s.handleGetSessions)
		api.GET("/sessions/:id", s.handleGetSession)
		api.POST("/sessions", s.handleCreateSession)
		api.PUT("/sessions/:id/end", s.handleEndSession)

		// Observations
		api.GET("/observations", s.handleGetObservations)
		api.GET("/observations/search", s.handleSearchObservations)
		api.POST("/observations", s.handleCreateObservation)

		// Summaries
		api.GET("/summaries", s.handleGetSummaries)
		api.POST("/summaries", s.handleCreateSummary)

		// Plans
		api.GET("/plans", s.handleGetPlans)
		api.POST("/plans", s.handleCreatePlan)

		// Team context
		api.GET("/team/context", s.handleTeamContext)

		// Projects
		api.GET("/projects", s.getProjects)
	}
}

// handleHealth handles the health check endpoint.
func (s *Server) handleHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	dbStatus := "ok"
	if err := s.db.Health(ctx); err != nil {
		dbStatus = "error: " + err.Error()
	}

	dbType := os.Getenv("DO_DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	c.JSON(http.StatusOK, models.HealthResponse{
		Status:   "ok",
		DBType:   dbType,
		DBStatus: dbStatus,
		Version:  version,
	})
}

// handleContextInject handles context injection for SessionStart hook.
func (s *Server) handleContextInject(c *gin.Context) {
	ctx := c.Request.Context()

	userName := c.Query("user")
	if userName == "" {
		userName = os.Getenv("DO_USER_NAME")
	}
	if userName == "" {
		userName = "default"
	}

	// Get latest session
	session, err := s.db.GetLatestSession(ctx, userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	// Get recent observations
	limitStr := c.DefaultQuery("obs_limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	observations, err := s.db.GetRecentObservations(ctx, userName, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	// Get active plan
	plan, err := s.db.GetActivePlan(ctx, userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	// Get team context
	teamContext, err := s.db.GetTeamContext(ctx, userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	// Build markdown response
	markdown := buildContextMarkdown(session, observations, plan, teamContext)

	c.JSON(http.StatusOK, models.ContextInjectResponse{
		Session:      session,
		Observations: observations,
		ActivePlan:   plan,
		TeamContext:  teamContext,
		Markdown:     markdown,
	})
}

// handleCreateSession handles session creation.
func (s *Server) handleCreateSession(c *gin.Context) {
	var req models.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	session := &models.Session{
		ID:        req.ID,
		UserName:  req.UserName,
		StartedAt: time.Now(),
	}

	if err := s.db.CreateSession(c.Request.Context(), session); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, session)
}

// handleEndSession handles session ending.
func (s *Server) handleEndSession(c *gin.Context) {
	id := c.Param("id")

	var req models.EndSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body
		req = models.EndSessionRequest{}
	}

	if err := s.db.EndSession(c.Request.Context(), id, req.Summary); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ended"})
}

// handleGetSessions handles session list retrieval.
func (s *Server) handleGetSessions(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}

	sessions, err := s.db.GetRecentSessions(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// handleGetSession handles single session retrieval.
func (s *Server) handleGetSession(c *gin.Context) {
	id := c.Param("id")

	session, err := s.db.GetSession(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	if session == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "Session not found",
		})
		return
	}

	c.JSON(http.StatusOK, session)
}

// handleGetObservations handles observation list retrieval.
func (s *Server) handleGetObservations(c *gin.Context) {
	sessionID := c.Query("session_id")
	obsType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	observations, err := s.db.GetObservationsFiltered(c.Request.Context(), sessionID, obsType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, observations)
}

// handleSearchObservations handles observation search.
func (s *Server) handleSearchObservations(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Query parameter 'q' is required",
		})
		return
	}

	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	results, err := s.db.SearchObservations(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

// handleGetSummaries handles summary list retrieval.
func (s *Server) handleGetSummaries(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	days, _ := strconv.Atoi(daysStr)
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	summaries, err := s.db.GetAllSummaries(c.Request.Context(), days, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, summaries)
}

// handleGetPlans handles plan list retrieval.
func (s *Server) handleGetPlans(c *gin.Context) {
	sessionID := c.Query("session_id")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	plans, err := s.db.GetAllPlans(c.Request.Context(), sessionID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// handleCreateObservation handles observation creation.
func (s *Server) handleCreateObservation(c *gin.Context) {
	var req models.CreateObservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Default importance to 3
	if req.Importance <= 0 || req.Importance > 5 {
		req.Importance = 3
	}

	// Convert tags to JSON
	var tagsJSON string
	if len(req.Tags) > 0 {
		tagsBytes, _ := json.Marshal(req.Tags)
		tagsJSON = string(tagsBytes)
	}

	obs := &models.Observation{
		SessionID:  req.SessionID,
		AgentName:  req.AgentName,
		Type:       req.Type,
		Content:    req.Content,
		Importance: req.Importance,
		Tags:       tagsJSON,
	}

	if err := s.db.CreateObservation(c.Request.Context(), obs); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, obs)
}

// handleCreateSummary handles summary creation.
func (s *Server) handleCreateSummary(c *gin.Context) {
	var req models.CreateSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	summary := &models.Summary{
		SessionID: req.SessionID,
		Type:      req.Type,
		Content:   req.Content,
	}

	if err := s.db.CreateSummary(c.Request.Context(), summary); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, summary)
}

// handleCreatePlan handles plan creation.
func (s *Server) handleCreatePlan(c *gin.Context) {
	var req models.CreatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	plan := &models.Plan{
		SessionID: req.SessionID,
		Title:     req.Title,
		Content:   req.Content,
		Status:    "draft",
		FilePath:  req.FilePath,
	}

	if err := s.db.CreatePlan(c.Request.Context(), plan); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, plan)
}

// handleTeamContext handles team context retrieval.
func (s *Server) handleTeamContext(c *gin.Context) {
	userName := c.Query("exclude_user")
	if userName == "" {
		userName = os.Getenv("DO_USER_NAME")
	}

	contexts, err := s.db.GetTeamContext(c.Request.Context(), userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"team": contexts})
}

// getProjects handles project list retrieval.
func (s *Server) getProjects(c *gin.Context) {
	projects, err := s.db.GetProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// buildContextMarkdown builds a markdown representation of the context.
func buildContextMarkdown(session *models.Session, observations []models.Observation, plan *models.Plan, team []models.TeamContext) string {
	var md string

	md += "# Do Worker Context\n\n"

	// Session info
	if session != nil {
		md += "## Last Session\n"
		md += "- ID: " + session.ID + "\n"
		md += "- Started: " + session.StartedAt.Format(time.RFC3339) + "\n"
		if session.EndedAt != nil {
			md += "- Ended: " + session.EndedAt.Format(time.RFC3339) + "\n"
		}
		if session.Summary != "" {
			md += "- Summary: " + session.Summary + "\n"
		}
		md += "\n"
	}

	// Active plan
	if plan != nil {
		md += "## Active Plan\n"
		md += "**" + plan.Title + "**\n\n"
		md += plan.Content + "\n\n"
	}

	// Recent observations
	if len(observations) > 0 {
		md += "## Recent Observations\n"
		for _, obs := range observations {
			importance := ""
			if obs.Importance >= 4 {
				importance = " [HIGH]"
			}
			md += "- [" + obs.Type + "]" + importance + " " + obs.Content
			if obs.AgentName != "" {
				md += " (by " + obs.AgentName + ")"
			}
			md += "\n"
		}
		md += "\n"
	}

	// Team context
	if len(team) > 0 {
		md += "## Team Activity\n"
		for _, t := range team {
			md += "- **" + t.UserName + "**: " + t.Summary
			if t.ActivePlan != "" {
				md += " [Working on: " + t.ActivePlan + "]"
			}
			md += "\n"
		}
	}

	return md
}
