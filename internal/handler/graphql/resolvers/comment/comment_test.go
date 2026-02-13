package comment

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/service"
	mockComment "github.com/Saracomethstein/ozon-test-task/internal/service/comment/mocks"
)

func TestCommentResolver_Children(t *testing.T) {
	t.Parallel()

	commentID := "123"
	commentIDInt := int64(123)
	first := int32(2)
	after := "cursor123"
	totalCount := int32(3)
	hasNextPage := true
	endCursor := "nextCursor"

	child1 := &models.Comment{
		ID:        4,
		PostID:    123,
		ParentID:  &commentIDInt,
		Author:    "Child1",
		Text:      "Child comment 1",
		CreatedAt: "2023-01-01T12:00:00Z",
	}
	child2 := &models.Comment{
		ID:        5,
		PostID:    123,
		ParentID:  &commentIDInt,
		Author:    "Child2",
		Text:      "Child comment 2",
		CreatedAt: "2023-01-01T13:00:00Z",
	}

	serviceConnection := &models.CommentConnection{
		Edges: []*models.CommentEdge{
			{Cursor: "cursor1", Node: child1},
			{Cursor: "cursor2", Node: child2},
		},
		PageInfo: &models.PageInfo{
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
		TotalCount: totalCount,
	}

	tests := []struct {
		name        string
		obj         *graphql.Comment
		first       *int32
		after       *string
		mockSetup   func(mockSvc *mockComment.MockUseCase)
		expected    *graphql.CommentConnection
		expectedErr string
	}{
		{
			name: "success_without_pagination",
			obj: &graphql.Comment{
				ID: commentID,
			},
			first: nil,
			after: nil,
			mockSetup: func(mockSvc *mockComment.MockUseCase) {
				mockSvc.EXPECT().
					Children(mock.Anything, commentIDInt, (*int32)(nil), (*string)(nil)).
					Return(serviceConnection, nil)
			},
			expected: &graphql.CommentConnection{
				Edges: []*graphql.CommentEdge{
					{
						Cursor: "cursor1",
						Node: &graphql.Comment{
							ID:        "4",
							PostID:    "123",
							ParentID:  &commentID,
							Author:    child1.Author,
							Text:      child1.Text,
							CreatedAt: child1.CreatedAt,
						},
					},
					{
						Cursor: "cursor2",
						Node: &graphql.Comment{
							ID:        "5",
							PostID:    "123",
							ParentID:  &commentID,
							Author:    child2.Author,
							Text:      child2.Text,
							CreatedAt: child2.CreatedAt,
						},
					},
				},
				PageInfo: &graphql.PageInfo{
					EndCursor:   &endCursor,
					HasNextPage: hasNextPage,
				},
				TotalCount: totalCount,
			},
		},
		{
			name: "success_with_pagination",
			obj: &graphql.Comment{
				ID: commentID,
			},
			first: &first,
			after: &after,
			mockSetup: func(mockSvc *mockComment.MockUseCase) {
				mockSvc.EXPECT().
					Children(mock.Anything, commentIDInt, &first, &after).
					Return(serviceConnection, nil)
			},
			expected: &graphql.CommentConnection{
				Edges: []*graphql.CommentEdge{
					{
						Cursor: "cursor1",
						Node: &graphql.Comment{
							ID:        "4",
							PostID:    "123",
							ParentID:  &commentID,
							Author:    child1.Author,
							Text:      child1.Text,
							CreatedAt: child1.CreatedAt,
						},
					},
					{
						Cursor: "cursor2",
						Node: &graphql.Comment{
							ID:        "5",
							PostID:    "123",
							ParentID:  &commentID,
							Author:    child2.Author,
							Text:      child2.Text,
							CreatedAt: child2.CreatedAt,
						},
					},
				},
				PageInfo: &graphql.PageInfo{
					EndCursor:   &endCursor,
					HasNextPage: hasNextPage,
				},
				TotalCount: totalCount,
			},
		},
		{
			name: "service_error",
			obj: &graphql.Comment{
				ID: commentID,
			},
			first: nil,
			after: nil,
			mockSetup: func(mockSvc *mockComment.MockUseCase) {
				mockSvc.EXPECT().
					Children(mock.Anything, commentIDInt, (*int32)(nil), (*string)(nil)).
					Return(nil, errors.New("db error"))
			},
			expectedErr: "db error",
		},
		{
			name: "invalid_comment_ID",
			obj: &graphql.Comment{
				ID: "abc",
			},
			first:       nil,
			after:       nil,
			mockSetup:   func(mockSvc *mockComment.MockUseCase) {},
			expectedErr: "invalid comment ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCommentService := mockComment.NewMockUseCase(t)
			resolver := &commentResolver{
				service: &service.Container{
					CommentService: mockCommentService,
				},
			}

			tt.mockSetup(mockCommentService)

			got, err := resolver.Children(context.Background(), tt.obj, tt.first, tt.after)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			mockCommentService.AssertExpectations(t)
		})
	}
}
