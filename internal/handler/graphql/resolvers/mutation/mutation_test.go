package mutation

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

func TestMutationResolver_AddComment(t *testing.T) {
	t.Parallel()

	postID := "123"
	postIDInt := int64(123)
	parentID := "456"
	parentIDInt := int64(456)
	author := "John Doe"
	text := "Test comment"
	createdAt := "2023-01-01T12:00:00Z"

	tests := []struct {
		name        string
		input       graphql.AddCommentInput
		mockSetup   func(mockSvc *mockComment.MockUseCase)
		expected    *graphql.Comment
		expectedErr string
	}{
		{
			name: "success_without_parent",
			input: graphql.AddCommentInput{
				PostID:   postID,
				Author:   author,
				Text:     text,
				ParentID: nil,
			},
			mockSetup: func(mockSvc *mockComment.MockUseCase) {
				mockSvc.EXPECT().
					AddComment(mock.Anything, models.AddCommentInput{
						PostID:   postID,
						Author:   author,
						Text:     text,
						ParentID: nil,
					}).
					Return(&models.Comment{
						ID:        1,
						PostID:    postIDInt,
						Author:    author,
						Text:      text,
						ParentID:  nil,
						CreatedAt: createdAt,
					}, nil)
			},
			expected: &graphql.Comment{
				ID:        "1",
				PostID:    postID,
				Author:    author,
				Text:      text,
				ParentID:  nil,
				CreatedAt: createdAt,
			},
		},
		{
			name: "success_with_parent",
			input: graphql.AddCommentInput{
				PostID:   postID,
				Author:   author,
				Text:     text,
				ParentID: &parentID,
			},
			mockSetup: func(mockSvc *mockComment.MockUseCase) {
				mockSvc.EXPECT().
					AddComment(mock.Anything, models.AddCommentInput{
						PostID:   postID,
						Author:   author,
						Text:     text,
						ParentID: &parentID,
					}).
					Return(&models.Comment{
						ID:        2,
						PostID:    postIDInt,
						Author:    author,
						Text:      text,
						ParentID:  &parentIDInt,
						CreatedAt: createdAt,
					}, nil)
			},
			expected: &graphql.Comment{
				ID:        "2",
				PostID:    postID,
				Author:    author,
				Text:      text,
				ParentID:  &parentID,
				CreatedAt: createdAt,
			},
		},
		{
			name: "validation_error",
			input: graphql.AddCommentInput{
				PostID:   postID,
				Author:   "",
				Text:     text,
				ParentID: nil,
			},
			mockSetup:   func(mockSvc *mockComment.MockUseCase) {},
			expectedErr: "autor, text or postID cannot be empty",
		},
		{
			name: "validation_error",
			input: graphql.AddCommentInput{
				PostID:   postID,
				Author:   author,
				Text:     "",
				ParentID: nil,
			},
			mockSetup:   func(mockSvc *mockComment.MockUseCase) {},
			expectedErr: "autor, text or postID cannot be empty",
		},
		{
			name: "validation_error",
			input: graphql.AddCommentInput{
				PostID:   "",
				Author:   author,
				Text:     text,
				ParentID: nil,
			},
			mockSetup:   func(mockSvc *mockComment.MockUseCase) {},
			expectedErr: "autor, text or postID cannot be empty",
		},
		{
			name: "service_error",
			input: graphql.AddCommentInput{
				PostID:   postID,
				Author:   author,
				Text:     text,
				ParentID: nil,
			},
			mockSetup: func(mockSvc *mockComment.MockUseCase) {
				mockSvc.EXPECT().
					AddComment(mock.Anything, models.AddCommentInput{
						PostID:   postID,
						Author:   author,
						Text:     text,
						ParentID: nil,
					}).
					Return(nil, errors.New("internal error"))
			},
			expectedErr: "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCommentService := mockComment.NewMockUseCase(t)
			resolver := &mutationResolver{
				service: &service.Container{
					CommentService: mockCommentService,
				},
			}

			tt.mockSetup(mockCommentService)

			got, err := resolver.AddComment(context.Background(), tt.input)

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

func TestMutationResolver_CreatePost(t *testing.T) {
	t.Parallel()

	title := "Test Title"
	author := "John Doe"
	body := "Test body content"
	allowComments := true
	createdAt := "2023-01-01T12:00:00Z"
	postID := int64(1)

	tests := []struct {
		name        string
		input       graphql.CreatePostInput
		mockSetup   func(mockSvc *mockPost.MockUseCase)
		expected    *graphql.Post
		expectedErr string
	}{
		{
			name: "success",
			input: graphql.CreatePostInput{
				Title:         title,
				Author:        author,
				Body:          body,
				AllowComments: &allowComments,
			},
			mockSetup: func(mockSvc *mockPost.MockUseCase) {
				mockSvc.EXPECT().
					CreatePost(mock.Anything, models.CreatePostInput{
						Title:         title,
						Author:        author,
						Body:          body,
						AllowComments: &allowComments,
					}).
					Return(&models.Post{
						ID:            postID,
						Title:         title,
						Author:        author,
						Body:          body,
						AllowComments: allowComments,
						CreatedAt:     createdAt,
					}, nil)
			},
			expected: &graphql.Post{
				ID:            "1",
				Title:         title,
				Author:        author,
				Body:          body,
				AllowComments: allowComments,
				CreatedAt:     createdAt,
			},
		},
		{
			name: "validation_error",
			input: graphql.CreatePostInput{
				Title:         "",
				Author:        author,
				Body:          body,
				AllowComments: &allowComments,
			},
			mockSetup:   func(mockSvc *mockPost.MockUseCase) {},
			expectedErr: "title, author and body are required fields",
		},
		{
			name: "validation_error",
			input: graphql.CreatePostInput{
				Title:         title,
				Author:        "",
				Body:          body,
				AllowComments: &allowComments,
			},
			mockSetup:   func(mockSvc *mockPost.MockUseCase) {},
			expectedErr: "title, author and body are required fields",
		},
		{
			name: "validation_error",
			input: graphql.CreatePostInput{
				Title:         title,
				Author:        author,
				Body:          "",
				AllowComments: &allowComments,
			},
			mockSetup:   func(mockSvc *mockPost.MockUseCase) {},
			expectedErr: "title, author and body are required fields",
		},
		{
			name: "service_error",
			input: graphql.CreatePostInput{
				Title:         title,
				Author:        author,
				Body:          body,
				AllowComments: &allowComments,
			},
			mockSetup: func(mockSvc *mockPost.MockUseCase) {
				mockSvc.EXPECT().
					CreatePost(mock.Anything, models.CreatePostInput{
						Title:         title,
						Author:        author,
						Body:          body,
						AllowComments: &allowComments,
					}).
					Return(nil, errors.New("database error"))
			},
			expectedErr: "failed to create post: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockPostService := mockPost.NewMockUseCase(t)
			resolver := &mutationResolver{
				service: &service.Container{
					PostService: mockPostService,
				},
			}

			tt.mockSetup(mockPostService)

			got, err := resolver.CreatePost(context.Background(), tt.input)

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

func TestMutationResolver_SetPostCommentsAllowed(t *testing.T) {
	t.Parallel()

	postID := "123"
	postIDInt := int64(123)
	allow := true
	title := "Test Title"
	author := "John Doe"
	body := "Test body"
	createdAt := "2023-01-01T12:00:00Z"

	tests := []struct {
		name        string
		postID      string
		allow       bool
		mockSetup   func(mockSvc *mockPost.MockUseCase)
		expected    *graphql.Post
		expectedErr string
	}{
		{
			name:   "success",
			postID: postID,
			allow:  allow,
			mockSetup: func(mockSvc *mockPost.MockUseCase) {
				mockSvc.EXPECT().
					SetPostCommentsAllowed(mock.Anything, postID, allow).
					Return(&models.Post{
						ID:            postIDInt,
						Title:         title,
						Author:        author,
						Body:          body,
						AllowComments: allow,
						CreatedAt:     createdAt,
					}, nil)
			},
			expected: &graphql.Post{
				ID:            postID,
				Title:         title,
				Author:        author,
				Body:          body,
				AllowComments: allow,
				CreatedAt:     createdAt,
			},
		},
		{
			name:        "validation_error",
			postID:      "",
			allow:       allow,
			mockSetup:   func(mockSvc *mockPost.MockUseCase) {},
			expectedErr: "postId cannot be empty",
		},
		{
			name:   "service_error",
			postID: postID,
			allow:  allow,
			mockSetup: func(mockSvc *mockPost.MockUseCase) {
				mockSvc.EXPECT().
					SetPostCommentsAllowed(mock.Anything, postID, allow).
					Return(nil, errors.New("db error"))
			},
			expectedErr: "failed to set post comments allowed: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockPostService := mockPost.NewMockUseCase(t)
			resolver := &mutationResolver{
				service: &service.Container{
					PostService: mockPostService,
				},
			}

			tt.mockSetup(mockPostService)

			got, err := resolver.SetPostCommentsAllowed(context.Background(), tt.postID, tt.allow)

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
