package search

import (
	"context"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

//go:generate mockgen -source=index.go -destination=mock_index.go -package=search Index
type Index interface {
	Search(ctx context.Context, keyword string, limit int) SearchResult
}

type searchIndex struct {
	Modules    []string
	Teams      []string
	Interfaces []string
	Databases  []string
	Flows      []string
}

func NewSearchIndex(ctx context.Context, cataloger catalogrepo.Cataloger) Index {
	modules, _ := cataloger.ListModules(ctx, "")
	interfaces, _ := cataloger.ListInterfaces(ctx, "")
	databases, _ := cataloger.ListDatabases(ctx)
	teams, _ := cataloger.ListTeams(ctx)
	flows, _ := cataloger.ListFlows(ctx)

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
	}
}

type SearchResult struct {
	Modules    []string
	Teams      []string
	Interfaces []string
	Databases  []string
	Flows      []string
}

func (idx *searchIndex) Search(ctx context.Context, keyword string, limit int) SearchResult {
	return SearchResult{
		Modules:    matchesToSlice(fuzzy.Find(keyword, idx.Modules), limit),
		Teams:      matchesToSlice(fuzzy.Find(keyword, idx.Teams), limit),
		Interfaces: matchesToSlice(fuzzy.Find(keyword, idx.Interfaces), limit),
		Databases:  matchesToSlice(fuzzy.Find(keyword, idx.Databases), limit),
		Flows:      matchesToSlice(fuzzy.Find(keyword, idx.Flows), limit*2),
	}
}

func matchesToSlice(matches fuzzy.Matches, limit int) []string {
	slice := lo.Map(matches, func(item fuzzy.Match, index int) string {
		return item.Str
	})
	// Limit to top 5 per category
	return slice[0:min(len(slice), limit)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
