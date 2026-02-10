package mutation

import (
	"context"
	"fmt"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
)

func (r *mutationResolver) SetPostCommentsAllowed(ctx context.Context, postID string, allow bool) (*graphql.Post, error) {
	panic(fmt.Errorf("not implemented: SetPostCommentsAllowed - setPostCommentsAllowed"))
}
