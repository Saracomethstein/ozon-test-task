package post

import (
	"context"
	"time"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (s *Post) CreatePost(ctx context.Context, in models.CreatePostInput) (*models.Post, error) {
	createAt := time.Now().UTC().Format(time.RFC3339)

	post, err := s.repo.Save(ctx, models.Post{
		Title:         in.Title,
		Author:        in.Author,
		Body:          in.Body,
		AllowComments: *in.AllowComments,
		CreatedAt:     createAt,
	})
	if err != nil {
		return nil, err
	}

	return &post, nil
}
