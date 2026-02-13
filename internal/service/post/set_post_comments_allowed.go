package post

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (s *Post) SetPostCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error) {
	out, err := s.repo.SetPostCommentsAllowed(ctx, postID, allow)
	if err != nil {
		return nil, err
	}
	return out, nil
}
