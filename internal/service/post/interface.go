package post

import (
	"context"

	"github.com/Saracomethstein/ozon-test-task/internal/models"
)

type UseCase interface {
	// TODO: rename func for package 'post'

	CreatePost(ctx context.Context, in models.CreatePostInput) (*models.Post, error)
	SetPostCommentsAllowed(ctx context.Context, postID string, allow bool) (*models.Post, error)
	GetPostById(ctx context.Context, postID string) (*models.Post, error)
	GetPosts(ctx context.Context, first *int32, after *string) (*models.PostConnection, error)
}
