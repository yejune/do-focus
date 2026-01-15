package context

import (
	"fmt"
	"strings"
	"time"

	"github.com/do-focus/worker/pkg/models"
)

// Renderer converts context data to markdown format.
type Renderer struct {
	maxObservations int
	maxTeamMembers  int
}

// NewRenderer creates a new renderer with default settings.
func NewRenderer() *Renderer {
	return &Renderer{
		maxObservations: 20,
		maxTeamMembers:  5,
	}
}

// RenderContext generates markdown from context data.
func (r *Renderer) RenderContext(ctx *models.ContextInjectResponse) string {
	var sb strings.Builder

	sb.WriteString("# Do Memory Context\n\n")

	// Render session
	if ctx.Session != nil {
		r.renderSession(&sb, ctx.Session)
	}

	// Render active plan
	if ctx.ActivePlan != nil {
		r.renderPlan(&sb, ctx.ActivePlan)
	}

	// Render observations
	if len(ctx.Observations) > 0 {
		r.renderObservations(&sb, ctx.Observations)
	}

	// Render team context
	if len(ctx.TeamContext) > 0 {
		r.renderTeamContext(&sb, ctx.TeamContext)
	}

	return sb.String()
}

// renderSession renders session information.
func (r *Renderer) renderSession(sb *strings.Builder, session *models.Session) {
	sb.WriteString("## Current Session\n\n")
	sb.WriteString(fmt.Sprintf("- **ID**: `%s`\n", session.ID))
	sb.WriteString(fmt.Sprintf("- **User**: %s\n", session.UserName))
	sb.WriteString(fmt.Sprintf("- **Started**: %s\n", session.StartedAt.Format(time.RFC3339)))

	if session.EndedAt != nil {
		sb.WriteString(fmt.Sprintf("- **Ended**: %s\n", session.EndedAt.Format(time.RFC3339)))
		duration := session.EndedAt.Sub(session.StartedAt)
		sb.WriteString(fmt.Sprintf("- **Duration**: %s\n", formatDuration(duration)))
	} else {
		duration := time.Since(session.StartedAt)
		sb.WriteString(fmt.Sprintf("- **Active for**: %s\n", formatDuration(duration)))
	}

	if session.Summary != "" {
		sb.WriteString(fmt.Sprintf("\n**Summary**: %s\n", session.Summary))
	}
	sb.WriteString("\n")
}

// renderPlan renders the active plan.
func (r *Renderer) renderPlan(sb *strings.Builder, plan *models.Plan) {
	sb.WriteString("## Active Plan\n\n")
	sb.WriteString(fmt.Sprintf("### %s\n\n", plan.Title))
	sb.WriteString(fmt.Sprintf("**Status**: %s\n\n", plan.Status))

	if plan.FilePath != "" {
		sb.WriteString(fmt.Sprintf("**File**: `%s`\n\n", plan.FilePath))
	}

	sb.WriteString(plan.Content)
	sb.WriteString("\n\n")
}

// renderObservations renders recent observations.
func (r *Renderer) renderObservations(sb *strings.Builder, observations []models.Observation) {
	sb.WriteString("## Recent Observations\n\n")

	// Group by type
	byType := make(map[string][]models.Observation)
	for _, obs := range observations {
		byType[obs.Type] = append(byType[obs.Type], obs)
	}

	// Render high priority first
	highPriority := filterHighPriority(observations)
	if len(highPriority) > 0 {
		sb.WriteString("### High Priority\n\n")
		for _, obs := range highPriority {
			r.renderObservation(sb, obs)
		}
		sb.WriteString("\n")
	}

	// Render by type
	typeOrder := []string{"decision", "insight", "learning", "pattern"}
	for _, t := range typeOrder {
		if obs, ok := byType[t]; ok {
			sb.WriteString(fmt.Sprintf("### %s\n\n", capitalize(t)))
			count := 0
			for _, o := range obs {
				if o.Importance >= 4 {
					continue // Already rendered in high priority
				}
				r.renderObservation(sb, o)
				count++
				if count >= 5 {
					break
				}
			}
			sb.WriteString("\n")
		}
	}
}

// renderObservation renders a single observation.
func (r *Renderer) renderObservation(sb *strings.Builder, obs models.Observation) {
	importance := ""
	if obs.Importance >= 4 {
		importance = " **[!]**"
	}

	agent := ""
	if obs.AgentName != "" {
		agent = fmt.Sprintf(" _(by %s)_", obs.AgentName)
	}

	sb.WriteString(fmt.Sprintf("- %s%s%s\n", obs.Content, importance, agent))
}

// renderTeamContext renders team member activity.
func (r *Renderer) renderTeamContext(sb *strings.Builder, team []models.TeamContext) {
	sb.WriteString("## Team Activity\n\n")

	for i, t := range team {
		if i >= r.maxTeamMembers {
			break
		}

		sb.WriteString(fmt.Sprintf("### %s\n\n", t.UserName))
		sb.WriteString(fmt.Sprintf("- **Last Active**: %s\n", t.LastActivity.Format("2006-01-02 15:04")))

		if t.Summary != "" {
			sb.WriteString(fmt.Sprintf("- **Last Work**: %s\n", t.Summary))
		}

		if t.ActivePlan != "" {
			sb.WriteString(fmt.Sprintf("- **Working On**: %s\n", t.ActivePlan))
		}

		sb.WriteString("\n")
	}
}

// Helper functions

func filterHighPriority(observations []models.Observation) []models.Observation {
	var result []models.Observation
	for _, obs := range observations {
		if obs.Importance >= 4 {
			result = append(result, obs)
		}
	}
	return result
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if minutes > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dh", hours)
}
