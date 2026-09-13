package ports

import "setupwizard/internal/core/domain"

type AppRepository interface {
	GetAll() ([]*domain.App, error)
}
