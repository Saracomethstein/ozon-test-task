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

	connection, err := r.service.CommentService.GetComments(ctx, postID, first, after)
	if err != nil {
		return nil, err
	}

	edges := make([]*graphql.CommentEdge, 0, len(connection.Edges))
	for _, edge := range connection.Edges {
		var parentIDPtr *string
		if edge.Node.ParentID != nil {
			pid := strconv.FormatInt(*edge.Node.ParentID, 10)
			parentIDPtr = &pid
		}

		node := &graphql.Comment{
			ID:        strconv.FormatInt(edge.Node.ID, 10),
			PostID:    strconv.FormatInt(edge.Node.PostID, 10),
			ParentID:  parentIDPtr,
			Author:    edge.Node.Author,
			Text:      edge.Node.Text,
			Path:      edge.Node.Path,
			CreatedAt: edge.Node.CreatedAt,
		}

		edges = append(edges, &graphql.CommentEdge{
			Cursor: edge.Cursor,
			Node:   node,
		})
	}

	pageInfo := &graphql.PageInfo{
		EndCursor:   connection.PageInfo.EndCursor,
		HasNextPage: connection.PageInfo.HasNextPage,
	}

	commentsConnection := &graphql.CommentConnection{
		Edges:      edges,
		PageInfo:   pageInfo,
		TotalCount: connection.TotalCount,
	}

	return commentsConnection, nil
}
