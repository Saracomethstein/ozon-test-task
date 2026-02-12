package mutation

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/pkg/errors"
)

func (r *mutationResolver) SetPostCommentsAllowed(ctx context.Context, postID string, allow bool) (*graphql.Post, error) {
	if postID == "" {
		return nil, errors.New("postId cannot be empty")
	}

	out, err := r.service.PostService.SetPostCommentsAllowed(ctx, postID, allow)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set post comments allowed")
	}

	return &graphql.Post{
		ID:            out.ID,
		Title:         out.Title,
		Body:          out.Body,
		Author:        out.Author,
		AllowComments: out.AllowComments,
		CreatedAt:     out.CreatedAt,
		// Comments: out.Comments,
	}, nil
}
