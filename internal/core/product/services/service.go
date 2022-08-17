package services

import (
	"losh/internal/core/product/models"
)

type Service struct {
	repo Repository

	// used to cache licenses
	licenses map[string]*models.License
	nameToID map[string]string

	// used to cache struct fields to speed up copying
	// structFieldsCache map[reflect.Type]map[string]structField
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		// structFieldsCache: make(map[reflect.Type]map[string]structField),
	}
}
