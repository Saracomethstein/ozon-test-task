package post

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func TestPostRepository_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		postID        int64
		mockSetup     func(mock pgxmock.PgxPoolIface)
		expectedPost  *models.Post
		expectedError error
	}{
		{
			name:   "success",
			postID: 123,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "title", "body", "author", "allow_comments", "created_at",
				}).AddRow(
					int64(123),
					"Test Title",
					"Test Body",
					"Test Author",
					true,
					"2026-02-12T19:57:26Z",
				)

				mock.ExpectQuery(`select id, title, body, author, allow_comments, created_at from posts where id = \$1`).
					WithArgs(int64(123)).
					WillReturnRows(rows)
			},
			expectedPost: &models.Post{
				ID:            123,
				Title:         "Test Title",
				Body:          "Test Body",
				Author:        "Test Author",
				AllowComments: true,
				CreatedAt:     "2026-02-12T19:57:26Z",
			},
			expectedError: nil,
		},
		{
			name:   "post_not_found",
			postID: 999,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select id, title, body, author, allow_comments, created_at from posts where id = \$1`).
					WithArgs(int64(999)).
					WillReturnError(pgx.ErrNoRows)
			},
			expectedPost:  nil,
			expectedError: errors.New("post not found"),
		},
		{
			name:   "db_error",
			postID: 500,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select id, title, body, author, allow_comments, created_at from posts where id = \$1`).
					WithArgs(int64(500)).
					WillReturnError(errors.New("db error"))
			},
			expectedPost:  nil,
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			repo := New(mock)

			tt.mockSetup(mock)

			result, err := repo.GetByID(context.Background(), tt.postID)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedError.Error())
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPost, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostRepository_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		afterCreated  *string
		afterID       int64
		limit         int32
		mockSetup     func(mock pgxmock.PgxPoolIface)
		expectedPosts []*models.Post
		expectError   bool
	}{
		{
			name:         "success_multiple_posts",
			afterCreated: nil,
			afterID:      0,
			limit:        2,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "title", "body", "author", "allow_comments", "created_at",
				}).
					AddRow(int64(2), "Title2", "Body2", "Author2", true, "2026-02-12T20:00:00Z").
					AddRow(int64(1), "Title1", "Body1", "Author1", false, "2026-02-12T19:00:00Z")

				mock.ExpectQuery(`select id, title, body, author, allow_comments, created_at from posts`).
					WithArgs(pgxmock.AnyArg(), int64(0), int32(2)).
					WillReturnRows(rows)
			},
			expectedPosts: []*models.Post{
				{
					ID:            2,
					Title:         "Title2",
					Body:          "Body2",
					Author:        "Author2",
					AllowComments: true,
					CreatedAt:     "2026-02-12T20:00:00Z",
				},
				{
					ID:            1,
					Title:         "Title1",
					Body:          "Body1",
					Author:        "Author1",
					AllowComments: false,
					CreatedAt:     "2026-02-12T19:00:00Z",
				},
			},
			expectError: false,
		},
		{
			name:         "empty_result",
			afterCreated: nil,
			afterID:      0,
			limit:        5,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "title", "body", "author", "allow_comments", "created_at",
				})

				mock.ExpectQuery(`select id, title, body, author, allow_comments, created_at from posts`).
					WithArgs(pgxmock.AnyArg(), int64(0), int32(5)).
					WillReturnRows(rows)
			},
			expectedPosts: []*models.Post{},
			expectError:   false,
		},
		{
			name:         "query_error",
			afterCreated: nil,
			afterID:      0,
			limit:        5,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select id, title, body, author, allow_comments, created_at from posts`).
					WithArgs(pgxmock.AnyArg(), int64(0), int32(5)).
					WillReturnError(errors.New("db error"))
			},
			expectedPosts: nil,
			expectError:   true,
		},
		{
			name:         "scan_error",
			afterCreated: nil,
			afterID:      0,
			limit:        1,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "title", "body", "author", "allow_comments", "created_at",
				}).AddRow("wrong_type", "Title", "Body", "Author", true, "2026-02-12T20:00:00Z")

				mock.ExpectQuery(`select id, title, body, author, allow_comments, created_at from posts`).
					WithArgs(pgxmock.AnyArg(), int64(0), int32(1)).
					WillReturnRows(rows)
			},
			expectedPosts: nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			repo := New(mock)

			tt.mockSetup(mock)

			result, err := repo.Get(context.Background(), tt.afterCreated, tt.afterID, tt.limit)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedPosts, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostRepository_TotalCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mockSetup     func(mock pgxmock.PgxPoolIface)
		expectedCount int64
		expectedError error
	}{
		{
			name: "success",
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"count"}).
					AddRow(int64(42))

				mock.ExpectQuery(`select count\(\*\) from posts`).
					WillReturnRows(rows)
			},
			expectedCount: 42,
			expectedError: nil,
		},
		{
			name: "query_error",
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select count\(\*\) from posts`).
					WillReturnError(errors.New("db error"))
			},
			expectedCount: 0,
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			repo := New(mock)

			tt.mockSetup(mock)

			count, err := repo.TotalCount(context.Background())

			if tt.expectedError != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedCount, count)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostRepository_Save(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		inputPost   models.Post
		mockSetup   func(mock pgxmock.PgxPoolIface, input models.Post)
		expectError bool
	}{
		{
			name: "success",
			inputPost: models.Post{
				Title:         "Test Title",
				Body:          "Test Body",
				Author:        "Adel",
				AllowComments: true,
				CreatedAt:     "2026-02-12T21:00:00Z",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, input models.Post) {
				rows := pgxmock.NewRows([]string{
					"id", "title", "body", "author", "allow_comments", "created_at",
				}).AddRow(
					int64(100),
					input.Title,
					input.Body,
					input.Author,
					input.AllowComments,
					input.CreatedAt,
				)

				mock.ExpectQuery(`insert into posts`).
					WithArgs(
						input.Title,
						input.Body,
						input.Author,
						input.AllowComments,
						input.CreatedAt,
					).
					WillReturnRows(rows)
			},
			expectError: false,
		},
		{
			name: "db_error",
			inputPost: models.Post{
				Title:         "Test Title",
				Body:          "Test Body",
				Author:        "Adel",
				AllowComments: true,
				CreatedAt:     "2026-02-12T21:00:00Z",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, input models.Post) {
				mock.ExpectQuery(`insert into posts`).
					WithArgs(
						input.Title,
						input.Body,
						input.Author,
						input.AllowComments,
						input.CreatedAt,
					).
					WillReturnError(errors.New("insert error"))
			},
			expectError: true,
		},
		{
			name: "scan_error",
			inputPost: models.Post{
				Title:         "Test Title",
				Body:          "Test Body",
				Author:        "Adel",
				AllowComments: true,
				CreatedAt:     "2026-02-12T21:00:00Z",
			},
			mockSetup: func(mock pgxmock.PgxPoolIface, input models.Post) {
				rows := pgxmock.NewRows([]string{
					"id", "title", "body", "author", "allow_comments", "created_at",
				}).AddRow(
					"wrong_id_type",
					input.Title,
					input.Body,
					input.Author,
					input.AllowComments,
					input.CreatedAt,
				)

				mock.ExpectQuery(`insert into posts`).
					WithArgs(
						input.Title,
						input.Body,
						input.Author,
						input.AllowComments,
						input.CreatedAt,
					).
					WillReturnRows(rows)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			repo := New(mock)

			tt.mockSetup(mock, tt.inputPost)

			result, err := repo.Save(context.Background(), tt.inputPost)

			if tt.expectError {
				require.Error(t, err)
				require.Equal(t, models.Post{}, result)
			} else {
				require.NoError(t, err)

				require.Equal(t, tt.inputPost.Title, result.Title)
				require.Equal(t, tt.inputPost.Body, result.Body)
				require.Equal(t, tt.inputPost.Author, result.Author)
				require.Equal(t, tt.inputPost.AllowComments, result.AllowComments)
				require.Equal(t, tt.inputPost.CreatedAt, result.CreatedAt)

				require.NotZero(t, result.ID)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPostRepository_SetCommentsAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		postID      int64
		allow       bool
		mockSetup   func(mock pgxmock.PgxPoolIface)
		expected    *models.Post
		expectError bool
	}{
		{
			name:   "success",
			postID: 10,
			allow:  true,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "title", "body", "author", "allow_comments", "created_at",
				}).AddRow(
					int64(10),
					"Updated Title",
					"Body",
					"Adel",
					true,
					"2026-02-12T22:00:00Z",
				)

				mock.ExpectQuery(`update posts`).
					WithArgs(int64(10), true).
					WillReturnRows(rows)
			},
			expected: &models.Post{
				ID:            10,
				Title:         "Updated Title",
				Body:          "Body",
				Author:        "Adel",
				AllowComments: true,
				CreatedAt:     "2026-02-12T22:00:00Z",
			},
			expectError: false,
		},
		{
			name:   "post_not_found",
			postID: 999,
			allow:  false,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`update posts`).
					WithArgs(int64(999), false).
					WillReturnError(pgx.ErrNoRows)
			},
			expected:    nil,
			expectError: true,
		},
		{
			name:   "db_error",
			postID: 15,
			allow:  false,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`update posts`).
					WithArgs(int64(15), false).
					WillReturnError(errors.New("db error"))
			},
			expected:    nil,
			expectError: true,
		},
		{
			name:   "scan_error",
			postID: 20,
			allow:  true,
			mockSetup: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{
					"id", "title", "body", "author", "allow_comments", "created_at",
				}).AddRow(
					"wrong_id_type",
					"Title",
					"Body",
					"Author",
					true,
					"2026-02-12T22:00:00Z",
				)

				mock.ExpectQuery(`update posts`).
					WithArgs(int64(20), true).
					WillReturnRows(rows)
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			repo := New(mock)

			tt.mockSetup(mock)

			result, err := repo.SetCommentsAllowed(
				context.Background(),
				tt.postID,
				tt.allow,
			)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, result)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
