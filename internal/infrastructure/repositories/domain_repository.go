package repositories

import (
	"database/sql"
	"fmt"

	"backend/internal/domain/entities"

	"github.com/google/uuid"
)

type DomainRepository interface {
	GetByID(id uuid.UUID) (*entities.Domain, error)
	Create(domain *entities.Domain) error
	List() ([]*entities.Domain, error)
	ListWithPagination(search string, page, limit int) (*DomainListResult, error)
	Update(domain *entities.Domain) error
	Delete(id uuid.UUID) error
}

type DomainListResult struct {
	Domains    []*entities.Domain `json:"domains"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

type domainRepository struct {
	db *sql.DB
}

func NewDomainRepository(db *sql.DB) DomainRepository {
	return &domainRepository{db: db}
}

func (r *domainRepository) GetByID(id uuid.UUID) (*entities.Domain, error) {
	var domain entities.Domain
	err := r.db.QueryRow("SELECT domain_id, name, domain FROM domains WHERE domain_id = $1", id).Scan(&domain.DomainID, &domain.Name, &domain.Domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

func (r *domainRepository) Create(domain *entities.Domain) error {
	domain.DomainID = uuid.New()
	err := r.db.QueryRow("INSERT INTO domains (domain_id, name, domain) VALUES ($1, $2, $3) RETURNING domain_id", domain.DomainID, domain.Name, domain.Domain).Scan(&domain.DomainID)
	return err
}

func (r *domainRepository) List() ([]*entities.Domain, error) {
	rows, err := r.db.Query("SELECT domain_id, name, domain FROM domains ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []*entities.Domain
	for rows.Next() {
		var domain entities.Domain
		err := rows.Scan(&domain.DomainID, &domain.Name, &domain.Domain)
		if err != nil {
			return nil, err
		}
		domains = append(domains, &domain)
	}
	return domains, nil
}

func (r *domainRepository) ListWithPagination(search string, page, limit int) (*DomainListResult, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Build the query with search condition
	baseQuery := "SELECT domain_id, name, domain FROM domains"
	countQuery := "SELECT COUNT(*) FROM domains"
	var args []interface{}
	var whereClause string

	if search != "" {
		whereClause = " WHERE name ILIKE $1 OR domain ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery+whereClause, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Get paginated results
	query := baseQuery + whereClause + " ORDER BY name LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []*entities.Domain
	for rows.Next() {
		var domain entities.Domain
		err := rows.Scan(&domain.DomainID, &domain.Name, &domain.Domain)
		if err != nil {
			return nil, err
		}
		domains = append(domains, &domain)
	}

	// Calculate total pages
	totalPages := (total + limit - 1) / limit

	return &DomainListResult{
		Domains:    domains,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *domainRepository) Update(domain *entities.Domain) error {
	_, err := r.db.Exec("UPDATE domains SET name = $1, domain = $2 WHERE domain_id = $3", domain.Name, domain.Domain, domain.DomainID)
	return err
}

func (r *domainRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM domains WHERE domain_id = $1", id)
	return err
}
