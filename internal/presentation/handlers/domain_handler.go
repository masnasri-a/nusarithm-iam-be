package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/application/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateDomainRequest struct {
	Name   string `json:"name" binding:"required"`
	Domain string `json:"domain" binding:"required"`
}

type UpdateDomainRequest struct {
	Name   string `json:"name" binding:"required"`
	Domain string `json:"domain" binding:"required"`
}

type DomainHandler struct {
	domainService services.DomainService
}

func NewDomainHandler(domainService services.DomainService) *DomainHandler {
	return &DomainHandler{domainService: domainService}
}

// GetDomain godoc
//
//	@Summary		Get a domain
//	@Description	Get domain by ID
//	@Tags			domains
//	@Accept			json
//	@Produce		json
//	@Param			domainId	path		string			true	"Domain ID"
//	@Success		200	{object}	entities.Domain
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/domains/{domainId} [get]
func (h *DomainHandler) GetDomain(c *gin.Context) {
	idStr := c.Param("domainId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	domain, err := h.domainService.GetDomainByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}
	c.JSON(http.StatusOK, domain)
}

// CreateDomain godoc
//
//	@Summary		Create a domain
//	@Description	Create a new domain
//	@Tags			domains
//	@Accept			json
//	@Produce		json
//	@Param			domain	body		CreateDomainRequest	true	"Domain data"
//	@Success		201		{object}	entities.Domain
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/domains [post]
func (h *DomainHandler) CreateDomain(c *gin.Context) {
	var req CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	domain, err := h.domainService.CreateDomain(req.Name, req.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create domain"})
		return
	}
	c.JSON(http.StatusCreated, domain)
}

// ListDomains godoc
//
//	@Summary		List all domains
//	@Description	Get all domains with pagination and search
//	@Tags			domains
//	@Accept			json
//	@Produce		json
//	@Param			search	query		string	false	"Search term for domain name"
//	@Param			page	query		int		false	"Page number (default: 1)"
//	@Param			limit	query		int		false	"Items per page (default: 10, max: 100)"
//	@Success		200		{object}	repositories.DomainListResult
//	@Failure		500		{object}	map[string]string
//	@Router			/domains [get]
func (h *DomainHandler) ListDomains(c *gin.Context) {
	// Parse query parameters
	search := c.DefaultQuery("search", "")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	result, err := h.domainService.ListDomainsWithPagination(search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list domains"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// UpdateDomain godoc
//
//	@Summary		Update a domain
//	@Description	Update domain by ID
//	@Tags			domains
//	@Accept			json
//	@Produce		json
//	@Param			domainId	path		string					true	"Domain ID"
//	@Param			domain	body		UpdateDomainRequest	true	"Domain data"
//	@Success		200		{object}	entities.Domain
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/domains/{domainId} [put]
func (h *DomainHandler) UpdateDomain(c *gin.Context) {
	idStr := c.Param("domainId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var req UpdateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domain, err := h.domainService.UpdateDomain(id, req.Name, req.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update domain"})
		return
	}
	c.JSON(http.StatusOK, domain)
}

// DeleteDomain godoc
//
//	@Summary		Delete a domain
//	@Description	Delete domain by ID
//	@Tags			domains
//	@Accept			json
//	@Produce		json
//	@Param			domainId	path		string			true	"Domain ID"
//	@Success		204	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/domains/{domainId} [delete]
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	idStr := c.Param("domainId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	err = h.domainService.DeleteDomain(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete domain"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Domain deleted successfully"})
}
