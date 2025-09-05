package services

import (
	"backend/internal/domain/entities"
	"backend/internal/infrastructure/repositories"

	"github.com/google/uuid"
)

type RoleService interface {
	GetRoleByID(id uuid.UUID) (*entities.Role, error)
	GetRolesByDomainID(domainID uuid.UUID) ([]*entities.Role, error)
	CreateRole(domainID uuid.UUID, roleName string, roleClaims map[string]interface{}) (*entities.Role, error)
	UpdateRole(id uuid.UUID, roleName string, roleClaims map[string]interface{}) (*entities.Role, error)
	DeleteRole(id uuid.UUID) error
	ListRolesWithPagination(search string, domainID uuid.UUID, page, limit int) (*repositories.RoleListResult, error)
}

type roleService struct {
	repo repositories.RoleRepository
}

func NewRoleService(repo repositories.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) GetRoleByID(id uuid.UUID) (*entities.Role, error) {
	return s.repo.GetByID(id)
}

func (s *roleService) GetRolesByDomainID(domainID uuid.UUID) ([]*entities.Role, error) {
	return s.repo.GetByDomainID(domainID)
}

func (s *roleService) CreateRole(domainID uuid.UUID, roleName string, roleClaims map[string]interface{}) (*entities.Role, error) {
	if roleClaims == nil {
		roleClaims = make(map[string]interface{})
	}

	role := &entities.Role{
		DomainID:   domainID,
		RoleName:   roleName,
		RoleClaims: roleClaims,
	}
	err := s.repo.Create(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleService) UpdateRole(id uuid.UUID, roleName string, roleClaims map[string]interface{}) (*entities.Role, error) {
	if roleClaims == nil {
		roleClaims = make(map[string]interface{})
	}

	role := &entities.Role{
		ID:         id,
		RoleName:   roleName,
		RoleClaims: roleClaims,
	}
	err := s.repo.Update(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (s *roleService) DeleteRole(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *roleService) ListRolesWithPagination(search string, domainID uuid.UUID, page, limit int) (*repositories.RoleListResult, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.repo.ListWithPagination(search, domainID, page, limit)
}
