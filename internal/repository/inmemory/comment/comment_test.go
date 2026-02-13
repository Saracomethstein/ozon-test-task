package comment

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
	"github.com/Saracomethstein/ozon-test-task/internal/repository"
	"github.com/Saracomethstein/ozon-test-task/internal/repository/inmemory/post"
)

func TestCommentRepo_Add(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	now := time.Now().UTC()

	t.Run("successful_add_root_comment", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)

		comment := models.Comment{
			PostID:    postID,
			ParentID:  nil,
			Author:    "Alice",
			Text:      "Root comment",
			CreatedAt: now.Format(time.RFC3339),
		}

		got, err := repo.Add(ctx, comment)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), got.ID)
		assert.Equal(t, postID, got.PostID)
		assert.Nil(t, got.ParentID)
		assert.Equal(t, "Alice", got.Author)
		assert.Equal(t, "Root comment", got.Text)
		assert.Equal(t, now.Format(time.RFC3339), got.CreatedAt)

		total, _ := repo.TotalCount(ctx, postID)
		assert.Equal(t, int64(1), total)

		roots, _ := repo.GetRootByPost(ctx, postID, nil, 0, 10)
		require.Len(t, roots, 1)
		assert.Equal(t, got.ID, roots[0].ID)
	})

	t.Run("successful_add_child_comment", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)
		parent := addComment(t, repo, postID, nil, "Parent", "Parent text", now)

		child := models.Comment{
			PostID:    postID,
			ParentID:  &parent.ID,
			Author:    "Bob",
			Text:      "Child comment",
			CreatedAt: now.Add(time.Hour).Format(time.RFC3339),
		}

		got, err := repo.Add(ctx, child)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), got.ID)
		assert.Equal(t, parent.ID, *got.ParentID)

		children, _ := repo.GetChild(ctx, parent.ID, nil, 0, 10)
		require.Len(t, children, 1)
		assert.Equal(t, got.ID, children[0].ID)
	})

	t.Run("multiple_adds_increment_id", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)

		c1 := addComment(t, repo, postID, nil, "A", "text1", now)
		c2 := addComment(t, repo, postID, nil, "B", "text2", now.Add(time.Hour))

		assert.Equal(t, int64(1), c1.ID)
		assert.Equal(t, int64(2), c2.ID)
	})
}

func TestCommentRepo_CheckAllowComments(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("post_allows_comments", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)

		allow, err := repo.CheckAllowComments(ctx, postID)

		assert.NoError(t, err)
		assert.True(t, allow)
	})

	t.Run("post_disallows_comments", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, false)

		allow, err := repo.CheckAllowComments(ctx, postID)

		assert.NoError(t, err)
		assert.False(t, allow)
	})

	t.Run("post_not_found", func(t *testing.T) {
		repo, _ := setupCommentRepo(t)

		allow, err := repo.CheckAllowComments(ctx, 999)

		assert.Error(t, err)
		assert.EqualError(t, err, "post not found")
		assert.False(t, allow)
	})
}

func TestCommentRepo_CheckParentExists(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	now := time.Now().UTC()

	t.Run("parent_exists", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)
		parent := addComment(t, repo, postID, nil, "Parent", "text", now)

		gotPostID, err := repo.CheckParentExists(ctx, parent.ID)

		assert.NoError(t, err)
		assert.Equal(t, postID, gotPostID)
	})

	t.Run("parent_not_found", func(t *testing.T) {
		repo, _ := setupCommentRepo(t)

		gotPostID, err := repo.CheckParentExists(ctx, 999)

		assert.Error(t, err)
		assert.EqualError(t, err, "parent comment not found")
		assert.Equal(t, int64(0), gotPostID)
	})
}

func TestCommentRepo_TotalCount(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	now := time.Now().UTC()

	t.Run("no_comments_for_post", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)

		count, err := repo.TotalCount(ctx, postID)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("with_comments", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)
		addComment(t, repo, postID, nil, "A", "text1", now)
		addComment(t, repo, postID, nil, "B", "text2", now.Add(time.Hour))

		count, err := repo.TotalCount(ctx, postID)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})
}

