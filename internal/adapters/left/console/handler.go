package console

import (
	"setupwizard/internal/core/domain"
	"setupwizard/internal/core/ports"
)

type Handler struct {
	appService ports.AppService
}

func NewHandler(appService ports.AppService) *Handler {
	return &Handler{
		appService: appService,
	}
}

func (h *Handler) GetApps() ([]*domain.App, error) {
	return h.appService.GetApps()
}
