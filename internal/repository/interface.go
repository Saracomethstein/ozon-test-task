package repository

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type UseCase interface {
	SavePost(ctx context.Context, post models.Post) (models.Post, error)
	SetPostCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error)
	GetPostById(ctx context.Context, postID int64) (*models.Post, error)
	GetPosts(ctx context.Context, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Post, error)
	TotalCount(ctx context.Context) (int64, error)

	AddComment(ctx context.Context, comment models.Comment) (*models.Comment, error)
	CheckPostAllowComments(ctx context.Context, postID int64) (bool, error)
	CheckParentCommentExists(ctx context.Context, parentID int64) (int64, error)
	GetCommentPath(ctx context.Context, id int64) (string, error)
	SetCommentPath(ctx context.Context, id int64, path string) error
}
