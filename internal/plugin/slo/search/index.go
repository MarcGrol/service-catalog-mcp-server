package search

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/plugin/slo/repo"
)

// Index defines the interface for a search index.
//
//go:generate go tool mockgen -source=index.go -destination=mock_index.go -package=search Index
type Index interface {
	Search(ctx context.Context, keyword string, limit int) Result
}

type searchIndex struct {
	SLOs         []string
	Teams        []string
	Applications []string
	Webapps      []string
	Services     []string
	Components   []string
	Methods      []string
}

// NewSearchIndex creates a new search index.
func NewSearchIndex(ctx context.Context, r repo.SLORepo) Index {
	slos, err := r.ListSLOs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing slos for search index")
	}

	sloNames := lo.Uniq(lo.Map(slos, func(slo repo.SLO, index int) string {
		return slo.UID
	}))
	teams := lo.Uniq(lo.Map(slos, func(m repo.SLO, index int) string {
		return m.Team
	}))
	applications := lo.Uniq(lo.Map(slos, func(slo repo.SLO, index int) string {
		return slo.Application
	}))
	webapps := lo.Uniq(lo.Map(slos, func(slo repo.SLO, index int) string {
		return slo.PromQLWebapp
	}))
	services1 := lo.Uniq(lo.Map(slos, func(slo repo.SLO, index int) string {
		return slo.Service
	}))
	services2 := lo.Uniq(lo.Map(slos, func(slo repo.SLO, index int) string {
		return slo.PromQLService
	}))
	services := lo.Uniq(append(services1, services2...))
	components := lo.Uniq(lo.Map(slos, func(slo repo.SLO, index int) string {
		return slo.Component
	}))
	methods := lo.Uniq(lo.Map(slos, func(slo repo.SLO, index int) string {
		return slo.PromQLMethods
	}))

	return &searchIndex{
		SLOs:         sloNames,
		Teams:        teams,
		Applications: applications,
		Webapps:      webapps,
		Services:     services,
		Components:   components,
		Methods:      methods,
	}
}

// Result represents the search results.
type Result struct {
	SLOs         []string
	Teams        []string
	Applications []string
	Webapps      []string
	Services     []string
	Components   []string
	Methods      []string
}

func (idx *searchIndex) Search(ctx context.Context, keyword string, limit int) Result {
	return Result{
		SLOs:         matchesToSlice(fuzzy.Find(keyword, idx.SLOs), limit),
		Teams:        matchesToSlice(fuzzy.Find(keyword, idx.Teams), limit),
		Applications: matchesToSlice(fuzzy.Find(keyword, idx.Applications), limit),
		Webapps:      matchesToSlice(fuzzy.Find(keyword, idx.Webapps), limit),
		Services:     matchesToSlice(fuzzy.Find(keyword, idx.Services), limit),
		Components:   matchesToSlice(fuzzy.Find(keyword, idx.Components), limit),
		Methods:      matchesToSlice(fuzzy.Find(keyword, idx.Methods), limit),
	}
}

func matchesToSlice(matches fuzzy.Matches, limit int) []string {
	filtered := lo.Filter(matches, func(item fuzzy.Match, index int) bool {
		return item.Score > 0
	})
	slice := lo.Map(filtered, func(item fuzzy.Match, index int) string {
		return item.Str
	})
	// Limit to top 'limit' per category
	return slice[0:min(len(slice), limit)]
}
