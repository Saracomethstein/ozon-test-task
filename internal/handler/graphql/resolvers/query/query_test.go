package query

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
	mockPost "github.com/Saracomethstein/ozon-test-task/internal/service/post/mocks"
)

func TestQueryResolver_Posts(t *testing.T) {
	t.Parallel()

	first := int32(2)
	after := "cursor123"
	totalCount := int32(5)
	hasNextPage := true
	endCursor := "nextCursor"

	post1 := &models.Post{
		ID:            1,
		Title:         "Post 1",
		Body:          "Body 1",
		Author:        "Author 1",
		AllowComments: true,
		CreatedAt:     "2023-01-01T12:00:00Z",
	}
	post2 := &models.Post{
		ID:            2,
		Title:         "Post 2",
		Body:          "Body 2",
		Author:        "Author 2",
		AllowComments: false,
		CreatedAt:     "2023-01-02T12:00:00Z",
	}

	serviceConnection := &models.PostConnection{
		Edges: []*models.PostEdge{
			{Cursor: "cursor1", Node: post1},
			{Cursor: "cursor2", Node: post2},
		},
		PageInfo: &models.PageInfo{
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
		TotalCount: totalCount,
	}

	tests := []struct {
		name        string
		first       *int32
		after       *string
		mockSetup   func(mockSvc *mockPost.MockUseCase)
		expected    *graphql.PostConnection
		expectedErr string
	}{
		{
			name:  "success_without_pagination",
			first: nil,
			after: nil,
			mockSetup: func(mockSvc *mockPost.MockUseCase) {
				mockSvc.EXPECT().
					GetPosts(mock.Anything, (*int32)(nil), (*string)(nil)).
					Return(serviceConnection, nil)
			},
			expected: &graphql.PostConnection{
				Edges: []*graphql.PostEdge{
					{
						Cursor: "cursor1",
						Node: &graphql.Post{
							ID:            "1",
							Title:         post1.Title,
							Body:          post1.Body,
							Author:        post1.Author,
							AllowComments: post1.AllowComments,
							CreatedAt:     post1.CreatedAt,
						},
					},
					{
						Cursor: "cursor2",
						Node: &graphql.Post{
							ID:            "2",
							Title:         post2.Title,
							Body:          post2.Body,
							Author:        post2.Author,
							AllowComments: post2.AllowComments,
							CreatedAt:     post2.CreatedAt,
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
			name:  "success_with_pagination",
			first: &first,
			after: &after,
			mockSetup: func(mockSvc *mockPost.MockUseCase) {
				mockSvc.EXPECT().
					GetPosts(mock.Anything, &first, &after).
					Return(serviceConnection, nil)
			},
			expected: &graphql.PostConnection{
				Edges: []*graphql.PostEdge{
					{
						Cursor: "cursor1",
						Node: &graphql.Post{
							ID:            "1",
							Title:         post1.Title,
							Body:          post1.Body,
							Author:        post1.Author,
							AllowComments: post1.AllowComments,
							CreatedAt:     post1.CreatedAt,
						},
					},
					{
						Cursor: "cursor2",
						Node: &graphql.Post{
							ID:            "2",
							Title:         post2.Title,
							Body:          post2.Body,
							Author:        post2.Author,
							AllowComments: post2.AllowComments,
							CreatedAt:     post2.CreatedAt,
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
			name:  "service_error",
			first: nil,
			after: nil,
			mockSetup: func(mockSvc *mockPost.MockUseCase) {
				mockSvc.EXPECT().
					GetPosts(mock.Anything, (*int32)(nil), (*string)(nil)).
					Return(nil, errors.New("db error"))
			},
			expectedErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockPostService := mockPost.NewMockUseCase(t)
			resolver := &queryResolver{
				service: &service.Container{
					PostService: mockPostService,
				},
			}

			tt.mockSetup(mockPostService)

			got, err := resolver.Posts(context.Background(), tt.first, tt.after)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			mockPostService.AssertExpectations(t)
		})
	}
}

func TestQueryResolver_Post(t *testing.T) {
	t.Parallel()

	postID := "123"
	postIDInt := int64(123)
	post := &models.Post{
		ID:            postIDInt,
		Title:         "Test Post",
		Body:          "Test Body",
		Author:        "John Doe",
		AllowComments: true,
		CreatedAt:     "2023-01-01T12:00:00Z",
	}

	tests := []struct {
		name        string
		id          string
		mockSetup   func(mockSvc *mockPost.MockUseCase)
		expected    *graphql.Post
		expectedErr string
	}{
		{
			name: "success",
			id:   postID,
			mockSetup: func(mockSvc *mockPost.MockUseCase) {
				mockSvc.EXPECT().
					GetPostById(mock.Anything, postID).
					Return(post, nil)
			},
			expected: &graphql.Post{
				ID:            postID,
				Title:         post.Title,
				Body:          post.Body,
				Author:        post.Author,
				AllowComments: post.AllowComments,
				CreatedAt:     post.CreatedAt,
			},
		},
		{
			name: "service_error",
			id:   postID,
			mockSetup: func(mockSvc *mockPost.MockUseCase) {
				mockSvc.EXPECT().
					GetPostById(mock.Anything, postID).
					Return(nil, errors.New("not found"))
			},
			expectedErr: "failed to get post by ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockPostService := mockPost.NewMockUseCase(t)
			resolver := &queryResolver{
				service: &service.Container{
					PostService: mockPostService,
				},
			}

			tt.mockSetup(mockPostService)

			got, err := resolver.Post(context.Background(), tt.id)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}

			mockPostService.AssertExpectations(t)
		})
	}
}

func TestQueryResolver_CommentsByPost(t *testing.T) {
	t.Parallel()

	postID := "123"
	postIDInt := int64(123)
	first := int32(2)
	after := "cursorX"
	totalCount := int32(3)
	hasNextPage := true
	endCursor := "nextCursor"

	comment1 := &models.Comment{
		ID:        1,
		PostID:    postIDInt,
		ParentID:  nil,
		Author:    "Alice",
		Text:      "Comment 1",
		CreatedAt: "2023-01-01T10:00:00Z",
	}
	comment2 := &models.Comment{
		ID:        2,
		PostID:    postIDInt,
		ParentID:  nil,
		Author:    "Bob",
		Text:      "Comment 2",
		CreatedAt: "2023-01-01T11:00:00Z",
	}

	serviceConnection := &models.CommentConnection{
		Edges: []*models.CommentEdge{
			{Cursor: "cursor1", Node: comment1},
			{Cursor: "cursor2", Node: comment2},
		},
		PageInfo: &models.PageInfo{
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
		TotalCount: totalCount,
	}

	tests := []struct {
		name        string
		postID      string
		first       *int32
		after       *string
		mockSetup   func(mockSvc *mockComment.MockUseCase)
		expected    *graphql.CommentConnection
		expectedErr string
	}{
		{
			name:   "success_without_pagination",
			postID: postID,
			first:  nil,
			after:  nil,
			mockSetup: func(mockSvc *mockComment.MockUseCase) {
				mockSvc.EXPECT().
					GetRootComments(mock.Anything, postID, (*int32)(nil), (*string)(nil)).
					Return(serviceConnection, nil)
			},
			expected: &graphql.CommentConnection{
				Edges: []*graphql.CommentEdge{
					{
						Cursor: "cursor1",
						Node: &graphql.Comment{
							ID:        "1",
							PostID:    postID,
							ParentID:  nil,
							Author:    comment1.Author,
							Text:      comment1.Text,
							CreatedAt: comment1.CreatedAt,
						},
					},
					{
						Cursor: "cursor2",
						Node: &graphql.Comment{
							ID:        "2",
							PostID:    postID,
							ParentID:  nil,
							Author:    comment2.Author,
							Text:      comment2.Text,
							CreatedAt: comment2.CreatedAt,
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
			name:   "success_with_pagination",
			postID: postID,
			first:  &first,
			after:  &after,
			mockSetup: func(mockSvc *mockComment.MockUseCase) {
				mockSvc.EXPECT().
					GetRootComments(mock.Anything, postID, &first, &after).
					Return(serviceConnection, nil)
			},
			expected: &graphql.CommentConnection{
				Edges: []*graphql.CommentEdge{
					{
						Cursor: "cursor1",
						Node: &graphql.Comment{
							ID:        "1",
							PostID:    postID,
							ParentID:  nil,
							Author:    comment1.Author,
							Text:      comment1.Text,
							CreatedAt: comment1.CreatedAt,
						},
					},
					{
						Cursor: "cursor2",
						Node: &graphql.Comment{
							ID:        "2",
							PostID:    postID,
							ParentID:  nil,
							Author:    comment2.Author,
							Text:      comment2.Text,
							CreatedAt: comment2.CreatedAt,
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
			name:        "validation_error",
			postID:      "",
			first:       nil,
			after:       nil,
			mockSetup:   func(mockSvc *mockComment.MockUseCase) {},
			expectedErr: "postID cannot be empty",
		},
		{
			name:   "service_error",
			postID: postID,
			first:  nil,
			after:  nil,
			mockSetup: func(mockSvc *mockComment.MockUseCase) {
				mockSvc.EXPECT().
					GetRootComments(mock.Anything, postID, (*int32)(nil), (*string)(nil)).
					Return(nil, errors.New("db error"))
			},
			expectedErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCommentService := mockComment.NewMockUseCase(t)
			resolver := &queryResolver{
				service: &service.Container{
					CommentService: mockCommentService,
				},
			}

			tt.mockSetup(mockCommentService)

			got, err := resolver.CommentsByPost(context.Background(), tt.postID, tt.first, tt.after)

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
