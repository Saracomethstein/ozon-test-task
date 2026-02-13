package mutation

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
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
		ID:            strconv.FormatInt(out.ID, 10),
		Title:         out.Title,
		Body:          out.Body,
		Author:        out.Author,
		AllowComments: out.AllowComments,
		CreatedAt:     out.CreatedAt,
	}, nil
}
