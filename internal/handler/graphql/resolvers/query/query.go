package query

import (
	"github.com/Saracomethstein/ozon-test-task/internal/service"
)

type queryResolver struct {
	service *service.Container
}

func New(service *service.Container) *queryResolver {
	return &queryResolver{service}
}
