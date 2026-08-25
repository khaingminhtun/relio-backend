package relationship

import (
	"github.com/gin-gonic/gin"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/httpx"
	"github.com/khaingminhtun/relio-backend/internal/shared/response"
)

type RelationshipMemberHandler struct {
	service RelationshipMemberService
}

func NewRelationshipMemberHandler(
	service RelationshipMemberService,
) *RelationshipMemberHandler {
	return &RelationshipMemberHandler{
		service: service,
	}
}

// GET /:relationshipId/members
func (h *RelationshipMemberHandler) List(c *gin.Context) {
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

	members, err := h.service.ListByRelationshipID(
		c.Request.Context(),
		relationshipID,
		userID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, members)
}

// GET /:relationshipId/members/:memberId
func (h *RelationshipMemberHandler) GetByID(c *gin.Context) {

	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid relationship id",
			err,
		))
		return
	}
	memberID, err := httpx.ParamInt64(c, "memberId")
	if err != nil {
		c.Error(err)
		return
	}

	member, err := h.service.GetByID(
		c.Request.Context(),
		relationshipID,
		memberID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, member)
}

// GET /:relationshipId/members/me
func (h *RelationshipMemberHandler) GetMe(c *gin.Context) {
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
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid user id",
			err,
		))
		return
	}

	member, err := h.service.GetByRelationshipAndUser(
		c.Request.Context(),
		relationshipID,
		userID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, member)
}

// PATCH /members/:memberId
func (h *RelationshipMemberHandler) UpdateMember(c *gin.Context) {
	relationshipID, err := httpx.ParamInt64(c, "relationshipId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid relationship id",
			err,
		))
		return
	}
	var req UpdateMemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid request body",
			err,
		))
		return
	}

	memberID, err := httpx.ParamInt64(c, "memberId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid member id",
			err,
		))
		return
	}

	member, err := h.service.UpdateMember(
		c.Request.Context(),
		relationshipID,
		memberID,
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, member)
}

// DELETE /:relationshipId/members/:userId
func (h *RelationshipMemberHandler) RemoveMember(c *gin.Context) {
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
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid user id",
			err,
		))
		return
	}

	if err := h.service.RemoveMember(
		c.Request.Context(),
		relationshipID,
		userID,
	); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}

// POST /:relationshipId/members/leave
func (h *RelationshipMemberHandler) Leave(c *gin.Context) {
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

	if err := h.service.RemoveMember(
		c.Request.Context(),
		relationshipID,
		userID,
	); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}
