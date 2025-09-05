package services

import (
	"backend/internal/domain/entities"
	"backend/internal/infrastructure/repositories"

	"github.com/google/uuid"
)

type DomainService interface {
	GetDomainByID(id uuid.UUID) (*entities.Domain, error)
	CreateDomain(name, domainStr string) (*entities.Domain, error)
	ListDomains() ([]*entities.Domain, error)
	ListDomainsWithPagination(search string, page, limit int) (*repositories.DomainListResult, error)
	UpdateDomain(id uuid.UUID, name, domainStr string) (*entities.Domain, error)
	DeleteDomain(id uuid.UUID) error
}

type domainService struct {
	repo repositories.DomainRepository
}

func NewDomainService(repo repositories.DomainRepository) DomainService {
	return &domainService{repo: repo}
}

func (s *domainService) GetDomainByID(id uuid.UUID) (*entities.Domain, error) {
	return s.repo.GetByID(id)
}

func (s *domainService) CreateDomain(name, domainStr string) (*entities.Domain, error) {
	domain := &entities.Domain{
		Name:   name,
		Domain: domainStr,
	}
	err := s.repo.Create(domain)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (s *domainService) ListDomains() ([]*entities.Domain, error) {
	return s.repo.List()
}

func (s *domainService) ListDomainsWithPagination(search string, page, limit int) (*repositories.DomainListResult, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.repo.ListWithPagination(search, page, limit)
}

func (s *domainService) UpdateDomain(id uuid.UUID, name, domainStr string) (*entities.Domain, error) {
	domain := &entities.Domain{
		DomainID: id,
		Name:     name,
		Domain:   domainStr,
	}
	err := s.repo.Update(domain)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

func (s *domainService) DeleteDomain(id uuid.UUID) error {
	return s.repo.Delete(id)
}
