package repository

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type CommentUC interface {
	Add(ctx context.Context, comment models.Comment) (*models.Comment, error)
	CheckAllowComments(ctx context.Context, postID int64) (bool, error)
	CheckParentExists(ctx context.Context, parentID int64) (int64, error)
	GetRootByPost(ctx context.Context, postID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error)
	TotalCount(ctx context.Context, postID int64) (int64, error)
	GetChild(ctx context.Context, parentID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error)
	GetChildBatch(ctx context.Context, parentIDs []int64) ([]*models.Comment, error)
}

type PostUC interface {
	Save(ctx context.Context, post models.Post) (models.Post, error)
	SetCommentsAllowed(ctx context.Context, postID int64, allow bool) (*models.Post, error)
	GetByID(ctx context.Context, postID int64) (*models.Post, error)
	Get(ctx context.Context, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Post, error)
	TotalCount(ctx context.Context) (int64, error)
}
