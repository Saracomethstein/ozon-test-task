package post

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

func (s *postService) GetPosts(ctx context.Context, first *int32, after *string) (*models.PostConnection, error) {
	limit := s.getLimit(first)

	cursorPos := s.parseCursor(after)
	if cursorPos.err != nil {
		return nil, cursorPos.err
	}

	posts, err := s.repo.DB.Post.GetPosts(ctx, cursorPos.afterCreatedAt, cursorPos.afterID, limit+1)
	if err != nil {
		return nil, err
	}

	hasNextPage, pagePosts := s.extractPage(posts, limit)

	edges := s.buildEdges(pagePosts)

	totalCount, err := s.repo.DB.Post.TotalCount(ctx)
	if err != nil {
		return nil, err
	}

	return s.buildConnection(edges, hasNextPage, totalCount), nil
}

func (s *postService) getLimit(first *int32) int32 {
	if first != nil && *first > 0 {
		return *first
	}

	return defaultPageLimit
}

func (s *postService) parseCursor(after *string) cursorPosition {
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

func (s *postService) extractPage(posts []*models.Post, limit int32) (hasNextPage bool, page []*models.Post) {
	if len(posts) > int(limit) {
		return true, posts[:limit]
	}

	return false, posts
}

func (s *postService) buildEdges(posts []*models.Post) []*models.PostEdge {
	edges := make([]*models.PostEdge, 0, len(posts))

	for _, post := range posts {
		edge := &models.PostEdge{
			Cursor: s.generateCursor(post),
			Node:   post,
		}
		edges = append(edges, edge)
	}

	return edges
}

// rm after chage internal models for post
func (s *postService) generateCursor(post *models.Post) string {
	postID, _ := strconv.ParseInt(post.ID, 10, 64)
	return cursor.Encode(post.CreatedAt, postID)
}

func (s *postService) getEndCursor(edges []*models.PostEdge) *string {
	if len(edges) == 0 {
		return nil
	}

	return &edges[len(edges)-1].Cursor
}

func (s *postService) buildConnection(edges []*models.PostEdge, hasNextPage bool, totalCount int64) *models.PostConnection {
	return &models.PostConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			EndCursor:   s.getEndCursor(edges),
			HasNextPage: hasNextPage,
		},
		TotalCount: int32(totalCount),
	}
}
