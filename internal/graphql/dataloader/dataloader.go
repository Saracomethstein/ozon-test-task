package dataloader

import (
	"context"
	"strconv"

	"github.com/graph-gophers/dataloader"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type ctxKey string

const Key = ctxKey("dataloader.comment.children")

func (l *CommentLoader) BatchGetChildren(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	parentIDs := make([]int64, len(keys))
	for i, key := range keys {
		id, err := strconv.ParseInt(key.String(), 10, 64)
		if err != nil {
			errResults := make([]*dataloader.Result, len(keys))
			for j := range errResults {
				errResults[j] = &dataloader.Result{Error: err}
			}
			return errResults
		}
		parentIDs[i] = id
	}

	comments, err := l.repo.GetChildCommentsBatch(ctx, parentIDs)
	if err != nil {
		errResults := make([]*dataloader.Result, len(keys))
		for i := range errResults {
			errResults[i] = &dataloader.Result{Error: err}
		}
		return errResults
	}

	groups := make(map[int64][]*models.Comment)
	for _, c := range comments {
		if c.ParentID != nil {
			groups[*c.ParentID] = append(groups[*c.ParentID], c)
		}
	}

	results := make([]*dataloader.Result, len(keys))
	for i, pid := range parentIDs {
		results[i] = &dataloader.Result{Data: groups[pid]}
	}

	return results
}
