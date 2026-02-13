package post

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
)

func TestPostRepo_Save(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Now().Format(time.RFC3339)

	t.Run("successful_save", func(t *testing.T) {
		repo := New()
		input := models.Post{
			Title:         "First Post",
			Body:          "Content",
			Author:        "Alice",
			AllowComments: true,
			CreatedAt:     now,
		}

		got, err := repo.Save(ctx, input)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), got.ID)
		assert.Equal(t, input.Title, got.Title)
		assert.Equal(t, input.Body, got.Body)
		assert.Equal(t, input.Author, got.Author)
		assert.Equal(t, input.AllowComments, got.AllowComments)
		assert.Equal(t, input.CreatedAt, got.CreatedAt)
	})

	t.Run("multiple_saves_increment_id", func(t *testing.T) {
		repo := New()

		p1, _ := repo.Save(ctx, models.Post{Title: "A", CreatedAt: now})
		p2, _ := repo.Save(ctx, models.Post{Title: "B", CreatedAt: now})

		assert.Equal(t, int64(1), p1.ID)
		assert.Equal(t, int64(2), p2.ID)
	})
}

func TestPostRepo_GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Now().Format(time.RFC3339)

	t.Run("existing_post", func(t *testing.T) {
		repo := New()
		saved, _ := repo.Save(ctx, models.Post{Title: "Test", CreatedAt: now})

		got, err := repo.GetByID(ctx, saved.ID)

		assert.NoError(t, err)
		assert.Equal(t, saved.ID, got.ID)
		assert.Equal(t, "Test", got.Title)
	})

	t.Run("non-existing_post", func(t *testing.T) {
		repo := New()

		got, err := repo.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.EqualError(t, err, "post not found")
		assert.Nil(t, got)
	})

	t.Run("returns_copy", func(t *testing.T) {
		repo := New()
		saved, _ := repo.Save(ctx, models.Post{Title: "Original", CreatedAt: now})

		got, _ := repo.GetByID(ctx, saved.ID)
		got.Title = "Modified"

		original, _ := repo.GetByID(ctx, saved.ID)
		assert.Equal(t, "Original", original.Title)
	})
}

func TestPostRepo_SetCommentsAllowed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Now().Format(time.RFC3339)

	t.Run("update_existing_post", func(t *testing.T) {
		repo := New()
		saved, _ := repo.Save(ctx, models.Post{Title: "Test", AllowComments: false, CreatedAt: now})

		updated, err := repo.SetCommentsAllowed(ctx, saved.ID, true)

		assert.NoError(t, err)
		assert.True(t, updated.AllowComments)
		assert.Equal(t, saved.ID, updated.ID)

		got, _ := repo.GetByID(ctx, saved.ID)
		assert.True(t, got.AllowComments)
	})

	t.Run("post_not_found", func(t *testing.T) {
		repo := New()

		updated, err := repo.SetCommentsAllowed(ctx, 999, true)

		assert.Error(t, err)
		assert.EqualError(t, err, "post not found")
		assert.Nil(t, updated)
	})

	t.Run("return_copy", func(t *testing.T) {
		repo := New()
		saved, _ := repo.Save(ctx, models.Post{Title: "Test", AllowComments: false, CreatedAt: now})

		updated, _ := repo.SetCommentsAllowed(ctx, saved.ID, true)
		updated.Title = "Hacked"

		original, _ := repo.GetByID(ctx, saved.ID)
		assert.Equal(t, "Test", original.Title)
	})
}

func TestPostRepo_TotalCount(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Now().Format(time.RFC3339)

	t.Run("empty_repo", func(t *testing.T) {
		repo := New()
		count, err := repo.TotalCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("with_posts", func(t *testing.T) {
		repo := New()
		repo.Save(ctx, models.Post{Title: "1", CreatedAt: now})
		repo.Save(ctx, models.Post{Title: "2", CreatedAt: now})

		count, err := repo.TotalCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})
}

func TestPostRepo_Get(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	now := time.Now().UTC()

	posts := []models.Post{
		{ID: 1, Title: "Post 1", CreatedAt: now.Format(time.RFC3339)},
		{ID: 2, Title: "Post 2", CreatedAt: now.Add(-1 * time.Hour).Format(time.RFC3339)},
		{ID: 3, Title: "Post 3", CreatedAt: now.Add(-2 * time.Hour).Format(time.RFC3339)},
		{ID: 4, Title: "Post 4", CreatedAt: now.Add(-3 * time.Hour).Format(time.RFC3339)},
	}

	fillRepo := func() repository.PostUC {
		repo := New()
		for _, p := range posts {
			repo.Save(ctx, models.Post{
				Title:     p.Title,
				Body:      "test post",
				Author:    "test autor",
				CreatedAt: p.CreatedAt,
			})
		}
		return repo
	}

	t.Run("first_page_without_cursor", func(t *testing.T) {
		repo := fillRepo()
		limit := int32(2)

		got, err := repo.Get(ctx, nil, 0, limit)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "Post 1", got[0].Title)
		assert.Equal(t, "Post 2", got[1].Title)
	})

	t.Run("second_page_with_cursor", func(t *testing.T) {
		repo := fillRepo()
		afterCreated := posts[1].CreatedAt
		afterID := int64(2)
		limit := int32(2)

		got, err := repo.Get(ctx, &afterCreated, afterID, limit)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "Post 3", got[0].Title)
		assert.Equal(t, "Post 4", got[1].Title)
	})

	t.Run("cursor_points_to_last_element", func(t *testing.T) {
		repo := fillRepo()
		afterCreated := posts[3].CreatedAt
		afterID := int64(4)
		limit := int32(2)

		got, err := repo.Get(ctx, &afterCreated, afterID, limit)

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("limit_larger_than_remaining_posts", func(t *testing.T) {
		repo := fillRepo()
		afterCreated := posts[1].CreatedAt
		afterID := int64(2)
		limit := int32(10)

		got, err := repo.Get(ctx, &afterCreated, afterID, limit)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "Post 3", got[0].Title)
		assert.Equal(t, "Post 4", got[1].Title)
	})

	t.Run("empty_repo", func(t *testing.T) {
		repo := New()
		got, err := repo.Get(ctx, nil, 0, 10)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("returns_copies", func(t *testing.T) {
		repo := fillRepo()
		got, _ := repo.Get(ctx, nil, 0, 1)
		require.Len(t, got, 1)
		got[0].Title = "Changed"

		original, _ := repo.GetByID(ctx, 1)
		assert.Equal(t, "Post 1", original.Title)
	})
}
