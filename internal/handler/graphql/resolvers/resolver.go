package resolvers

import (
	"github.com/Saracomethstein/ozon-test-task/generated/graphql"
	"github.com/Saracomethstein/ozon-test-task/internal/handler/graphql/resolvers/comment"
	"github.com/Saracomethstein/ozon-test-task/internal/handler/graphql/resolvers/mutation"
	"github.com/Saracomethstein/ozon-test-task/internal/handler/graphql/resolvers/post"
	"github.com/Saracomethstein/ozon-test-task/internal/handler/graphql/resolvers/query"
	"github.com/Saracomethstein/ozon-test-task/internal/handler/graphql/resolvers/subscription"
	"github.com/Saracomethstein/ozon-test-task/internal/service"
)

type Resolver struct {
	service *service.Container
}

func New(services *service.Container) *Resolver {
	return &Resolver{service: services}
}

func (r *Resolver) Query() graphql.QueryResolver {
	return query.New(r.service)
}

func (r *Resolver) Mutation() graphql.MutationResolver {
	return mutation.New(r.service)
}

func (r *Resolver) Subscription() graphql.SubscriptionResolver {
	return subscription.New(r.service)
}

func (r *Resolver) Post() graphql.PostResolver {
	return post.New(r.service)
}

func (r *Resolver) Comment() graphql.CommentResolver {
	return comment.New(r.service)
}
