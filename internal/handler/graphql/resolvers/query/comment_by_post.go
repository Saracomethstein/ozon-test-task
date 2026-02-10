package query

import (
	"context"
	"fmt"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
)

func (r *queryResolver) CommentsByPost(ctx context.Context, postID string, first *int32, after *string) (*graphql.CommentConnection, error) {
	panic(fmt.Errorf("not implemented: CommentsByPost - commentsByPost"))
}
