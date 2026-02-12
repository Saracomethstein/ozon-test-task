package comment

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type UseCase interface {
	// TODO: rename func for package 'comment'

	AddComment(ctx context.Context, comment models.Comment) (*models.Comment, error)
	CheckPostAllowComments(ctx context.Context, postID int64) (bool, error)
	CheckParentCommentExists(ctx context.Context, parentID int64) (int64, error)

	GetRootCommentsByPost(ctx context.Context, postID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error)
	TotalCountComments(ctx context.Context, postID int64) (int64, error)
	GetChildComments(ctx context.Context, parentID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error)
	GetChildCommentsBatch(ctx context.Context, parentIDs []int64) ([]*models.Comment, error)
}
