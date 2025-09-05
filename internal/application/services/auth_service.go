package services

import (
	"crypto/sha256"
	"fmt"
	"time"

	"backend/internal/domain/entities"
	"backend/internal/infrastructure/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(domainID uuid.UUID, username, password string) (*LoginResponse, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	GetProfile(userID uuid.UUID) (*UserProfile, error)
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        *UserProfile `json:"user"`
}

type UserProfile struct {
	ID        uuid.UUID      `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      *RoleProfile   `json:"role"`
	Domain    *DomainProfile `json:"domain"`
}

type RoleProfile struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Claims      map[string]interface{} `json:"claims"`
}

type DomainProfile struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type TokenClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	DomainID uuid.UUID `json:"domain_id"`
	Username string    `json:"username"`
	RoleID   uuid.UUID `json:"role_id"`
	jwt.RegisteredClaims
}

type authService struct {
	userRepo    repositories.UserRepository
	roleRepo    repositories.RoleRepository
	domainRepo  repositories.DomainRepository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

func NewAuthService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository, domainRepo repositories.DomainRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		domainRepo:  domainRepo,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: 24 * time.Hour, // 24 hours
	}
}

func (s *authService) Login(domainID uuid.UUID, username, password string) (*LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user belongs to the specified domain
	if user.DomainID != domainID {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password
	if !s.verifyPassword(user.PasswordHash, password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Get user profile with role and domain
	userProfile, err := s.buildUserProfile(user)
	if err != nil {
		return nil, fmt.Errorf("failed to build user profile: %w", err)
	}

	return &LoginResponse{
		AccessToken: token,
		User:        userProfile,
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

func (s *authService) GetProfile(userID uuid.UUID) (*UserProfile, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return s.buildUserProfile(user)
}

func (s *authService) generateToken(user *entities.User) (string, error) {
	claims := TokenClaims{
		UserID:   user.ID,
		DomainID: user.DomainID,
		Username: user.Username,
		RoleID:   user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "nusarithm-iam",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *authService) verifyPassword(hashedPassword, password string) bool {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash) == hashedPassword
}

func (s *authService) buildUserProfile(user *entities.User) (*UserProfile, error) {
	// Get role information
	role, err := s.roleRepo.GetByID(user.RoleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Get domain information
	domain, err := s.domainRepo.GetByID(user.DomainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain: %w", err)
	}

	return &UserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role: &RoleProfile{
			ID:          role.ID,
			Name:        role.RoleName,
			Description: "", // Role doesn't have description field
			Claims:      role.RoleClaims,
		},
		Domain: &DomainProfile{
			ID:          domain.DomainID,
			Name:        domain.Name,
			Description: domain.Domain, // Using domain field as description
		},
	}, nil
}
