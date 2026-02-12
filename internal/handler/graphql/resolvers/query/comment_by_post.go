package query

import (
	"context"
	"strconv"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/pkg/errors"
)

func (r *queryResolver) CommentsByPost(ctx context.Context, postID string, first *int32, after *string) (*graphql.CommentConnection, error) {
	if postID == "" {
		return nil, errors.New("postID cannot be empty")
	}

	connection, err := r.service.CommentService.GetRootComments(ctx, postID, first, after)
	if err != nil {
		return nil, err
	}

	edges := make([]*graphql.CommentEdge, len(connection.Edges))
	for i, edge := range connection.Edges {
		node := &graphql.Comment{
			ID:        strconv.FormatInt(edge.Node.ID, 10),
			PostID:    strconv.FormatInt(edge.Node.PostID, 10),
			ParentID:  nil,
			Author:    edge.Node.Author,
			Text:      edge.Node.Text,
			CreatedAt: edge.Node.CreatedAt,
		}
		edges[i] = &graphql.CommentEdge{
			Cursor: edge.Cursor,
			Node:   node,
		}
	}

	pageInfo := &graphql.PageInfo{
		EndCursor:   connection.PageInfo.EndCursor,
		HasNextPage: connection.PageInfo.HasNextPage,
	}

	return &graphql.CommentConnection{
		Edges:      edges,
		PageInfo:   pageInfo,
		TotalCount: connection.TotalCount,
	}, nil
}
