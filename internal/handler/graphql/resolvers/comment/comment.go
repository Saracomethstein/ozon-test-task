package comment

import "github.com/Saracomethstein/ozon-test-task/internal/service"

type commentResolver struct {
	service *service.Container
}

func New(service *service.Container) *commentResolver {
	return &commentResolver{service}
}
