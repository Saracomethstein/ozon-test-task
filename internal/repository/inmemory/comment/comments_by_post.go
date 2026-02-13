package comment

import (
	"context"
	"sort"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *comment) GetRootByPost(ctx context.Context, postID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids, ok := r.byPost[postID]
	if !ok {
		return []*models.Comment{}, nil
	}

	var roots []*models.Comment
	for _, id := range ids {
		c := r.comments[id]
		if c.ParentID == nil {
			roots = append(roots, c)
		}
	}

	sort.Slice(roots, func(i, j int) bool {
		if roots[i].CreatedAt == roots[j].CreatedAt {
			return roots[i].ID > roots[j].ID
		}
		return roots[i].CreatedAt > roots[j].CreatedAt
	})

	startIdx := 0
	if afterCreatedAt != nil && afterID > 0 {
		for i, c := range roots {
			if c.CreatedAt == *afterCreatedAt && c.ID == afterID {
				startIdx = i + 1
				break
			}
		}
	}
	if startIdx >= len(roots) {
		return []*models.Comment{}, nil
	}

	endIdx := min(startIdx+int(limit), len(roots))

	result := make([]*models.Comment, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		clone := *roots[i]
		result[i-startIdx] = &clone
	}

	return result, nil
}

func (r *comment) GetChild(ctx context.Context, parentID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids, ok := r.byParent[parentID]
	if !ok {
		return []*models.Comment{}, nil
	}

	var children []*models.Comment
	for _, id := range ids {
		children = append(children, r.comments[id])
	}

	sort.Slice(children, func(i, j int) bool {
		if children[i].CreatedAt == children[j].CreatedAt {
			return children[i].ID > children[j].ID
		}

		return children[i].CreatedAt > children[j].CreatedAt
	})

	startIdx := 0
	if afterCreatedAt != nil && afterID > 0 {
		for i, c := range children {
			if c.CreatedAt == *afterCreatedAt && c.ID == afterID {
				startIdx = i + 1
				break
			}
		}
	}
	if startIdx >= len(children) {
		return []*models.Comment{}, nil
	}

	endIdx := min(startIdx+int(limit), len(children))

	result := make([]*models.Comment, endIdx-startIdx)
	for i := startIdx; i < endIdx; i++ {
		clone := *children[i]
		result[i-startIdx] = &clone
	}

	return result, nil
}

func (r *comment) GetChildBatch(ctx context.Context, parentIDs []int64) ([]*models.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	need := make(map[int64]bool, len(parentIDs))
	for _, pid := range parentIDs {
		need[pid] = true
	}

	var result []*models.Comment
	for _, c := range r.comments {
		if c.ParentID != nil && need[*c.ParentID] {
			clone := *c
			result = append(result, &clone)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		pi := *result[i].ParentID
		pj := *result[j].ParentID

		if pi != pj {
			return pi < pj
		}

		if result[i].CreatedAt == result[j].CreatedAt {
			return result[i].ID > result[j].ID
		}

		return result[i].CreatedAt > result[j].CreatedAt
	})

	return result, nil
}

func (r *comment) TotalCount(ctx context.Context, postID int64) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids, ok := r.byPost[postID]
	if !ok {
		return 0, nil
	}

	return int64(len(ids)), nil
}
