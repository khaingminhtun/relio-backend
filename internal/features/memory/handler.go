package memory

import (
	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/httpx"
	"github.com/khaingminhtun/relio-backend/internal/shared/response"
)

type MemoryHandler struct {
	service MemoryService
}

func NewMemoryHandler(service MemoryService) *MemoryHandler {
	return &MemoryHandler{service: service}
}

// POST /relationships/:relationshipId/memories
func (h *MemoryHandler) Create(c *gin.Context) {
	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid relationship id",
			err,
		))
		return
	}

	var req CreateMemoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid request body",
			err,
		))
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	result, err := h.service.Create(
		c.Request.Context(),
		userID,
		relationshipID,
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, result)
}

// GET /relationships/:relationshipId/memories
func (h *MemoryHandler) List(c *gin.Context) {
	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid relationship id",
			err,
		))
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	memories, err := h.service.List(
		c.Request.Context(),
		userID,
		relationshipID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, memories)
}

// GET /relationships/:relationshipId/memories/:memoryId
func (h *MemoryHandler) Get(c *gin.Context) {
	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid relationship id",
			err,
		))
		return
	}

	memoryID, err := httpx.ParamInt64(c, "memoryId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid memory id",
			err,
		))
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	memory, err := h.service.GetByID(
		c.Request.Context(),
		userID,
		relationshipID,
		memoryID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, memory)
}

// PATCH /relationships/:relationshipId/memories/:memoryId
func (h *MemoryHandler) Update(c *gin.Context) {
	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid relationship id",
			err,
		))
		return
	}

	memoryID, err := httpx.ParamInt64(c, "memoryId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid memory id",
			err,
		))
		return
	}

	var req UpdateMemoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid request body",
			err,
		))
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	result, err := h.service.Update(
		c.Request.Context(),
		userID,
		relationshipID,
		memoryID,
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, result)
}

// DELETE /relationships/:relationshipId/memories/:memoryId
func (h *MemoryHandler) Delete(c *gin.Context) {
	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid relationship id",
			err,
		))
		return
	}

	memoryID, err := httpx.ParamInt64(c, "memoryId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid memory id",
			err,
		))
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Delete(
		c.Request.Context(),
		userID,
		relationshipID,
		memoryID,
	); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}
