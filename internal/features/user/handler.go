package user

import (
	"github.com/gin-gonic/gin"

	"github.com/khaingminhtun/production-go-api/internal/shared/httpx"
	"github.com/khaingminhtun/production-go-api/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ============================================================
// Create
// ============================================================

// CreateUser handles:
//
// POST /users
//
// Creates a new user.
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.service.CreateUser(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, result)
}

// ============================================================
// Get
// ============================================================

// GetUser handles:
//
// GET /users/:id
//
// Returns a user by ID.
func (h *Handler) GetUser(c *gin.Context) {
	userID, err := httpx.ParamInt64(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	result, err := h.service.GetUser(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, result)
}

// ============================================================
// List
// ============================================================

// ListUsers handles:
//
// GET /users?offset=0&limit=20
//
// Returns a paginated list of users.
func (h *Handler) ListUsers(c *gin.Context) {
	offset, err := httpx.QueryInt(c, "offset", 0)
	if err != nil {
		response.BadRequest(c, "invalid offset")
		return
	}

	limit, err := httpx.QueryInt(c, "limit", 20)
	if err != nil {
		response.BadRequest(c, "invalid limit")
		return
	}

	result, err := h.service.ListUsers(
		c.Request.Context(),
		offset,
		limit,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, result)
}

// ============================================================
// Update
// ============================================================

// UpdateUser handles:
//
// PATCH /users/:id
//
// Updates a user.
func (h *Handler) UpdateUser(c *gin.Context) {
	userID, err := httpx.ParamInt64(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.service.UpdateUser(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, result)
}

// ============================================================
// Delete
// ============================================================

// DeleteUser handles:
//
// DELETE /users/:id
//
// Deletes a user.
func (h *Handler) DeleteUser(c *gin.Context) {
	userID, err := httpx.ParamInt64(c, "id")
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.service.DeleteUser(
		c.Request.Context(),
		userID,
	); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}
