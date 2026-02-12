package comment

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type UseCase interface {
	AddComment(ctx context.Context, in models.AddCommentInput) (*models.Comment, error)
	GetComments(ctx context.Context, postID string, first *int32, after *string) (*models.CommentConnection, error)
}
