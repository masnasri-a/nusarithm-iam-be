package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/application/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	RoleName   string                 `json:"role_name" binding:"required"`
	RoleClaims map[string]interface{} `json:"role_claims"`
}

type UpdateRoleRequest struct {
	RoleName   string                 `json:"role_name" binding:"required"`
	RoleClaims map[string]interface{} `json:"role_claims"`
}

type RoleHandler struct {
	roleService services.RoleService
}

func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// GetRole godoc
//
//	@Summary		Get a role
//	@Description	Get role by ID
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Role ID"
//	@Success		200	{object}	entities.Role
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Router			/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	role, err := h.roleService.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	c.JSON(http.StatusOK, role)
}

// GetRolesByDomain godoc
//
//	@Summary		Get roles by domain
//	@Description	Get all roles for a specific domain
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			domainId	path		string			true	"Domain ID"
//	@Success		200			{array}		entities.Role
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/domains/{domainId}/roles [get]
func (h *RoleHandler) GetRolesByDomain(c *gin.Context) {
	domainIdStr := c.Param("domainId")
	domainID, err := uuid.Parse(domainIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain UUID"})
		return
	}
	roles, err := h.roleService.GetRolesByDomainID(domainID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get roles"})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// ListRoles godoc
//
//	@Summary		List roles with pagination
//	@Description	Get roles with pagination and search
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			domainId	query		string	false	"Domain ID to filter roles"
//	@Param			search		query		string	false	"Search term for role name"
//	@Param			page		query		int		false	"Page number (default: 1)"
//	@Param			limit		query		int		false	"Items per page (default: 10, max: 100)"
//	@Success		200			{object}	repositories.RoleListResult
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	// Parse query parameters
	search := c.DefaultQuery("search", "")
	domainIdStr := c.DefaultQuery("domainId", "")
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

	var domainID uuid.UUID
	if domainIdStr != "" {
		domainID, err = uuid.Parse(domainIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain UUID"})
			return
		}
	}

	result, err := h.roleService.ListRolesWithPagination(search, domainID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list roles"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// CreateRole godoc
//
//	@Summary		Create a role
//	@Description	Create a new role
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			domainId	path		string					true	"Domain ID"
//	@Param			role		body		CreateRoleRequest		true	"Role data"
//	@Success		201			{object}	entities.Role
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/domains/{domainId}/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	domainIdStr := c.Param("domainId")
	domainID, err := uuid.Parse(domainIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid domain UUID"})
		return
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.roleService.CreateRole(domainID, req.RoleName, req.RoleClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}
	c.JSON(http.StatusCreated, role)
}

// UpdateRole godoc
//
//	@Summary		Update a role
//	@Description	Update role by ID
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Role ID"
//	@Param			role	body		UpdateRoleRequest		true	"Role data"
//	@Success		200		{object}	entities.Role
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.roleService.UpdateRole(id, req.RoleName, req.RoleClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}
	c.JSON(http.StatusOK, role)
}

// DeleteRole godoc
//
//	@Summary		Delete a role
//	@Description	Delete role by ID
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Role ID"
//	@Success		204	{object}	map[string]string
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	err = h.roleService.DeleteRole(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Role deleted successfully"})
}
