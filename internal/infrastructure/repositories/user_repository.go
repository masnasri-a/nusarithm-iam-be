package repositories

import (
	"database/sql"
	"fmt"

	"backend/internal/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByID(id uuid.UUID) (*entities.User, error)
	GetByUsername(username string) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	GetByDomainID(domainID uuid.UUID) ([]*entities.User, error)
	Create(user *entities.User) error
	Update(user *entities.User) error
	UpdatePassword(id uuid.UUID, hashedPassword string) error
	Delete(id uuid.UUID) error
	ListWithPagination(search string, domainID uuid.UUID, page, limit int) (*UserListResult, error)
}

type UserListResult struct {
	Users      []*entities.User `json:"users"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow(`
		SELECT id, domain_id, role_id, first_name, last_name, username, email, password_hash, created_at, updated_at
		FROM users WHERE id = $1`, id).Scan(
		&user.ID, &user.DomainID, &user.RoleID, &user.FirstName, &user.LastName,
		&user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow(`
		SELECT id, domain_id, role_id, first_name, last_name, username, email, password_hash, created_at, updated_at
		FROM users WHERE username = $1`, username).Scan(
		&user.ID, &user.DomainID, &user.RoleID, &user.FirstName, &user.LastName,
		&user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.QueryRow(`
		SELECT id, domain_id, role_id, first_name, last_name, username, email, password_hash, created_at, updated_at
		FROM users WHERE email = $1`, email).Scan(
		&user.ID, &user.DomainID, &user.RoleID, &user.FirstName, &user.LastName,
		&user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByDomainID(domainID uuid.UUID) ([]*entities.User, error) {
	rows, err := r.db.Query(`
		SELECT id, domain_id, role_id, first_name, last_name, username, email, password_hash, created_at, updated_at
		FROM users WHERE domain_id = $1 ORDER BY username`, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.ID, &user.DomainID, &user.RoleID, &user.FirstName, &user.LastName,
			&user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepository) Create(user *entities.User) error {
	user.ID = uuid.New()
	err := r.db.QueryRow(`
		INSERT INTO users (id, domain_id, role_id, first_name, last_name, username, email, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		user.ID, user.DomainID, user.RoleID, user.FirstName, user.LastName,
		user.Username, user.Email, user.PasswordHash).Scan(&user.ID)
	return err
}

func (r *userRepository) Update(user *entities.User) error {
	_, err := r.db.Exec(`
		UPDATE users SET first_name = $1, last_name = $2, username = $3, email = $4, role_id = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6`, user.FirstName, user.LastName, user.Username, user.Email, user.RoleID, user.ID)
	return err
}

func (r *userRepository) UpdatePassword(id uuid.UUID, hashedPassword string) error {
	_, err := r.db.Exec(`
		UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`, hashedPassword, id)
	return err
}

func (r *userRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *userRepository) ListWithPagination(search string, domainID uuid.UUID, page, limit int) (*UserListResult, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Build the query with search condition
	baseQuery := "SELECT id, domain_id, role_id, first_name, last_name, username, email, password_hash, created_at, updated_at FROM users WHERE domain_id = $1"
	countQuery := "SELECT COUNT(*) FROM users WHERE domain_id = $1"
	args := []interface{}{domainID}
	var whereClause string

	if search != "" {
		whereClause = " AND (username ILIKE $" + fmt.Sprintf("%d", len(args)+1) +
			" OR email ILIKE $" + fmt.Sprintf("%d", len(args)+1) +
			" OR first_name ILIKE $" + fmt.Sprintf("%d", len(args)+1) +
			" OR last_name ILIKE $" + fmt.Sprintf("%d", len(args)+1) + ")"
		args = append(args, "%"+search+"%")
	}

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get paginated results
	query := baseQuery + whereClause + " ORDER BY username LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.ID, &user.DomainID, &user.RoleID, &user.FirstName, &user.LastName,
			&user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	return &UserListResult{
		Users:      users,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
