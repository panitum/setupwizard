package json

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"setupwizard/internal/core/domain"
)

type JsonRepository struct {
	storage []domain.App
}

func NewJsonRepository() *JsonRepository {
	return &JsonRepository{storage: []domain.App{}}
}

func (rep *JsonRepository) GetAll() ([]*domain.App, error) {
	if err := rep.loadApps(); err != nil {
		return []*domain.App{}, err
	}

	if len(rep.storage) == 0 {
		return []*domain.App{}, errors.New("apps not found")
	}

	apps := make([]*domain.App, 0, len(rep.storage))
	for i := range rep.storage {
		apps = append(apps, &rep.storage[i])
	}

	return apps, nil
}

func (rep *JsonRepository) loadApps() error {
	var req struct {
		Apps []domain.App `json:"apps"`
	}

	data, err := os.ReadFile(getFullPath())
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	rep.storage = req.Apps

	return nil
}

func getFullPath() string {
	return filepath.Join("config", "config.json")
}
