package comment

import (
	"context"
	"testing"
	"time"

	"github.com/graph-gophers/dataloader"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/repository/mocks"
	"github.com/Saracomethstein/ozon-test-task/internal/utils/cursor"
)

func TestService_AddComment(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	now := time.Now().UTC().Format(time.RFC3339)

	parentID := int64(10)
	parentIDStr := "10"
	postID := int64(1)
	postIDStr := "1"

	tests := []struct {
		name        string
		input       models.AddCommentInput
		setupMock   func(repo *mocks.MockCommentUC)
		want        *models.Comment
		wantErr     bool
		expectedErr string
	}{
		{
			name: "successful_root_comment",
			input: models.AddCommentInput{
				PostID:   postIDStr,
				ParentID: nil,
				Author:   "Alice",
				Text:     "Hello",
			},
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("CheckAllowComments", context.Background(), postID).Return(true, nil)
				repo.On("Add", mock.Anything, mock.MatchedBy(func(c models.Comment) bool {
					return c.PostID == postID && c.ParentID == nil && c.Author == "Alice" && c.Text == "Hello" && c.CreatedAt != ""
				})).Return(&models.Comment{
					ID:        1,
					PostID:    postID,
					ParentID:  nil,
					Author:    "Alice",
					Text:      "Hello",
					CreatedAt: now,
				}, nil)
			},
			want: &models.Comment{
				ID:        1,
				PostID:    postID,
				ParentID:  nil,
				Author:    "Alice",
				Text:      "Hello",
				CreatedAt: now,
			},
		},
		{
			name: "successful_child_comment",
			input: models.AddCommentInput{
				PostID:   postIDStr,
				ParentID: &parentIDStr,
				Author:   "Bob",
				Text:     "Reply",
			},
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("CheckAllowComments", mock.Anything, postID).Return(true, nil)
				repo.On("CheckParentExists", mock.Anything, parentID).Return(postID, nil)
				repo.On("Add", mock.Anything, mock.MatchedBy(func(c models.Comment) bool {
					return c.PostID == postID && *c.ParentID == parentID && c.Author == "Bob" && c.Text == "Reply"
				})).Return(&models.Comment{
					ID:        2,
					PostID:    postID,
					ParentID:  &parentID,
					Author:    "Bob",
					Text:      "Reply",
					CreatedAt: now,
				}, nil)
			},
			want: &models.Comment{
				ID:        2,
				PostID:    postID,
				ParentID:  &parentID,
				Author:    "Bob",
				Text:      "Reply",
				CreatedAt: now,
			},
		},
		{
			name: "invalid_postID_format",
			input: models.AddCommentInput{
				PostID:   "abc",
				ParentID: nil,
			},
			setupMock:   func(repo *mocks.MockCommentUC) {},
			wantErr:     true,
			expectedErr: "invalid postID format",
		},
		{
			name: "invalid_postID",
			input: models.AddCommentInput{
				PostID:   "0",
				ParentID: nil,
			},
			setupMock:   func(repo *mocks.MockCommentUC) {},
			wantErr:     true,
			expectedErr: "postID must be greater 0",
		},
		{
			name: "comments_not_allowed",
			input: models.AddCommentInput{
				PostID:   postIDStr,
				ParentID: nil,
			},
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("CheckAllowComments", mock.Anything, postID).Return(false, nil)
			},
			wantErr:     true,
			expectedErr: "comments not allowed for this post",
		},
		{
			name: "checkAllowComments_error",
			input: models.AddCommentInput{
				PostID:   postIDStr,
				ParentID: nil,
			},
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("CheckAllowComments", mock.Anything, postID).Return(false, errors.New("db error"))
			},
			wantErr:     true,
			expectedErr: "db error",
		},
		{
			name: "invalid_parentID_format",
			input: models.AddCommentInput{
				PostID:   postIDStr,
				ParentID: strPtr("abc"),
			},
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("CheckAllowComments", mock.Anything, postID).Return(true, nil)
			},
			wantErr:     true,
			expectedErr: "invalid parentID format",
		},
		{
			name: "parent_comment_not_found",
			input: models.AddCommentInput{
				PostID:   postIDStr,
				ParentID: &parentIDStr,
			},
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("CheckAllowComments", mock.Anything, postID).Return(true, nil)
				repo.On("CheckParentExists", mock.Anything, parentID).Return(int64(0), errors.New("parent comment not found"))
			},
			wantErr:     true,
			expectedErr: "parent comment not found",
		},
		{
			name: "parent_comment_belongs_to_different_post",
			input: models.AddCommentInput{
				PostID:   postIDStr,
				ParentID: &parentIDStr,
			},
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("CheckAllowComments", mock.Anything, postID).Return(true, nil)
				repo.On("CheckParentExists", mock.Anything, parentID).Return(int64(2), nil)
			},
			wantErr:     true,
			expectedErr: "parent comment does not belong to this post",
		},
		{
			name: "repo_Add_error",
			input: models.AddCommentInput{
				PostID:   postIDStr,
				ParentID: nil,
			},
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("CheckAllowComments", mock.Anything, postID).Return(true, nil)
				repo.On("Add", mock.Anything, mock.Anything).Return(nil, errors.New("insert failed"))
			},
			wantErr:     true,
			expectedErr: "insert failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockCommentUC(t)
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			s := New(mockRepo)
			got, err := s.AddComment(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetRootComments(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	postID := int64(1)
	postIDStr := "1"

	now := time.Now().UTC()
	comment1 := &models.Comment{ID: 1, PostID: postID, Author: "A", Text: "1", CreatedAt: now.Format(time.RFC3339)}
	comment2 := &models.Comment{ID: 2, PostID: postID, Author: "B", Text: "2", CreatedAt: now.Add(-time.Hour).Format(time.RFC3339)}
	comment3 := &models.Comment{ID: 3, PostID: postID, Author: "C", Text: "3", CreatedAt: now.Add(-2 * time.Hour).Format(time.RFC3339)}

	tests := []struct {
		name        string
		postID      string
		first       *int32
		after       *string
		setupMock   func(repo *mocks.MockCommentUC)
		want        *models.CommentConnection
		wantErr     bool
		expectedErr string
	}{
		{
			name:   "first_page_without_cursor",
			postID: postIDStr,
			first:  nil,
			after:  nil,
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("GetRootByPost", mock.Anything, postID, (*string)(nil), int64(0), int32(21)).
					Return([]*models.Comment{comment1, comment2, comment3}, nil)
				repo.On("TotalCount", mock.Anything, postID).Return(int64(10), nil)
			},
			want: &models.CommentConnection{
				Edges: []*models.CommentEdge{
					{Cursor: cursor.Encode(comment1.CreatedAt, comment1.ID), Node: comment1},
					{Cursor: cursor.Encode(comment2.CreatedAt, comment2.ID), Node: comment2},
					{Cursor: cursor.Encode(comment3.CreatedAt, comment3.ID), Node: comment3},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursor.Encode(comment3.CreatedAt, comment3.ID)),
					HasNextPage: false,
				},
				TotalCount: 10,
			},
		},
		{
			name:   "with_next_page",
			postID: postIDStr,
			first:  int32Ptr(2),
			after:  nil,
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("GetRootByPost", mock.Anything, postID, (*string)(nil), int64(0), int32(3)).
					Return([]*models.Comment{comment1, comment2, comment3}, nil)
				repo.On("TotalCount", mock.Anything, postID).Return(int64(10), nil)
			},
			want: &models.CommentConnection{
				Edges: []*models.CommentEdge{
					{Cursor: cursor.Encode(comment1.CreatedAt, comment1.ID), Node: comment1},
					{Cursor: cursor.Encode(comment2.CreatedAt, comment2.ID), Node: comment2},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursor.Encode(comment2.CreatedAt, comment2.ID)),
					HasNextPage: true,
				},
				TotalCount: 10,
			},
		},
		{
			name:   "with_cursor",
			postID: postIDStr,
			first:  int32Ptr(2),
			after:  strPtr(cursor.Encode(comment2.CreatedAt, comment2.ID)),
			setupMock: func(repo *mocks.MockCommentUC) {
				afterCreatedAt := comment2.CreatedAt
				afterID := comment2.ID
				repo.On("GetRootByPost", mock.Anything, postID, &afterCreatedAt, afterID, int32(3)).
					Return([]*models.Comment{comment3}, nil)
				repo.On("TotalCount", mock.Anything, postID).Return(int64(10), nil)
			},
			want: &models.CommentConnection{
				Edges: []*models.CommentEdge{
					{Cursor: cursor.Encode(comment3.CreatedAt, comment3.ID), Node: comment3},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursor.Encode(comment3.CreatedAt, comment3.ID)),
					HasNextPage: false,
				},
				TotalCount: 10,
			},
		},
		{
			name:   "cursor_beyond_last",
			postID: postIDStr,
			first:  int32Ptr(2),
			after:  strPtr(cursor.Encode(comment3.CreatedAt, comment3.ID)),
			setupMock: func(repo *mocks.MockCommentUC) {
				afterCreatedAt := comment3.CreatedAt
				afterID := comment3.ID
				repo.On("GetRootByPost", mock.Anything, postID, &afterCreatedAt, afterID, int32(3)).
					Return([]*models.Comment{}, nil)
				repo.On("TotalCount", mock.Anything, postID).Return(int64(10), nil)
			},
			want: &models.CommentConnection{
				Edges:      []*models.CommentEdge{},
				PageInfo:   &models.PageInfo{EndCursor: nil, HasNextPage: false},
				TotalCount: 10,
			},
		},
		{
			name:   "invalid_postID",
			postID: "abc",
			first:  nil,
			after:  nil,
			setupMock: func(repo *mocks.MockCommentUC) {
			},
			wantErr:     true,
			expectedErr: "invalid postID format",
		},
		{
			name:        "invalid_postID",
			postID:      "0",
			first:       nil,
			after:       nil,
			wantErr:     true,
			expectedErr: "postID must be greater 0",
		},
		{
			name:        "invalid_cursor_format",
			postID:      postIDStr,
			first:       int32Ptr(2),
			after:       strPtr("invalid"),
			setupMock:   func(repo *mocks.MockCommentUC) {},
			wantErr:     true,
			expectedErr: "invalid cursor format",
		},
		{
			name:   "repo_GetRootByPost_error",
			postID: postIDStr,
			first:  int32Ptr(2),
			after:  nil,
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("GetRootByPost", mock.Anything, postID, (*string)(nil), int64(0), int32(3)).
					Return(nil, errors.New("db error"))
			},
			wantErr:     true,
			expectedErr: "db error",
		},
		{
			name:   "repo_TotalCount_error",
			postID: postIDStr,
			first:  int32Ptr(2),
			after:  nil,
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("GetRootByPost", mock.Anything, postID, (*string)(nil), int64(0), int32(3)).
					Return([]*models.Comment{comment1, comment2}, nil)
				repo.On("TotalCount", mock.Anything, postID).Return(int64(0), errors.New("count error"))
			},
			wantErr:     true,
			expectedErr: "count error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockCommentUC(t)
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			s := New(mockRepo)
			got, err := s.GetRootComments(ctx, tt.postID, tt.first, tt.after)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetChildComments(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	parentID := int64(10)
	parentIDStr := "10"

	comment1 := &models.Comment{ID: 1, ParentID: &parentID, Author: "A", Text: "c1", CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	comment2 := &models.Comment{ID: 2, ParentID: &parentID, Author: "B", Text: "c2", CreatedAt: time.Now().UTC().Add(-time.Hour).Format(time.RFC3339)}
	comment3 := &models.Comment{ID: 3, ParentID: &parentID, Author: "C", Text: "c3", CreatedAt: time.Now().UTC().Add(-2 * time.Hour).Format(time.RFC3339)}

	tests := []struct {
		name        string
		parentID    string
		first       *int32
		after       *string
		setupMock   func(repo *mocks.MockCommentUC)
		want        *models.CommentConnection
		wantErr     bool
		expectedErr string
	}{
		{
			name:     "first_page_without_cursor",
			parentID: parentIDStr,
			first:    nil,
			after:    nil,
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("GetChild", mock.Anything, parentID, (*string)(nil), int64(0), int32(21)).
					Return([]*models.Comment{comment1, comment2, comment3}, nil)
			},
			want: &models.CommentConnection{
				Edges: []*models.CommentEdge{
					{Cursor: cursor.Encode(comment1.CreatedAt, comment1.ID), Node: comment1},
					{Cursor: cursor.Encode(comment2.CreatedAt, comment2.ID), Node: comment2},
					{Cursor: cursor.Encode(comment3.CreatedAt, comment3.ID), Node: comment3},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursor.Encode(comment3.CreatedAt, comment3.ID)),
					HasNextPage: false,
				},
				TotalCount: 3,
			},
		},
		{
			name:     "with_next_page",
			parentID: parentIDStr,
			first:    int32Ptr(2),
			after:    nil,
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("GetChild", mock.Anything, parentID, (*string)(nil), int64(0), int32(3)).
					Return([]*models.Comment{comment1, comment2, comment3}, nil)
			},
			want: &models.CommentConnection{
				Edges: []*models.CommentEdge{
					{Cursor: cursor.Encode(comment1.CreatedAt, comment1.ID), Node: comment1},
					{Cursor: cursor.Encode(comment2.CreatedAt, comment2.ID), Node: comment2},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursor.Encode(comment2.CreatedAt, comment2.ID)),
					HasNextPage: true,
				},
				TotalCount: 2,
			},
		},
		{
			name:     "with_cursor",
			parentID: parentIDStr,
			first:    int32Ptr(2),
			after:    strPtr(cursor.Encode(comment2.CreatedAt, comment2.ID)),
			setupMock: func(repo *mocks.MockCommentUC) {
				afterCreatedAt := comment2.CreatedAt
				afterID := comment2.ID
				repo.On("GetChild", mock.Anything, parentID, &afterCreatedAt, afterID, int32(3)).
					Return([]*models.Comment{comment3}, nil)
			},
			want: &models.CommentConnection{
				Edges: []*models.CommentEdge{
					{Cursor: cursor.Encode(comment3.CreatedAt, comment3.ID), Node: comment3},
				},
				PageInfo: &models.PageInfo{
					EndCursor:   strPtr(cursor.Encode(comment3.CreatedAt, comment3.ID)),
					HasNextPage: false,
				},
				TotalCount: 1,
			},
		},
		{
			name:        "invalid_parentID",
			parentID:    "abc",
			first:       nil,
			after:       nil,
			wantErr:     true,
			expectedErr: "invalid parentID",
		},
		{
			name:        "invalid_parentID",
			parentID:    "0",
			first:       nil,
			after:       nil,
			wantErr:     true,
			expectedErr: "invalid parentID",
		},
		{
			name:        "invalid_cursor",
			parentID:    parentIDStr,
			first:       int32Ptr(2),
			after:       strPtr("bad"),
			setupMock:   func(repo *mocks.MockCommentUC) {},
			wantErr:     true,
			expectedErr: "invalid cursor format",
		},
		{
			name:     "repo_GetChild_error",
			parentID: parentIDStr,
			first:    int32Ptr(2),
			after:    nil,
			setupMock: func(repo *mocks.MockCommentUC) {
				repo.On("GetChild", mock.Anything, parentID, (*string)(nil), int64(0), int32(3)).
					Return(nil, errors.New("db error"))
			},
			wantErr:     true,
			expectedErr: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockCommentUC(t)
			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			s := New(mockRepo)
			got, err := s.GetChildComments(ctx, tt.parentID, tt.first, tt.after)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

type mockDataloader struct {
	data map[string]interface{}
	err  error
}

func (m *mockDataloader) Load(ctx context.Context, key dataloader.StringKey) dataloader.Thunk {
	return func() (interface{}, error) {
		if m.err != nil {
			return nil, m.err
		}
		return m.data[string(key)], nil
	}
}

func strPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32 { return &i }
