package comment

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (r *comment) Add(ctx context.Context, comment models.Comment) (*models.Comment, error) {
	return nil, nil
}
func (r *comment) CheckAllowComments(ctx context.Context, postID int64) (bool, error) {
	return false, nil
}
func (r *comment) CheckParentExists(ctx context.Context, parentID int64) (int64, error) {
	return 0, nil
}
func (r *comment) GetRootByPost(ctx context.Context, postID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error) {
	return nil, nil
}
func (r *comment) TotalCount(ctx context.Context, postID int64) (int64, error) {
	return 0, nil
}
func (r *comment) GetChild(ctx context.Context, parentID int64, afterCreatedAt *string, afterID int64, limit int32) ([]*models.Comment, error) {
	return nil, nil
}
func (r *comment) GetChildBatch(ctx context.Context, parentIDs []int64) ([]*models.Comment, error) {
	return nil, nil
}
