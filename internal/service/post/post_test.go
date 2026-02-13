package post

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/repository/mocks"
	"github.com/Saracomethstein/ozon-test-task/internal/utils/cursor"
)

func TestPostService_CreatePost(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	allowComments := true

	tests := []struct {
		name        string
		input       models.CreatePostInput
		setupMock   func(repo *mocks.MockPostUC)
		want        *models.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name: "successful_creation",
			input: models.CreatePostInput{
				Title:         "Test Title",
				Author:        "Test Author",
				Body:          "Test Body",
				AllowComments: &allowComments,
			},
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("Save", mock.Anything, mock.MatchedBy(func(p models.Post) bool {
					return p.Title == "Test Title" &&
						p.Author == "Test Author" &&
						p.Body == "Test Body" &&
						p.AllowComments == true &&
						p.CreatedAt != ""
				})).Return(models.Post{
					ID:            1,
					Title:         "Test Title",
					Author:        "Test Author",
					Body:          "Test Body",
					AllowComments: true,
					CreatedAt:     "2023-01-01T00:00:00Z",
				}, nil)
			},
			want: &models.Post{
				ID:            1,
				Title:         "Test Title",
				Author:        "Test Author",
				Body:          "Test Body",
				AllowComments: true,
				CreatedAt:     "2023-01-01T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name: "repository_error",
			input: models.CreatePostInput{
				Title:         "Test Title",
				Author:        "Test Author",
				Body:          "Test Body",
				AllowComments: &allowComments,
			},
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("Save", mock.Anything, mock.Anything).Return(models.Post{}, errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockPostUC(t)
			tt.setupMock(mockRepo)

			s := New(mockRepo)
			got, err := s.CreatePost(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_GetPostById(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name        string
		postID      string
		setupMock   func(repo *mocks.MockPostUC)
		want        *models.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "successful_get",
			postID: "42",
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("GetByID", mock.Anything, int64(42)).Return(&models.Post{
					ID:            42,
					Title:         "Title",
					Author:        "Author",
					Body:          "Body",
					AllowComments: true,
					CreatedAt:     "2023-01-01T00:00:00Z",
				}, nil)
			},
			want: &models.Post{
				ID:            42,
				Title:         "Title",
				Author:        "Author",
				Body:          "Body",
				AllowComments: true,
				CreatedAt:     "2023-01-01T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name:        "empty_id",
			postID:      "",
			setupMock:   func(repo *mocks.MockPostUC) {},
			wantErr:     true,
			expectedErr: errors.New("post ID cannot be empty"),
		},
		{
			name:        "invalid_id_format",
			postID:      "abc",
			setupMock:   func(repo *mocks.MockPostUC) {},
			wantErr:     true,
			expectedErr: errors.New("invalid post ID format"),
		},
		{
			name:        "negative_id",
			postID:      "-5",
			setupMock:   func(repo *mocks.MockPostUC) {},
			wantErr:     true,
			expectedErr: errors.New("post ID must be a positive integer"),
		},
		{
			name:   "post_not_found",
			postID: "999",
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("post not found"))
			},
			wantErr:     true,
			expectedErr: errors.New("post not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockPostUC(t)
			tt.setupMock(mockRepo)

			s := New(mockRepo)
			got, err := s.GetPostById(ctx, tt.postID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_SetPostCommentsAllowed(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name        string
		postID      string
		allow       bool
		setupMock   func(repo *mocks.MockPostUC)
		want        *models.Post
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "successful_update",
			postID: "42",
			allow:  true,
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("SetCommentsAllowed", mock.Anything, int64(42), true).Return(&models.Post{
					ID:            42,
					Title:         "Title",
					Author:        "Author",
					Body:          "Body",
					AllowComments: true,
					CreatedAt:     "2023-01-01T00:00:00Z",
				}, nil)
			},
			want: &models.Post{
				ID:            42,
				Title:         "Title",
				Author:        "Author",
				Body:          "Body",
				AllowComments: true,
				CreatedAt:     "2023-01-01T00:00:00Z",
			},
			wantErr: false,
		},
		{
			name:        "empty_id",
			postID:      "",
			allow:       true,
			setupMock:   func(repo *mocks.MockPostUC) {},
			wantErr:     true,
			expectedErr: errors.New("post ID cannot be empty"),
		},
		{
			name:        "invalid_id_format",
			postID:      "abc",
			allow:       true,
			setupMock:   func(repo *mocks.MockPostUC) {},
			wantErr:     true,
			expectedErr: errors.New("invalid post ID format"),
		},
		{
			name:        "negative_id",
			postID:      "-5",
			allow:       true,
			setupMock:   func(repo *mocks.MockPostUC) {},
			wantErr:     true,
			expectedErr: errors.New("post ID must be a positive integer"),
		},
		{
			name:   "post_not_found",
			postID: "999",
			allow:  true,
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("SetCommentsAllowed", mock.Anything, int64(999), true).Return(nil, errors.New("post not found"))
			},
			wantErr:     true,
			expectedErr: errors.New("post not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockPostUC(t)
			tt.setupMock(mockRepo)

			s := New(mockRepo)
			got, err := s.SetPostCommentsAllowed(ctx, tt.postID, tt.allow)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPostService_GetPosts(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	post1 := &models.Post{
		ID:            1,
		Title:         "Post 1",
		Author:        "Author 1",
		Body:          "Body 1",
		AllowComments: true,
		CreatedAt:     "2023-01-01T12:00:00Z",
	}
	post2 := &models.Post{
		ID:            2,
		Title:         "Post 2",
		Author:        "Author 2",
		Body:          "Body 2",
		AllowComments: true,
		CreatedAt:     "2023-01-01T11:00:00Z",
	}
	post3 := &models.Post{
		ID:            3,
		Title:         "Post 3",
		Author:        "Author 3",
		Body:          "Body 3",
		AllowComments: true,
		CreatedAt:     "2023-01-01T10:00:00Z",
	}

	cursorFor := func(post *models.Post) string {
		return cursor.Encode(post.CreatedAt, post.ID)
	}

	tests := []struct {
		name          string
		first         *int32
		after         *string
		setupMock     func(repo *mocks.MockPostUC)
		expected      *models.PostConnection
		expectedError string
	}{
		{
			name:  "default_limit_without_cursor",
			first: nil,
			after: nil,
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("Get", mock.Anything, (*string)(nil), int64(0), int32(21)).
					Return([]*models.Post{post1, post2, post3}, nil)
				repo.On("TotalCount", mock.Anything).Return(int64(5), nil)
			},
			expected: &models.PostConnection{
				Edges: []*models.PostEdge{
					{Cursor: cursorFor(post1), Node: post1},
					{Cursor: cursorFor(post2), Node: post2},
					{Cursor: cursorFor(post3), Node: post3},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursorFor(post3)),
					HasNextPage: false,
				},
				TotalCount: 5,
			},
		},
		{
			name:  "default_limit_with_next_page",
			first: nil,
			after: nil,
			setupMock: func(repo *mocks.MockPostUC) {
				posts := make([]*models.Post, 22)
				for i := 0; i < 22; i++ {
					p := *post1
					p.ID = int64(i + 1)
					posts[i] = &p
				}
				repo.On("Get", mock.Anything, (*string)(nil), int64(0), int32(21)).
					Return(posts, nil)
				repo.On("TotalCount", mock.Anything).Return(int64(30), nil)
			},
			expected: &models.PostConnection{
				Edges: func() []*models.PostEdge {
					edges := make([]*models.PostEdge, 20)
					for i := 0; i < 20; i++ {
						p := *post1
						p.ID = int64(i + 1)
						edges[i] = &models.PostEdge{
							Cursor: cursor.Encode(p.CreatedAt, p.ID),
							Node:   &p,
						}
					}
					return edges
				}(),
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursor.Encode(post1.CreatedAt, 20)),
					HasNextPage: true,
				},
				TotalCount: 30,
			},
		},
		{
			name:  "custom_first_smaller_than_available",
			first: int32Ptr(2),
			after: nil,
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("Get", mock.Anything, (*string)(nil), int64(0), int32(3)).
					Return([]*models.Post{post1, post2}, nil)
				repo.On("TotalCount", mock.Anything).Return(int64(10), nil)
			},
			expected: &models.PostConnection{
				Edges: []*models.PostEdge{
					{Cursor: cursorFor(post1), Node: post1},
					{Cursor: cursorFor(post2), Node: post2},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursorFor(post2)),
					HasNextPage: false,
				},
				TotalCount: 10,
			},
		},
		{
			name:  "custom_first_with_next_page",
			first: int32Ptr(2),
			after: nil,
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("Get", mock.Anything, (*string)(nil), int64(0), int32(3)).
					Return([]*models.Post{post1, post2, post3}, nil)
				repo.On("TotalCount", mock.Anything).Return(int64(10), nil)
			},
			expected: &models.PostConnection{
				Edges: []*models.PostEdge{
					{Cursor: cursorFor(post1), Node: post1},
					{Cursor: cursorFor(post2), Node: post2},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursorFor(post2)),
					HasNextPage: true,
				},
				TotalCount: 10,
			},
		},
		{
			name:  "with_after_cursor",
			first: int32Ptr(2),
			after: strPtr(cursor.Encode("2023-01-01T11:00:00Z", 2)),
			setupMock: func(repo *mocks.MockPostUC) {
				afterCreatedAt := "2023-01-01T11:00:00Z"
				afterID := int64(2)
				repo.On("Get", mock.Anything, &afterCreatedAt, afterID, int32(3)).
					Return([]*models.Post{post3}, nil)
				repo.On("TotalCount", mock.Anything).Return(int64(10), nil)
			},
			expected: &models.PostConnection{
				Edges: []*models.PostEdge{
					{Cursor: cursorFor(post3), Node: post3},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursorFor(post3)),
					HasNextPage: false,
				},
				TotalCount: 10,
			},
		},
		{
			name:          "invalid_cursor_format",
			first:         int32Ptr(2),
			after:         strPtr("invalid-cursor"),
			setupMock:     func(repo *mocks.MockPostUC) {},
			expectedError: "invalid cursor format",
		},
		{
			name:  "empty_after_string",
			first: int32Ptr(2),
			after: strPtr(""),
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("Get", mock.Anything, (*string)(nil), int64(0), int32(3)).
					Return([]*models.Post{post1, post2}, nil)
				repo.On("TotalCount", mock.Anything).Return(int64(10), nil)
			},
			expected: &models.PostConnection{
				Edges: []*models.PostEdge{
					{Cursor: cursorFor(post1), Node: post1},
					{Cursor: cursorFor(post2), Node: post2},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursorFor(post2)),
					HasNextPage: false,
				},
				TotalCount: 10,
			},
		},
		{
			name:  "repo_Get_error",
			first: int32Ptr(2),
			after: nil,
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("Get", mock.Anything, (*string)(nil), int64(0), int32(3)).
					Return(nil, errors.New("db error"))
			},
			expectedError: "db error",
		},
		{
			name:  "repo_TotalCount_error",
			first: int32Ptr(2),
			after: nil,
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("Get", mock.Anything, (*string)(nil), int64(0), int32(3)).
					Return([]*models.Post{post1, post2}, nil)
				repo.On("TotalCount", mock.Anything).Return(int64(0), errors.New("count error"))
			},
			expectedError: "count error",
		},
		{
			name:  "zero_results",
			first: int32Ptr(2),
			after: nil,
			setupMock: func(repo *mocks.MockPostUC) {
				repo.On("Get", mock.Anything, (*string)(nil), int64(0), int32(3)).
					Return([]*models.Post{}, nil)
				repo.On("TotalCount", mock.Anything).Return(int64(0), nil)
			},
			expected: &models.PostConnection{
				Edges:      []*models.PostEdge{},
				PageInfo:   &models.PageInfo{EndCursor: nil, HasNextPage: false},
				TotalCount: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockPostUC(t)
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			s := New(mockRepo)
			got, err := s.GetPosts(ctx, tt.first, tt.after)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, got)

				assert.Equal(t, tt.expected.TotalCount, got.TotalCount)
				assert.Equal(t, tt.expected.PageInfo.HasNextPage, got.PageInfo.HasNextPage)
				assert.Equal(t, tt.expected.PageInfo.EndCursor, got.PageInfo.EndCursor)

				require.Len(t, got.Edges, len(tt.expected.Edges))
				for i, expectedEdge := range tt.expected.Edges {
					assert.Equal(t, expectedEdge.Cursor, got.Edges[i].Cursor)
					assert.Equal(t, expectedEdge.Node, got.Edges[i].Node)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func int32Ptr(v int32) *int32 { return &v }
func strPtr(v string) *string { return &v }
