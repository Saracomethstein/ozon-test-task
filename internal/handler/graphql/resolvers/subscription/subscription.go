package subscription

import (
	"github.com/Saracomethstein/ozon-test-task/internal/service"
)

type subscriptionResolver struct {
	service *service.Container
}

func New(service *service.Container) *subscriptionResolver {
	return &subscriptionResolver{service}
}
