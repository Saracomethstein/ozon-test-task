package mutation

import (
	"context"
	"fmt"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
)

func (r *mutationResolver) AddComment(ctx context.Context, input graphql.AddCommentInput) (*graphql.Comment, error) {
	panic(fmt.Errorf("not implemented: AddComment - addComment"))
}
