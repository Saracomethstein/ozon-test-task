package comment

import (
	"context"
	"strconv"

	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"

	myLoader "github.com/Saracomethstein/ozon-test-task/internal/graphql/dataloader"
	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/utils/cursor"
)

func (s *Service) Children(ctx context.Context, parentID int64, first *int32, after *string) (*models.CommentConnection, error) {
	loader, ok := ctx.Value(myLoader.Key).(dataloader.Interface)
	if !ok {
		return nil, errors.New("dataloader not found in context")
	}

	thunk := loader.Load(ctx, dataloader.StringKey(strconv.FormatInt(parentID, 10)))
	result, err := thunk()
	if err != nil {
		return nil, err
	}

	children, ok := result.([]*models.Comment)
	if !ok {
		return nil, errors.New("unexpected data type from dataloader")
	}

	limit := s.getLimit(first)
	var startIdx int

	if after != nil && *after != "" {
		createdAt, id, err := cursor.Decode(*after)
		if err != nil {
			return nil, errors.New("invalid cursor format")
		}
		for i, c := range children {
			if c.CreatedAt == createdAt && c.ID == id {
				startIdx = i + 1
				break
			}
		}
	}

	if startIdx >= len(children) {
		return &models.CommentConnection{
			Edges:      []*models.CommentEdge{},
			PageInfo:   &models.PageInfo{HasNextPage: false},
			TotalCount: int32(len(children)),
		}, nil
	}

	endIdx := startIdx + int(limit)
	hasNextPage := endIdx < len(children)
	if endIdx > len(children) {
		endIdx = len(children)
		hasNextPage = false
	}

	page := children[startIdx:endIdx]

	edges := make([]*models.CommentEdge, len(page))
	for i, c := range page {
		edges[i] = &models.CommentEdge{
			Cursor: cursor.Encode(c.CreatedAt, c.ID),
			Node:   c,
		}
	}

	var endCursor *string
	if len(edges) > 0 {
		endCursor = &edges[len(edges)-1].Cursor
	}

	return &models.CommentConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		},
		TotalCount: int32(len(children)),
	}, nil
}
