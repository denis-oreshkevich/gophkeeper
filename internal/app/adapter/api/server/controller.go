package server

import (
	"github.com/denis-oreshkevich/gophkeeper/internal/app/domain/service"
)

type Controller struct {
	crudSvc    service.CRUDService
	syncSvc    service.SyncService
	userClient service.UserClientService
}

func NewController(crudSvc service.CRUDService, syncSvc service.SyncService,
	userClient service.UserClientService,
) *Controller {
	return &Controller{
		crudSvc:    crudSvc,
		syncSvc:    syncSvc,
		userClient: userClient,
	}
}
