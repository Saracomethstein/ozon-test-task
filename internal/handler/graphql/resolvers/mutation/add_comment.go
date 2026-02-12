package mutation

import (
	"context"
	"strconv"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/pkg/errors"
)

func (r *mutationResolver) AddComment(ctx context.Context, input graphql.AddCommentInput) (*graphql.Comment, error) {
	if input.Author == "" || input.Text == "" || input.PostID == "" {
		return nil, errors.New("autor, text or postID cannot be empty")
	}

	comment, err := r.service.CommentService.AddComment(ctx, models.AddCommentInput{
		PostID:   input.PostID,
		ParentID: input.ParentID,
		Author:   input.Author,
		Text:     input.Text,
	})
	if err != nil {
		return nil, err
	}

	var parentIDPtr *string
	if comment.ParentID != nil {
		pid := strconv.FormatInt(*comment.ParentID, 10)
		parentIDPtr = &pid
	}

	gqlComment := &graphql.Comment{
		ID:        strconv.FormatInt(comment.ID, 10),
		PostID:    strconv.FormatInt(comment.PostID, 10),
		ParentID:  parentIDPtr,
		Author:    comment.Author,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
	}

	return gqlComment, nil
}
