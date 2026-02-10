package mutation

import (
	"github.com/Saracomethstein/ozon-test-task/internal/service"
)

type mutationResolver struct {
	service *service.Container
}

func New(service *service.Container) *mutationResolver {
	return &mutationResolver{service}
}
