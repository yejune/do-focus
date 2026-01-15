package server

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/do-focus/worker/pkg/models"
	"github.com/gin-gonic/gin"
)

// Version is set by main package
var Version = "dev"

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
		api.POST("/summaries/generate", s.handleGenerateSummary)

		// User Prompts
		api.GET("/prompts", s.handleGetUserPrompts)
		api.POST("/prompts", s.handleCreateUserPrompt)

		// FTS5 Search
		api.GET("/search", s.handleSearch)

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
		Version:  Version,
	})
}

// handleContextInject handles context injection for SessionStart hook.
// level parameter controls the amount of data returned:
// - level 1: minimal (session only)
// - level 2: standard (session + observations) [default]
// - level 3: full (session + observations + plan + team)
func (s *Server) handleContextInject(c *gin.Context) {
	ctx := c.Request.Context()

	userName := c.Query("user")
	if userName == "" {
		userName = os.Getenv("DO_USER_NAME")
	}
	if userName == "" {
		userName = "default"
	}

	// Parse level parameter (1-3, default 2)
	levelStr := c.DefaultQuery("level", "2")
	level, _ := strconv.Atoi(levelStr)
	if level < 1 {
		level = 1
	}
	if level > 3 {
		level = 3
	}

	// Get latest session (always included)
	session, err := s.db.GetLatestSession(ctx, userName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	var observations []models.Observation
	var plan *models.Plan
	var teamContext []models.TeamContext

	// Level 2+: Include observations
	if level >= 2 {
		limitStr := c.DefaultQuery("obs_limit", "20")
		limit, _ := strconv.Atoi(limitStr)
		if limit <= 0 {
			limit = 20
		}
		// Adjust limit based on level
		if level == 2 {
			if limit > 20 {
				limit = 20
			}
		}

		observations, err = s.db.GetRecentObservations(ctx, userName, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "database_error",
				Message: err.Error(),
			})
			return
		}
	}

	// Level 3: Include plan and team context
	if level >= 3 {
		plan, err = s.db.GetActivePlan(ctx, userName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "database_error",
				Message: err.Error(),
			})
			return
		}

		teamContext, err = s.db.GetTeamContext(ctx, userName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "database_error",
				Message: err.Error(),
			})
			return
		}
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

// handleCreateSession handles session creation (idempotent).
func (s *Server) handleCreateSession(c *gin.Context) {
	var req models.CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// Check if session already exists (idempotent)
	existing, _ := s.db.GetSession(c.Request.Context(), req.ID)
	if existing != nil {
		c.JSON(http.StatusOK, existing)
		return
	}

	session := &models.Session{
		ID:        req.ID,
		UserName:  req.UserName,
		ProjectID: req.ProjectID,
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

// handleGenerateSummary generates a rule-based summary from session observations.
func (s *Server) handleGenerateSummary(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.GenerateSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	// 1. Verify session exists
	session, err := s.db.GetSession(ctx, req.SessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "session_not_found",
			Message: "Session not found: " + req.SessionID,
		})
		return
	}

	// 2. Get observations for the session
	observations, err := s.db.GetObservationsFiltered(ctx, req.SessionID, "", 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	// 3. Generate rule-based summary
	summaryContent := generateRuleBasedSummary(observations, req.LastAssistantMessage)

	// 4. Save summary to DB
	summary := &models.Summary{
		SessionID: req.SessionID,
		Type:      "session",
		Content:   summaryContent,
	}

	if err := s.db.CreateSummary(ctx, summary); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, summary)
}

// generateRuleBasedSummary extracts key information from observations using rules.
func generateRuleBasedSummary(observations []models.Observation, lastMessage string) string {
	var completed []string
	var decisions []string
	var learnings []string
	var request string

	// Extract request from last message (first sentence)
	if lastMessage != "" {
		sentences := strings.SplitN(lastMessage, ".", 2)
		if len(sentences) > 0 {
			request = strings.TrimSpace(sentences[0])
		}
	}

	// Categorize observations by type
	for _, obs := range observations {
		switch obs.Type {
		case "decision":
			decisions = append(decisions, obs.Content)
			completed = append(completed, obs.Content)
		case "bugfix":
			completed = append(completed, "Fixed: "+obs.Content)
		case "feature":
			completed = append(completed, "Implemented: "+obs.Content)
		case "learning", "insight":
			learnings = append(learnings, obs.Content)
		case "pattern":
			learnings = append(learnings, "Pattern: "+obs.Content)
		}
	}

	// Build summary content
	var parts []string

	if request != "" {
		parts = append(parts, "## Request\n"+request)
	}

	if len(completed) > 0 {
		completedStr := "## Completed\n"
		for _, c := range completed {
			completedStr += "- " + c + "\n"
		}
		parts = append(parts, completedStr)
	}

	if len(decisions) > 0 {
		decisionsStr := "## Decisions\n"
		for _, d := range decisions {
			decisionsStr += "- " + d + "\n"
		}
		parts = append(parts, decisionsStr)
	}

	if len(learnings) > 0 {
		learningsStr := "## Learnings\n"
		for _, l := range learnings {
			learningsStr += "- " + l + "\n"
		}
		parts = append(parts, learningsStr)
	}

	if len(parts) == 0 {
		return "No significant observations recorded in this session."
	}

	return strings.Join(parts, "\n")
}

// handleGetUserPrompts handles user prompt list retrieval.
func (s *Server) handleGetUserPrompts(c *gin.Context) {
	ctx := c.Request.Context()

	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 100
	}

	sessionID := c.Query("session_id")

	prompts, err := s.db.GetUserPrompts(ctx, sessionID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, prompts)
}

// handleCreateUserPrompt handles user prompt creation.
func (s *Server) handleCreateUserPrompt(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.CreateUserPromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	prompt := &models.UserPrompt{
		SessionID:      req.SessionID,
		PromptNumber:   req.PromptNumber,
		PromptText:     req.PromptText,
		CreatedAtEpoch: time.Now().Unix(),
	}

	if err := s.db.CreateUserPrompt(ctx, prompt); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, prompt)
}

// handleSearch handles FTS5 full-text search across multiple types.
func (s *Server) handleSearch(c *gin.Context) {
	ctx := c.Request.Context()

	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_request",
			Message: "Query parameter 'q' is required",
		})
		return
	}

	// Parse types parameter (comma-separated)
	typesStr := c.DefaultQuery("types", "observation,prompt")
	types := strings.Split(typesStr, ",")
	for i, t := range types {
		types[i] = strings.TrimSpace(t)
	}

	// Parse limit
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	results, err := s.db.SearchFTS(ctx, query, types, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.SearchResponse{
		Results: results,
		Query:   query,
		Total:   len(results),
	})
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
