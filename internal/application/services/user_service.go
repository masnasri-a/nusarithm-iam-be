package services

import (
	"crypto/sha256"
	"fmt"

	"backend/internal/domain/entities"
	"backend/internal/infrastructure/repositories"

	"github.com/google/uuid"
)

type UserService interface {
	GetUserByID(id uuid.UUID) (*entities.User, error)
	GetUserByUsername(username string) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	GetUsersByDomainID(domainID uuid.UUID) ([]*entities.User, error)
	CreateUser(domainID, roleID uuid.UUID, firstName, lastName, username, email, password string) (*entities.User, error)
	UpdateUser(id uuid.UUID, firstName, lastName, username, email string, roleID uuid.UUID) (*entities.User, error)
	ResetUserPassword(id uuid.UUID, newPassword string) error
	DeleteUser(id uuid.UUID) error
	ListUsersWithPagination(search string, domainID uuid.UUID, page, limit int) (*repositories.UserListResult, error)
	VerifyPassword(hashedPassword, password string) bool
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(id uuid.UUID) (*entities.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetUserByUsername(username string) (*entities.User, error) {
	return s.repo.GetByUsername(username)
}

func (s *userService) GetUserByEmail(email string) (*entities.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) GetUsersByDomainID(domainID uuid.UUID) ([]*entities.User, error) {
	return s.repo.GetByDomainID(domainID)
}

func (s *userService) CreateUser(domainID, roleID uuid.UUID, firstName, lastName, username, email, password string) (*entities.User, error) {
	// Hash the password
	hashedPassword := s.hashPassword(password)

	user := &entities.User{
		DomainID:     domainID,
		RoleID:       roleID,
		FirstName:    firstName,
		LastName:     lastName,
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}
	err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateUser(id uuid.UUID, firstName, lastName, username, email string, roleID uuid.UUID) (*entities.User, error) {
	user := &entities.User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Email:     email,
		RoleID:    roleID,
	}
	err := s.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) ResetUserPassword(id uuid.UUID, newPassword string) error {
	// Hash the new password
	hashedPassword := s.hashPassword(newPassword)

	// Update the user's password hash
	return s.repo.UpdatePassword(id, hashedPassword)
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *userService) ListUsersWithPagination(search string, domainID uuid.UUID, page, limit int) (*repositories.UserListResult, error) {
	// Set default values
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	return s.repo.ListWithPagination(search, domainID, page, limit)
}

func (s *userService) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

func (s *userService) VerifyPassword(hashedPassword, password string) bool {
	return s.hashPassword(password) == hashedPassword
}
