package query

import (
	"context"
	"fmt"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
)

func (r *queryResolver) Post(ctx context.Context, id string) (*graphql.Post, error) {
	panic(fmt.Errorf("not implemented: Post - post"))
}
