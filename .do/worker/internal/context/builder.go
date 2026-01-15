// Package context provides context building utilities for session injection.
package context

import (
	"context"

	"github.com/do-focus/worker/internal/db"
	"github.com/do-focus/worker/pkg/models"
)

// Builder constructs context for session injection.
type Builder struct {
	db       db.Adapter
	renderer *Renderer
}

// NewBuilder creates a new context builder.
func NewBuilder(adapter db.Adapter) *Builder {
	return &Builder{
		db:       adapter,
		renderer: NewRenderer(),
	}
}

// BuildContext assembles the full context for a user.
func (b *Builder) BuildContext(ctx context.Context, userName string, opts BuildOptions) (*models.ContextInjectResponse, error) {
	resp := &models.ContextInjectResponse{}

	// Get latest session
	session, err := b.db.GetLatestSession(ctx, userName)
	if err != nil {
		return nil, err
	}
	resp.Session = session

	// Get recent observations
	limit := opts.ObservationLimit
	if limit <= 0 {
		limit = 20
	}
	observations, err := b.db.GetRecentObservations(ctx, userName, limit)
	if err != nil {
		return nil, err
	}
	resp.Observations = observations

	// Get active plan
	if opts.IncludePlan {
		plan, err := b.db.GetActivePlan(ctx, userName)
		if err != nil {
			return nil, err
		}
		resp.ActivePlan = plan
	}

	// Get team context
	if opts.IncludeTeam {
		team, err := b.db.GetTeamContext(ctx, userName)
		if err != nil {
			return nil, err
		}
		resp.TeamContext = team
	}

	// Render markdown
	resp.Markdown = b.renderer.RenderContext(resp)

	return resp, nil
}

// BuildOptions configures context building.
type BuildOptions struct {
	ObservationLimit int
	IncludePlan      bool
	IncludeTeam      bool
	IncludeSession   bool
}

// DefaultBuildOptions returns the default build options.
func DefaultBuildOptions() BuildOptions {
	return BuildOptions{
		ObservationLimit: 20,
		IncludePlan:      true,
		IncludeTeam:      true,
		IncludeSession:   true,
	}
}

// MinimalBuildOptions returns minimal build options.
func MinimalBuildOptions() BuildOptions {
	return BuildOptions{
		ObservationLimit: 10,
		IncludePlan:      false,
		IncludeTeam:      false,
		IncludeSession:   true,
	}
}

// FullBuildOptions returns full build options.
func FullBuildOptions() BuildOptions {
	return BuildOptions{
		ObservationLimit: 50,
		IncludePlan:      true,
		IncludeTeam:      true,
		IncludeSession:   true,
	}
}