func TestCommentRepo_GetRootByPost(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	fixedTime, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z")

	setupWithRoots := func(t *testing.T) (repository.CommentUC, int64) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)

		addComment(t, repo, postID, nil, "Root1", "text1", fixedTime)
		addComment(t, repo, postID, nil, "Root2", "text2", fixedTime.Add(-1*time.Hour))
		addComment(t, repo, postID, nil, "Root3", "text3", fixedTime.Add(-2*time.Hour))
		parent := addComment(t, repo, postID, nil, "Parent", "parent", fixedTime.Add(-3*time.Hour))
		addComment(t, repo, postID, &parent.ID, "Child", "child", fixedTime.Add(-4*time.Hour))

		return repo, postID
	}

	t.Run("first_page_without_cursor", func(t *testing.T) {
		repo, postID := setupWithRoots(t)
		limit := int32(2)

		got, err := repo.GetRootByPost(ctx, postID, nil, 0, limit)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, int64(1), got[0].ID)
		assert.Equal(t, "Root1", got[0].Author)
		assert.Equal(t, fixedTime.Format(time.RFC3339), got[0].CreatedAt)
		assert.Equal(t, int64(2), got[1].ID)
		assert.Equal(t, "Root2", got[1].Author)
		assert.Equal(t, fixedTime.Add(-1*time.Hour).Format(time.RFC3339), got[1].CreatedAt)
	})

	t.Run("second_page_with_cursor", func(t *testing.T) {
		repo, postID := setupWithRoots(t)
		afterCreatedAt := fixedTime.Add(-1 * time.Hour).Format(time.RFC3339)
		afterID := int64(2)
		limit := int32(2)

		got, err := repo.GetRootByPost(ctx, postID, &afterCreatedAt, afterID, limit)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, int64(3), got[0].ID)
		assert.Equal(t, "Root3", got[0].Author)
		assert.Equal(t, fixedTime.Add(-2*time.Hour).Format(time.RFC3339), got[0].CreatedAt)
	})

	t.Run("post_with_no comments", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)

		got, err := repo.GetRootByPost(ctx, postID, nil, 0, 10)

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("returns_copies", func(t *testing.T) {
		repo, postID := setupWithRoots(t)
		got, _ := repo.GetRootByPost(ctx, postID, nil, 0, 1)
		require.Len(t, got, 1)
		got[0].Author = "Hacked"

		original, _ := repo.GetRootByPost(ctx, postID, nil, 0, 10)
		assert.Equal(t, "Root1", original[0].Author)
	})
}

func TestCommentRepo_GetChild(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	now := time.Now().UTC()

	setupWithChildren := func(t *testing.T) (repository.CommentUC, int64, int64) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)

		parent := addComment(t, repo, postID, nil, "Parent", "parent", now.Add(-1*time.Hour))

		addComment(t, repo, postID, &parent.ID, "Child1", "text1", now)
		addComment(t, repo, postID, &parent.ID, "Child2", "text2", now.Add(-30*time.Minute))
		addComment(t, repo, postID, &parent.ID, "Child3", "text3", now.Add(-1*time.Hour))

		return repo, postID, parent.ID
	}

	t.Run("first_page_without_cursor", func(t *testing.T) {
		repo, _, parentID := setupWithChildren(t)
		limit := int32(2)

		got, err := repo.GetChild(ctx, parentID, nil, 0, limit)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, int64(2), got[0].ID)
		assert.Equal(t, int64(3), got[1].ID)
	})

	t.Run("second_page_with_cursor", func(t *testing.T) {
		repo, _, parentID := setupWithChildren(t)
		afterCreatedAt := now.Add(-30 * time.Minute).Format(time.RFC3339)
		afterID := int64(3)
		limit := int32(2)

		got, err := repo.GetChild(ctx, parentID, &afterCreatedAt, afterID, limit)

		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, int64(4), got[0].ID)
	})

	t.Run("parent_with_no_children", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)
		parent := addComment(t, repo, postID, nil, "Parent", "text", now)

		got, err := repo.GetChild(ctx, parent.ID, nil, 0, 10)

		require.NoError(t, err)
		assert.Empty(t, got)
	})
}

