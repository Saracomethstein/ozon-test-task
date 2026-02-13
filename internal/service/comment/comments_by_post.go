package comment

import (
	"context"
	"strconv"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/utils/cursor"
	"github.com/pkg/errors"
)

const (
	defaultPageLimit = 20
)

type cursorPosition struct {
	afterCreatedAt *string
	afterID        int64
	err            error
}

func (s *Service) GetRootComments(ctx context.Context, postID string, first *int32, after *string) (*models.CommentConnection, error) {
	pID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid postID format")
	}
	if pID <= 0 {
		return nil, errors.New("postID must be greater 0")
	}

	limit := s.getLimit(first)

	cursorPos := s.parseCursor(after)
	if cursorPos.err != nil {
		return nil, cursorPos.err
	}

	comments, err := s.repo.GetRootCommentsByPost(ctx, pID, cursorPos.afterCreatedAt, cursorPos.afterID, limit+1)
	if err != nil {
		return nil, err
	}

	hasNextPage, pageComments := s.extractPage(comments, limit)

	edges := s.buildEdges(pageComments)

	totalCount, err := s.repo.TotalCountComments(ctx, pID)
	if err != nil {
		return nil, err
	}

	return s.buildConnection(edges, hasNextPage, totalCount), nil
}

func (s *Service) GetChildComments(ctx context.Context, parentID string, first *int32, after *string) (*models.CommentConnection, error) {
	pID, err := strconv.ParseInt(parentID, 10, 64)
	if err != nil || pID <= 0 {
		return nil, errors.New("invalid parentID")
	}
	if pID <= 0 {
		return nil, errors.New("postID must be greater 0")
	}

	limit := s.getLimit(first)

	cursorPos := s.parseCursor(after)
	if cursorPos.err != nil {
		return nil, cursorPos.err
	}

	comments, err := s.repo.GetChildComments(ctx, pID, cursorPos.afterCreatedAt, cursorPos.afterID, limit+1)
	if err != nil {
		return nil, err
	}

	hasNextPage, pageComments := s.extractPage(comments, limit)

	edges := s.buildEdges(pageComments)

	totalCount := int64(len(pageComments))

	return s.buildConnection(edges, hasNextPage, totalCount), nil
}

func (s *Service) getLimit(first *int32) int32 {
	if first != nil && *first > 0 {
		return *first
	}

	return defaultPageLimit
}

func (s *Service) parseCursor(after *string) cursorPosition {
	if after == nil || *after == "" {
		return cursorPosition{}
	}

	createdAt, id, err := cursor.Decode(*after)
	if err != nil {
		return cursorPosition{
			err: errors.New("invalid cursor format"),
		}
	}

	return cursorPosition{
		afterCreatedAt: &createdAt,
		afterID:        id,
	}
}

func (s *Service) extractPage(comments []*models.Comment, limit int32) (hasNextPage bool, page []*models.Comment) {
	if len(comments) > int(limit) {
		return true, comments[:limit]
	}

	return false, comments
}

func (s *Service) buildEdges(comments []*models.Comment) []*models.CommentEdge {
	edges := make([]*models.CommentEdge, 0, len(comments))

	for _, c := range comments {
		edge := &models.CommentEdge{
			Cursor: cursor.Encode(c.CreatedAt, c.ID),
			Node:   c,
		}
		edges = append(edges, edge)
	}

	return edges
}

func (s *Service) getEndCursor(edges []*models.CommentEdge) *string {
	if len(edges) == 0 {
		return nil
	}

	return &edges[len(edges)-1].Cursor
}

func (s *Service) buildConnection(edges []*models.CommentEdge, hasNextPage bool, totalCount int64) *models.CommentConnection {
	return &models.CommentConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			EndCursor:   s.getEndCursor(edges),
			HasNextPage: hasNextPage,
		},
		TotalCount: int32(totalCount),
	}
}
