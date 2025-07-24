package search

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"

	"github.com/MarcGrol/service-catalog-mcp-server/internal/servicecatalog/catalogrepo"
)

// Index defines the interface for a search index.
//
//go:generate mockgen -source=index.go -destination=mock_index.go -package=search Index
type Index interface {
	Search(ctx context.Context, keyword string, limit int) Result
}

type searchIndex struct {
	Modules    []string
	Teams      []string
	Interfaces []string
	Databases  []string
	Flows      []string
	Kinds      []string
}

// NewSearchIndex creates a new search index.
func NewSearchIndex(ctx context.Context, cataloger catalogrepo.Cataloger) Index {
	modules, err := cataloger.ListModules(ctx, "")
	if err != nil {
		log.Error().Err(err).Msg("Error listing modules for search index")
	}
	interfaces, err := cataloger.ListInterfaces(ctx, "")
	if err != nil {
		log.Error().Err(err).Msg("Error listing interfaces for search index")
	}
	databases, err := cataloger.ListDatabases(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing databases for search index")
	}
	teams, err := cataloger.ListTeams(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing teams for search index")
	}
	flows, err := cataloger.ListFlows(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing flows for search index")
	}
	kinds, err := cataloger.ListKinds(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error listing kinds for search index")
	}

	return &searchIndex{
		Modules: lo.Map(modules, func(m catalogrepo.Module, index int) string {
			return m.ModuleID
		}),
		Interfaces: lo.Map(interfaces, func(m catalogrepo.Interface, index int) string {
			return m.InterfaceID
		}),
		Teams:     teams,
		Databases: databases,
		Flows:     flows,
		Kinds:     kinds,
	}
}

// Result represents the search results.
type Result struct {
	Modules    []string
	Teams      []string
	Interfaces []string
	Databases  []string
	Flows      []string
	Kinds      []string
}

const flowSearchLimitMultiplier = 2

func (idx *searchIndex) Search(ctx context.Context, keyword string, limit int) Result {
	return Result{
		Modules:    matchesToSlice(fuzzy.Find(keyword, idx.Modules), limit),
		Teams:      matchesToSlice(fuzzy.Find(keyword, idx.Teams), limit),
		Interfaces: matchesToSlice(fuzzy.Find(keyword, idx.Interfaces), limit),
		Databases:  matchesToSlice(fuzzy.Find(keyword, idx.Databases), limit),
		Flows:      matchesToSlice(fuzzy.Find(keyword, idx.Flows), limit*flowSearchLimitMultiplier),
		Kinds:      matchesToSlice(fuzzy.Find(keyword, idx.Kinds), limit*4),
	}
}

func matchesToSlice(matches fuzzy.Matches, limit int) []string {
	slice := lo.Map(matches, func(item fuzzy.Match, index int) string {
		return item.Str
	})
	// Limit to top 5 per category
	return slice[0:min(len(slice), limit)]
}
