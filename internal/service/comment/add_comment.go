package comment

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

const (
	pathIDWidth = 10
)

type parentInfo struct {
	id   *int64
	path string
}

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

	p, err := s.processParent(ctx, in.ParentID, postID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	comment, err := s.repo.DB.Comment.AddComment(ctx, models.Comment{
		PostID:    postID,
		ParentID:  p.id,
		Author:    in.Author,
		Text:      in.Text,
		CreatedAt: now,
	})
	if err != nil {
		return nil, err
	}

	formattedID := formatIDForPath(comment.ID)

	if p.id == nil {
		comment.Path = formattedID
	} else {
		comment.Path = p.path + "." + formattedID
	}

	if err := s.repo.DB.Comment.SetCommentPath(ctx, comment.ID, comment.Path); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) processParent(ctx context.Context, parentIDStr *string, postID int64) (*parentInfo, error) {
	if parentIDStr == nil || *parentIDStr == "" {
		return &parentInfo{id: nil, path: ""}, nil
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

	path, err := s.repo.DB.Comment.GetCommentPath(ctx, pid)
	if err != nil {
		return nil, err
	}

	return &parentInfo{
		id:   &pid,
		path: path,
	}, nil
}

func formatIDForPath(id int64) string {
	return fmt.Sprintf("%0*d", pathIDWidth, id)
}
