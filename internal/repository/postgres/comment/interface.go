package comment

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type UseCase interface {
	AddComment(ctx context.Context, comment models.Comment) (*models.Comment, error)
	CheckPostAllowComments(ctx context.Context, postID int64) (bool, error)
	CheckParentCommentExists(ctx context.Context, parentID int64) (int64, error)
	GetCommentPath(ctx context.Context, id int64) (string, error)
	SetCommentPath(ctx context.Context, id int64, path string) error

	GetCommentsByPost(ctx context.Context, postID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error)
	TotalCountComments(ctx context.Context, postID int64) (int64, error)
}