func TestCommentRepo_GetChildBatch(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	now := time.Now().UTC()

	setupForBatch := func(t *testing.T) (repository.CommentUC, []int64) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)

		parent1 := addComment(t, repo, postID, nil, "Parent1", "p1", now.Add(-2*time.Hour))
		parent2 := addComment(t, repo, postID, nil, "Parent2", "p2", now.Add(-3*time.Hour))

		addComment(t, repo, postID, &parent1.ID, "Child1A", "a", now)
		addComment(t, repo, postID, &parent1.ID, "Child1B", "b", now.Add(-30*time.Minute))
		addComment(t, repo, postID, &parent2.ID, "Child2A", "c", now.Add(-1*time.Hour))

		return repo, []int64{parent1.ID, parent2.ID}
	}

	t.Run("successful_batch", func(t *testing.T) {
		repo, parentIDs := setupForBatch(t)

		got, err := repo.GetChildBatch(ctx, parentIDs)

		require.NoError(t, err)
		require.Len(t, got, 3)

		assert.Equal(t, int64(3), got[0].ID)
		assert.Equal(t, int64(4), got[1].ID)
		assert.Equal(t, int64(5), got[2].ID)

		assert.Equal(t, parentIDs[0], *got[0].ParentID)
		assert.Equal(t, parentIDs[0], *got[1].ParentID)
		assert.Equal(t, parentIDs[1], *got[2].ParentID)
	})

	t.Run("some_parents_have_no_children", func(t *testing.T) {
		repo, postRepo := setupCommentRepo(t)
		postID := createTestPost(t, postRepo, true)
		parent1 := addComment(t, repo, postID, nil, "Parent1", "p1", now)
		parent2 := addComment(t, repo, postID, nil, "Parent2", "p2", now)

		got, err := repo.GetChildBatch(ctx, []int64{parent1.ID, parent2.ID})

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("empty_parentIDs_slice", func(t *testing.T) {
		repo, _ := setupCommentRepo(t)

		got, err := repo.GetChildBatch(ctx, []int64{})

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("returns_copies", func(t *testing.T) {
		repo, parentIDs := setupForBatch(t)
		got, _ := repo.GetChildBatch(ctx, parentIDs)
		require.NotEmpty(t, got)
		got[0].Author = "Hacked"

		original, _ := repo.GetChild(ctx, parentIDs[0], nil, 0, 10)
		assert.Equal(t, "Child1A", original[0].Author)
	})
}

func setupCommentRepo(t *testing.T) (repository.CommentUC, repository.PostUC) {
	t.Helper()
	postRepo := post.New()
	commentRepo := New(postRepo)
	return commentRepo, postRepo
}

func createTestPost(t *testing.T, postRepo repository.PostUC, allowComments bool) int64 {
	t.Helper()
	p, err := postRepo.Save(context.Background(), models.Post{
		Title:         "Test Post",
		Body:          "Content",
		Author:        "Tester",
		AllowComments: allowComments,
		CreatedAt:     time.Now().Format(time.RFC3339),
	})
	require.NoError(t, err)
	return p.ID
}

func addComment(t *testing.T, repo repository.CommentUC, postID int64, parentID *int64, author, text string, createdAt time.Time) *models.Comment {
	t.Helper()
	c, err := repo.Add(context.Background(), models.Comment{
		PostID:    postID,
		ParentID:  parentID,
		Author:    author,
		Text:      text,
		CreatedAt: createdAt.Format(time.RFC3339),
	})
	require.NoError(t, err)
	return c
}
