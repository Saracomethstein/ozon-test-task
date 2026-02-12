package comment

import (
	"context"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func (s *commentService) AddComment(ctx context.Context, in models.AddCommentInput) (*models.Comment, error) {
	postID, err := strconv.ParseInt(in.PostID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid postID format")
	}

	if postID <= 0 {
		return nil, errors.New("postID must be greater 0")
	}

	allow, err := s.repo.DB.Comment.CheckPostAllowComments(ctx, postID)
	if err != nil {
		return nil, err
	}
	if !allow {
		return nil, errors.New("comments not allowed for this post")
	}

	parentID, err := s.processParent(ctx, in.ParentID, postID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	comment, err := s.repo.DB.Comment.AddComment(ctx, models.Comment{
		PostID:    postID,
		ParentID:  parentID,
		Author:    in.Author,
		Text:      in.Text,
		CreatedAt: now,
	})
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) processParent(ctx context.Context, parentIDStr *string, postID int64) (*int64, error) {
	if parentIDStr == nil || *parentIDStr == "" {
		return nil, nil
	}

	pid, err := strconv.ParseInt(*parentIDStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid parentID format")
	}

	parentPostID, err := s.repo.DB.Comment.CheckParentCommentExists(ctx, pid)
	if err != nil {
		return nil, err
	}
	if parentPostID != postID {
		return nil, errors.New("parent comment does not belong to this post")
	}

	return &pid, nil
}
