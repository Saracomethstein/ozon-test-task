package query

import (
	"context"
	"fmt"

	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
)

func (r *queryResolver) Posts(ctx context.Context, first *int32, after *string) (*graphql.PostConnection, error) {
	panic(fmt.Errorf("not implemented: Posts - posts"))
}
