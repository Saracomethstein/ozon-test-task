package comment

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type UseCase interface {
	// TODO: rename func for package 'comment'

	AddComment(ctx context.Context, in models.AddCommentInput) (*models.Comment, error)
	GetRootComments(ctx context.Context, postID string, first *int32, after *string) (*models.CommentConnection, error)
	GetChildComments(ctx context.Context, parentID string, first *int32, after *string) (*models.CommentConnection, error)
	Children(ctx context.Context, parentID int64, first *int32, after *string) (*models.CommentConnection, error)
}
