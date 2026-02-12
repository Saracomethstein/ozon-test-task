package query

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
)

func (r *queryResolver) Posts(ctx context.Context, first *int32, after *string) (*graphql.PostConnection, error) {
	connection, err := r.service.PostService.GetPosts(ctx, first, after)
	if err != nil {
		return nil, err
	}

	edges := make([]*graphql.PostEdge, 0, len(connection.Edges))
	for _, edge := range connection.Edges {
		node := &graphql.Post{
			ID:            edge.Node.ID,
			Title:         edge.Node.Title,
			Body:          edge.Node.Body,
			Author:        edge.Node.Author,
			AllowComments: edge.Node.AllowComments,
			CreatedAt:     edge.Node.CreatedAt,
			// Comments: edge.Node.Comment,
		}

		edges = append(edges, &graphql.PostEdge{
			Cursor: edge.Cursor,
			Node:   node,
		})
	}

	pageInfo := &graphql.PageInfo{
		EndCursor:   connection.PageInfo.EndCursor,
		HasNextPage: connection.PageInfo.HasNextPage,
	}

	postConnection := &graphql.PostConnection{
		Edges:      edges,
		PageInfo:   pageInfo,
		TotalCount: connection.TotalCount,
	}

	return postConnection, nil
}
