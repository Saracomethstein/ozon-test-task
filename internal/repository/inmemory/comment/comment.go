package comment

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *comment) AddComment(ctx context.Context, comment models.Comment) (*models.Comment, error) {
	return nil, nil
}
func (r *comment) CheckPostAllowComments(ctx context.Context, postID int64) (bool, error) {
	return false, nil
}
func (r *comment) CheckParentCommentExists(ctx context.Context, parentID int64) (int64, error) {
	return 0, nil
}
func (r *comment) GetRootCommentsByPost(ctx context.Context, postID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error) {
	return nil, nil
}
func (r *comment) TotalCountComments(ctx context.Context, postID int64) (int64, error) {
	return 0, nil
}
func (r *comment) GetChildComments(ctx context.Context, parentID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error) {
	return nil, nil
}
func (r *comment) GetChildCommentsBatch(ctx context.Context, parentIDs []int64) ([]*models.Comment, error) {
	return nil, nil
}
