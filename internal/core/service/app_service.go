package service

import (
	"setupwizard/internal/core/domain"
	"setupwizard/internal/core/ports"
)

type AppService struct {
	rep ports.AppRepository
}

func NewAppService(rep ports.AppRepository) ports.AppService {
	return &AppService{rep: rep}
}

func (s *AppService) GetApps() ([]*domain.App, error) {
	return s.rep.GetAll()
}
