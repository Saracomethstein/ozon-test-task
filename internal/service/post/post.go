package post

import (
	"context"
	"strconv"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (s *postService) GetPostById(ctx context.Context, postID string) (*models.Post, error) {
	if postID == "" {
		return nil, errors.New("post ID cannot be empty")
	}

	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid post ID format")
	}

	if id <= 0 {
		return nil, errors.New("post ID must be a positive integer")
	}

	post, err := s.repo.DB.Post.GetPostById(ctx, id)
	if err != nil {
		return nil, err
	}

	return post, nil
}
