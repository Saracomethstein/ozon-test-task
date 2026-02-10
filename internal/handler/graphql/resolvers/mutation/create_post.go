package mutation

import (
	"context"
	"fmt"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input graphql.CreatePostInput) (*graphql.Post, error) {
	panic(fmt.Errorf("not implemented: CreatePost - createPost"))
}
