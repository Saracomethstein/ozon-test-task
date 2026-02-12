package post

import "github.com/Saracomethstein/ozon-test-task/internal/service"

type postResolver struct {
	service *service.Container
}

func New(service *service.Container) *postResolver {
	return &postResolver{service}
}
