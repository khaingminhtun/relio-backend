package relationship

import (
	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/httpx"
	"github.com/khaingminhtun/relio-backend/internal/shared/response"
)

type RelationshipHandler struct {
	service RelationshipService
}

func NewHandler(service RelationshipService) *RelationshipHandler {
	return &RelationshipHandler{service: service}
}

func (h *RelationshipHandler) CreateRelationship(c *gin.Context) {
	var req CreateRelationshipRequest

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
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, result)
}

func (h *RelationshipHandler) List(c *gin.Context) {
	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	relationships, err := h.service.List(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, relationships)
}

func (h *RelationshipHandler) GetByID(c *gin.Context) {
	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(err)
		return
	}

	result, err := h.service.GetByID(
		c.Request.Context(),
		relationshipID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, result)
}

func (h *RelationshipHandler) Update(c *gin.Context) {
	var req UpdateRelationshipRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid request body",
			err,
		))
		return
	}

	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(err)
		return
	}

	result, err := h.service.Update(
		c.Request.Context(),
		relationshipID,
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, result)
}

func (h *RelationshipHandler) Delete(c *gin.Context) {
	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.Delete(
		c.Request.Context(),
		relationshipID,
	); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}
