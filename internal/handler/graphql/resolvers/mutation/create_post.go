package mutation

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input graphql.CreatePostInput) (*graphql.Post, error) {
	// create validator for inupt data
	if input.Title == "" || input.Author == "" || input.Body == "" {
		return nil, errors.New("title, author and body are required fields")
	}

	out, err := r.service.PostService.CreatePost(ctx, models.CreatePostInput{
		Title:         input.Title,
		Body:          input.Body,
		Author:        input.Author,
		AllowComments: input.AllowComments,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create post")
	}

	return &graphql.Post{
		// create converter models.Post -> graphql.Post
		ID: out.ID,
	}, nil
}
