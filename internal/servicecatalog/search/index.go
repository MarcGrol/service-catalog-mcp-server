package search

import (
	"context"

	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"

	"github.com/MarcGrol/learnmcp/internal/servicecatalog/catalogrepo"
)

type Index interface {
	Search(ctx context.Context, keyword string) SearchResult
}

type searchIndex struct {
	Modules    []string
	Teams      []string
	Interfaces []string
	Databases  []string
}

func NewSearchIndex(ctx context.Context, cataloger catalogrepo.Cataloger) Index {
	modules, _ := cataloger.ListModules(ctx, "")
	interfaces, _ := cataloger.ListInterfaces(ctx)
	databases, _ := cataloger.ListDatabases(ctx)
	teams, _ := cataloger.ListTeams(ctx)

	return &searchIndex{
		Modules: lo.Map(modules, func(m catalogrepo.Module, index int) string {
			return m.ModuleID
		}),
		Interfaces: lo.Map(interfaces, func(m catalogrepo.Interface, index int) string {
			return m.InterfaceID
		}),
		Teams:     teams,
		Databases: databases,
	}
}

type SearchResult struct {
	Modules    []string
	Teams      []string
	Interfaces []string
	Databases  []string
}

func (idx *searchIndex) Search(ctx context.Context, keyword string) SearchResult {
	return SearchResult{
		Modules:    matchesToSlice(fuzzy.Find(keyword, idx.Modules)),
		Teams:      matchesToSlice(fuzzy.Find(keyword, idx.Teams)),
		Interfaces: matchesToSlice(fuzzy.Find(keyword, idx.Interfaces)),
		Databases:  matchesToSlice(fuzzy.Find(keyword, idx.Databases)),
	}
}

func matchesToSlice(matches fuzzy.Matches) []string {
	slice := lo.Map(matches, func(item fuzzy.Match, index int) string {
		return item.Str
	})
	// Limit to top 5 per category
	return slice[0:min(len(slice), 5)]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
