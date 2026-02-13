package comment

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	baseComment := models.Comment{
		PostID:    1,
		ParentID:  nil,
		Author:    "John Doe",
		Text:      "Test comment",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	tests := []struct {
		name        string
		comment     models.Comment
		setupMock   func(mock pgxmock.PgxPoolIface)
		wantID      int64
		wantErr     bool
		expectedErr error
	}{
		{
			name:    "success",
			comment: baseComment,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`insert into comments`).
					WithArgs(baseComment.PostID, baseComment.ParentID, baseComment.Author, baseComment.Text, baseComment.CreatedAt).
					WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int64(123)))
			},
			wantID:  123,
			wantErr: false,
		},
		{
			name:    "db_error",
			comment: baseComment,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`insert into comments`).
					WithArgs(baseComment.PostID, baseComment.ParentID, baseComment.Author, baseComment.Text, baseComment.CreatedAt).
					WillReturnError(errors.New("insert failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.setupMock(mock)

			r := New(mock)
			got, err := r.Add(context.Background(), tt.comment)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, got)
				assert.Equal(t, tt.wantID, got.ID)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCheckAllowComments(t *testing.T) {
	t.Parallel()

	postID := int64(1)

	tests := []struct {
		name        string
		postID      int64
		setupMock   func(mock pgxmock.PgxPoolIface)
		wantAllow   bool
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "comments_allowed",
			postID: postID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select allow_comments from posts`).
					WithArgs(postID).
					WillReturnRows(pgxmock.NewRows([]string{"allow_comments"}).AddRow(true))
			},
			wantAllow: true,
			wantErr:   false,
		},
		{
			name:   "comments_not_allowed",
			postID: postID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select allow_comments from posts`).
					WithArgs(postID).
					WillReturnRows(pgxmock.NewRows([]string{"allow_comments"}).AddRow(false))
			},
			wantAllow: false,
			wantErr:   false,
		},
		{
			name:   "post_not_found",
			postID: postID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select allow_comments from posts`).
					WithArgs(postID).
					WillReturnError(sql.ErrNoRows)
			},
			wantAllow:   false,
			wantErr:     true,
			expectedErr: errors.New("post not found"),
		},
		{
			name:   "db_error",
			postID: postID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select allow_comments from posts`).
					WithArgs(postID).
					WillReturnError(errors.New("connection error"))
			},
			wantAllow: false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.setupMock(mock)

			r := New(mock)
			got, err := r.CheckAllowComments(context.Background(), tt.postID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.False(t, got)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantAllow, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCheckParentExists(t *testing.T) {
	t.Parallel()

	parentID := int64(100)
	expectedPostID := int64(1)

	tests := []struct {
		name        string
		parentID    int64
		setupMock   func(mock pgxmock.PgxPoolIface)
		wantPostID  int64
		wantErr     bool
		expectedErr error
	}{
		{
			name:     "parent_exists",
			parentID: parentID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select post_id from comments`).
					WithArgs(parentID).
					WillReturnRows(pgxmock.NewRows([]string{"post_id"}).AddRow(expectedPostID))
			},
			wantPostID: expectedPostID,
			wantErr:    false,
		},
		{
			name:     "parent_not_found",
			parentID: parentID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select post_id from comments`).
					WithArgs(parentID).
					WillReturnError(sql.ErrNoRows)
			},
			wantPostID:  0,
			wantErr:     true,
			expectedErr: errors.New("parent comment not found"),
		},
		{
			name:     "db_error",
			parentID: parentID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select post_id from comments`).
					WithArgs(parentID).
					WillReturnError(errors.New("some db error"))
			},
			wantPostID: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.setupMock(mock)

			r := New(mock)
			got, err := r.CheckParentExists(context.Background(), tt.parentID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int64(0), got)
				if tt.expectedErr != nil {
					assert.EqualError(t, err, tt.expectedErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPostID, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetRootByPost(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	postID := int64(1)
	afterCreatedAt := "2023-01-01T00:00:00Z"
	afterID := int64(10)
	limit := int32(5)

	comment1 := testComment(1, postID, nil, "Alice", "First", now.Format(time.RFC3339))
	comment2 := testComment(2, postID, nil, "Bob", "Second", now.Add(-time.Hour).Format(time.RFC3339))

	tests := []struct {
		name      string
		postID    int64
		after     *string
		afterID   int64
		limit     int32
		setupMock func(pgxmock.PgxPoolIface)
		want      []*models.Comment
		wantErr   bool
	}{
		{
			name:    "success_with_results",
			postID:  postID,
			after:   &afterCreatedAt,
			afterID: afterID,
			limit:   limit,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "post_id", "parent_id", "author", "body", "created_at"}).
					AddRow(comment1.ID, comment1.PostID, comment1.ParentID, comment1.Author, comment1.Text, comment1.CreatedAt).
					AddRow(comment2.ID, comment2.PostID, comment2.ParentID, comment2.Author, comment2.Text, comment2.CreatedAt)
				mock.ExpectQuery(`select id, post_id, parent_id, author, body, created_at from comments where post_id = \$1 and parent_id is null and \(\$2::text is null or \(created_at, id\) < \(\$2::text, \$3::bigint\)\) order by created_at desc, id desc limit \$4`).
					WithArgs(postID, &afterCreatedAt, afterID, limit).
					WillReturnRows(rows)
			},
			want: []*models.Comment{&comment1, &comment2},
		},
		{
			name:    "success_with_nil_after",
			postID:  postID,
			after:   nil,
			afterID: 0,
			limit:   limit,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "post_id", "parent_id", "author", "body", "created_at"}).
					AddRow(comment1.ID, comment1.PostID, comment1.ParentID, comment1.Author, comment1.Text, comment1.CreatedAt)
				mock.ExpectQuery(`select id, post_id, parent_id, author, body, created_at from comments where post_id = \$1 and parent_id is null and \(\$2::text is null or \(created_at, id\) < \(\$2::text, \$3::bigint\)\) order by created_at desc, id desc limit \$4`).
					WithArgs(postID, (*string)(nil), int64(0), limit).
					WillReturnRows(rows)
			},
			want: []*models.Comment{&comment1},
		},
		{
			name:    "no_rows",
			postID:  postID,
			after:   &afterCreatedAt,
			afterID: afterID,
			limit:   limit,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select id, post_id, parent_id, author, body, created_at from comments where post_id = \$1 and parent_id is null and \(\$2::text is null or \(created_at, id\) < \(\$2::text, \$3::bigint\)\) order by created_at desc, id desc limit \$4`).
					WithArgs(postID, &afterCreatedAt, afterID, limit).
					WillReturnRows(pgxmock.NewRows([]string{"id", "post_id", "parent_id", "author", "body", "created_at"}))
			},
			want: []*models.Comment{},
		},
		{
			name:    "db_error",
			postID:  postID,
			after:   &afterCreatedAt,
			afterID: afterID,
			limit:   limit,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select id, post_id, parent_id, author, body, created_at from comments where post_id = \$1 and parent_id is null and \(\$2::text is null or \(created_at, id\) < \(\$2::text, \$3::bigint\)\) order by created_at desc, id desc limit \$4`).
					WithArgs(postID, &afterCreatedAt, afterID, limit).
					WillReturnError(errors.New("query failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.setupMock(mock)

			r := New(mock)
			got, err := r.GetRootByPost(context.Background(), tt.postID, tt.after, tt.afterID, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
func TestGetChild(t *testing.T) {
	t.Parallel()

	parentID := int64(100)
	afterCreatedAt := "2026-02-12T19:57:26Z"
	afterID := int64(10)
	limit := int32(5)

	comment1 := testComment(1, 1, &parentID, "Alice", "Child 1", "2026-02-12T19:57:26Z")
	comment2 := testComment(2, 1, &parentID, "Bob", "Child 2", "2026-02-12T18:57:26Z")

	tests := []struct {
		name      string
		parentID  int64
		after     *string
		afterID   int64
		limit     int32
		setupMock func(pgxmock.PgxPoolIface)
		want      []*models.Comment
		wantErr   bool
	}{
		{
			name:     "success_with_results",
			parentID: parentID,
			after:    &afterCreatedAt,
			afterID:  afterID,
			limit:    limit,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"id", "post_id", "parent_id", "author", "body", "created_at"}).
					AddRow(comment1.ID, comment1.PostID, comment1.ParentID, comment1.Author, comment1.Text, comment1.CreatedAt).
					AddRow(comment2.ID, comment2.PostID, comment2.ParentID, comment2.Author, comment2.Text, comment2.CreatedAt)
				mock.ExpectQuery(`select id, post_id, parent_id, author, body, created_at from comments where parent_id = \$1 and \(\$2::text is null or \(created_at, id\) < \(\$2::text, \$3::bigint\)\) order by created_at desc, id desc limit \$4`).
					WithArgs(parentID, &afterCreatedAt, afterID, limit).
					WillReturnRows(rows)
			},
			want: []*models.Comment{&comment1, &comment2},
		},
		{
			name:     "no_rows",
			parentID: parentID,
			after:    &afterCreatedAt,
			afterID:  afterID,
			limit:    limit,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select id, post_id, parent_id, author, body, created_at from comments where parent_id = \$1 and \(\$2::text is null or \(created_at, id\) < \(\$2::text, \$3::bigint\)\) order by created_at desc, id desc limit \$4`).
					WithArgs(parentID, &afterCreatedAt, afterID, limit).
					WillReturnRows(pgxmock.NewRows([]string{"id", "post_id", "parent_id", "author", "body", "created_at"}))
			},
			want: []*models.Comment{},
		},
		{
			name:     "db_error",
			parentID: parentID,
			after:    &afterCreatedAt,
			afterID:  afterID,
			limit:    limit,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select id, post_id, parent_id, author, body, created_at from comments where parent_id = \$1 and \(\$2::text is null or \(created_at, id\) < \(\$2::text, \$3::bigint\)\) order by created_at desc, id desc limit \$4`).
					WithArgs(parentID, &afterCreatedAt, afterID, limit).
					WillReturnError(errors.New("query failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.setupMock(mock)

			r := New(mock)
			got, err := r.GetChild(context.Background(), tt.parentID, tt.after, tt.afterID, tt.limit)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetChildBatch(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	parentIDs := []int64{100, 200}

	comment1 := models.Comment{
		ID:        1,
		PostID:    1,
		ParentID:  &parentIDs[0],
		Author:    "Alice",
		Text:      "Child of 100",
		CreatedAt: now.Format(time.RFC3339),
	}
	comment2 := models.Comment{
		ID:        2,
		PostID:    1,
		ParentID:  &parentIDs[0],
		Author:    "Bob",
		Text:      "Another child of 100",
		CreatedAt: now.Add(-1 * time.Hour).Format(time.RFC3339),
	}
	comment3 := models.Comment{
		ID:        3,
		PostID:    2,
		ParentID:  &parentIDs[1],
		Author:    "Charlie",
		Text:      "Child of 200",
		CreatedAt: now.Add(-2 * time.Hour).Format(time.RFC3339),
	}

	tests := []struct {
		name      string
		parentIDs []int64
		setupMock func(pgxmock.PgxPoolIface)
		want      []*models.Comment
		wantErr   bool
	}{
		{
			name:      "success_with_results",
			parentIDs: parentIDs,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				rows := pgxmock.NewRows([]string{"parent_id", "id", "post_id", "author", "body", "created_at"}).
					AddRow(comment1.ParentID, comment1.ID, comment1.PostID, comment1.Author, comment1.Text, comment1.CreatedAt).
					AddRow(comment2.ParentID, comment2.ID, comment2.PostID, comment2.Author, comment2.Text, comment2.CreatedAt).
					AddRow(comment3.ParentID, comment3.ID, comment3.PostID, comment3.Author, comment3.Text, comment3.CreatedAt)
				mock.ExpectQuery(`select parent_id, id, post_id, author, body, created_at from comments where parent_id = any\(\$1\) order by parent_id, created_at desc, id desc`).
					WithArgs(parentIDs).
					WillReturnRows(rows)
			},
			want: []*models.Comment{&comment1, &comment2, &comment3},
		},
		{
			name:      "no_rows",
			parentIDs: parentIDs,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select parent_id, id, post_id, author, body, created_at from comments where parent_id = any\(\$1\) order by parent_id, created_at desc, id desc`).
					WithArgs(parentIDs).
					WillReturnRows(pgxmock.NewRows([]string{"parent_id", "id", "post_id", "author", "body", "created_at"}))
			},
			want: []*models.Comment{},
		},
		{
			name:      "db_error",
			parentIDs: parentIDs,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select parent_id, id, post_id, author, body, created_at from comments where parent_id = any\(\$1\) order by parent_id, created_at desc, id desc`).
					WithArgs(parentIDs).
					WillReturnError(errors.New("batch query failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.setupMock(mock)

			r := New(mock)
			got, err := r.GetChildBatch(context.Background(), tt.parentIDs)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTotalCount(t *testing.T) {
	t.Parallel()

	postID := int64(1)

	tests := []struct {
		name        string
		postID      int64
		setupMock   func(pgxmock.PgxPoolIface)
		wantCount   int64
		wantErr     bool
		expectedErr error
	}{
		{
			name:   "success",
			postID: postID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select count\(\*\) from comments where post_id = \$1`).
					WithArgs(postID).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(42)))
			},
			wantCount: 42,
		},
		{
			name:   "zero_count",
			postID: postID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select count\(\*\) from comments where post_id = \$1`).
					WithArgs(postID).
					WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(0)))
			},
			wantCount: 0,
		},
		{
			name:   "db_error",
			postID: postID,
			setupMock: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(`select count\(\*\) from comments where post_id = \$1`).
					WithArgs(postID).
					WillReturnError(errors.New("count failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mock, err := pgxmock.NewPool()
			require.NoError(t, err)
			defer mock.Close()

			tt.setupMock(mock)

			r := New(mock)
			got, err := r.TotalCount(context.Background(), tt.postID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int64(0), got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func testComment(id, postID int64, parentID *int64, author, text string, createdAt string) models.Comment {
	return models.Comment{
		ID:        id,
		PostID:    postID,
		ParentID:  parentID,
		Author:    author,
		Text:      text,
		CreatedAt: createdAt,
	}
}
