package api

import (
	"github.com/denis-oreshkevich/gophkeeper/internal/server/service"
)

type Controller struct {
	svc service.ServerService
}

func NewController(svc service.ServerService) *Controller {
	return &Controller{
		svc: svc,
	}
}
