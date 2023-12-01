package search

import (
	"context"

	"github.com/stackrox/rox/central/policycategory/index"
	"github.com/stackrox/rox/central/policycategory/store"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/search"
)

// Searcher provides search functionality on existing alerts
//
//go:generate mockgen-wrapper
type Searcher interface {
	Search(ctx context.Context, q *v1.Query) ([]search.Result, error)
	Count(ctx context.Context, q *v1.Query) (int, error)
	SearchCategories(ctx context.Context, q *v1.Query) ([]*v1.SearchResult, error)
	SearchRawCategories(ctx context.Context, q *v1.Query) ([]*storage.PolicyCategory, error)
}

// New returns a new instance of Searcher for the given storage and indexer.
func New(storage store.Store, indexer index.Indexer) Searcher {
	return &searcherImpl{
		storage:  storage,
		indexer:  indexer,
		searcher: formatSearcher(indexer),
	}
}
