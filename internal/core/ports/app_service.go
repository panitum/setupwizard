package ports

import "setupwizard/internal/core/domain"

type AppService interface {
	GetApps() ([]*domain.App, error)
}
