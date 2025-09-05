package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/internal/domain/entities"

	"github.com/google/uuid"
)

type RoleRepository interface {
	GetByID(id uuid.UUID) (*entities.Role, error)
	GetByDomainID(domainID uuid.UUID) ([]*entities.Role, error)
	Create(role *entities.Role) error
	Update(role *entities.Role) error
	Delete(id uuid.UUID) error
	ListWithPagination(search string, domainID uuid.UUID, page, limit int) (*RoleListResult, error)
}

type RoleListResult struct {
	Roles      []*entities.Role `json:"roles"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByID(id uuid.UUID) (*entities.Role, error) {
	var role entities.Role
	var claimsJSON []byte

	err := r.db.QueryRow(`
		SELECT id, domain_id, role_name, role_claims, created_at, updated_at
		FROM roles WHERE id = $1`, id).Scan(
		&role.ID, &role.DomainID, &role.RoleName, &claimsJSON, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Parse JSONB claims
	if err := json.Unmarshal(claimsJSON, &role.RoleClaims); err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) GetByDomainID(domainID uuid.UUID) ([]*entities.Role, error) {
	rows, err := r.db.Query(`
		SELECT id, domain_id, role_name, role_claims, created_at, updated_at
		FROM roles WHERE domain_id = $1 ORDER BY role_name`, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entities.Role
	for rows.Next() {
		var role entities.Role
		var claimsJSON []byte

		err := rows.Scan(&role.ID, &role.DomainID, &role.RoleName, &claimsJSON, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Parse JSONB claims
		if err := json.Unmarshal(claimsJSON, &role.RoleClaims); err != nil {
			return nil, err
		}

		roles = append(roles, &role)
	}
	return roles, nil
}

func (r *roleRepository) Create(role *entities.Role) error {
	role.ID = uuid.New()

	// Convert claims to JSON
	claimsJSON, err := json.Marshal(role.RoleClaims)
	if err != nil {
		return err
	}

	err = r.db.QueryRow(`
		INSERT INTO roles (id, domain_id, role_name, role_claims)
		VALUES ($1, $2, $3, $4) RETURNING id`,
		role.ID, role.DomainID, role.RoleName, claimsJSON).Scan(&role.ID)
	return err
}

func (r *roleRepository) Update(role *entities.Role) error {
	// Convert claims to JSON
	claimsJSON, err := json.Marshal(role.RoleClaims)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		UPDATE roles SET role_name = $1, role_claims = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3`, role.RoleName, claimsJSON, role.ID)
	return err
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM roles WHERE id = $1", id)
	return err
}

func (r *roleRepository) ListWithPagination(search string, domainID uuid.UUID, page, limit int) (*RoleListResult, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Build the query with search condition
	baseQuery := "SELECT id, domain_id, role_name, role_claims, created_at, updated_at FROM roles WHERE domain_id = $1"
	countQuery := "SELECT COUNT(*) FROM roles WHERE domain_id = $1"
	args := []interface{}{domainID}
	var whereClause string

	if search != "" {
		whereClause = " AND role_name ILIKE $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, "%"+search+"%")
	}

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get paginated results
	query := baseQuery + whereClause + " ORDER BY role_name LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entities.Role
	for rows.Next() {
		var role entities.Role
		var claimsJSON []byte

		err := rows.Scan(&role.ID, &role.DomainID, &role.RoleName, &claimsJSON, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Parse JSONB claims
		if err := json.Unmarshal(claimsJSON, &role.RoleClaims); err != nil {
			return nil, err
		}

		roles = append(roles, &role)
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	return &RoleListResult{
		Roles:      roles,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}
