package subscription

import (
	"context"
	"fmt"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
)

func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *graphql.Comment, error) {
	panic(fmt.Errorf("not implemented: CommentAdded - commentAdded"))
}
