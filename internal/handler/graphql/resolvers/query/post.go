package query

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/pkg/errors"
)

func (r *queryResolver) Post(ctx context.Context, id string) (*graphql.Post, error) {
	post, err := r.service.PostService.GetPostById(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get post by ID: %w")
	}

	return &graphql.Post{
		ID:            post.ID,
		Title:         post.Title,
		Body:          post.Body,
		Author:        post.Author,
		AllowComments: post.AllowComments,
		CreatedAt:     post.CreatedAt,
		// Comments: post.Comments,
	}, nil
}
