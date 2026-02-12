package post

import (
	"context"
	"fmt"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
)

func (r *postResolver) Comments(ctx context.Context, obj *graphql.Post, first *int32, after *string) (*graphql.CommentConnection, error) {
	panic(fmt.Errorf("not implemented: Comments - comments"))
}
