package post

import (
	"context"
	"sort"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *post) Get(ctx context.Context, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Post, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	all := make([]*models.Post, 0, len(r.posts))
	for _, p := range r.posts {
		all = append(all, p)
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].CreatedAt == all[j].CreatedAt {
			return all[i].ID > all[j].ID
		}
		return all[i].CreatedAt > all[j].CreatedAt
	})

	startIdx := 0
	if afterCreatedAt != nil && afterID > 0 {
		for i, p := range all {
			if p.CreatedAt == *afterCreatedAt && p.ID == afterID {
				startIdx = i + 1
				break
			}
		}
	}

	if startIdx >= len(all) {
		return []*models.Post{}, nil
	}

	endIdx := min(startIdx+int(limit), len(all))

	result := make([]*models.Post, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		clone := *all[i]
		result[i-startIdx] = &clone
	}
	return result, nil
}

func (r *post) TotalCount(ctx context.Context) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return int64(len(r.posts)), nil
}
