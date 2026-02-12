package comment

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *commentResolver) Children(ctx context.Context, obj *graphql.Comment, first *int32, after *string) (*graphql.CommentConnection, error) {
	parentID, err := strconv.ParseInt(obj.ID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid comment ID")
	}

	conn, err := r.service.CommentService.Children(ctx, parentID, first, after)
	if err != nil {
		return nil, err
	}

	return convertToGraphQLCommentConnection(conn), nil
}

func convertToGraphQLCommentConnection(conn *models.CommentConnection) *graphql.CommentConnection {
	if conn == nil {
		return nil
	}

	edges := make([]*graphql.CommentEdge, len(conn.Edges))
	for i, edge := range conn.Edges {
		node := &graphql.Comment{
			ID:        strconv.FormatInt(edge.Node.ID, 10),
			PostID:    strconv.FormatInt(edge.Node.PostID, 10),
			Author:    edge.Node.Author,
			Text:      edge.Node.Text,
			CreatedAt: edge.Node.CreatedAt,
		}

		if edge.Node.ParentID != nil {
			pid := strconv.FormatInt(*edge.Node.ParentID, 10)
			node.ParentID = &pid
		}

		edges[i] = &graphql.CommentEdge{
			Cursor: edge.Cursor,
			Node:   node,
		}
	}

	return &graphql.CommentConnection{
		Edges: edges,
		PageInfo: &graphql.PageInfo{
			EndCursor:   conn.PageInfo.EndCursor,
			HasNextPage: conn.PageInfo.HasNextPage,
		},
		TotalCount: conn.TotalCount,
	}
}
