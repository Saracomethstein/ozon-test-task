package comment

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

const RootParent = int64(0)

func (r *comment) Add(ctx context.Context, comment models.Comment) (*models.Comment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.seq++
	id := r.seq

	clone := models.Comment{
		ID:        id,
		PostID:    comment.PostID,
		ParentID:  comment.ParentID,
		Author:    comment.Author,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
	}
	r.comments[id] = &clone
	r.byPost[clone.PostID] = append(r.byPost[clone.PostID], id)

	parentKey := RootParent
	if clone.ParentID != nil {
		parentKey = *clone.ParentID
	}
	r.byParent[parentKey] = append(r.byParent[parentKey], id)

	return &clone, nil
}

func (r *comment) CheckAllowComments(ctx context.Context, postID int64) (bool, error) {
	post, err := r.repoPost.GetByID(ctx, postID)
	if err != nil {
		return false, err
	}

	return post.AllowComments, nil
}

func (r *comment) CheckParentExists(ctx context.Context, parentID int64) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parent, ok := r.comments[parentID]
	if !ok {
		return 0, errors.New("parent comment not found")
	}

	return parent.PostID, nil
}
