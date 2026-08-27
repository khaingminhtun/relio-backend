package relationship

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/relio-backend/internal/shared/httpx"
	"github.com/khaingminhtun/relio-backend/internal/shared/response"
)

type InvitationHandler struct {
	service InvitationService
}

func NewInvitationHandler(service InvitationService) *InvitationHandler {
	return &InvitationHandler{
		service: service,
	}
}

// ============================================================
// Create Invitation
// ============================================================

// CreateInvitation handles:
//
// POST /relationships/:relationshipID/invitations
//
// Creates an invitation for an email address.
func (h *InvitationHandler) CreateInvitation(c *gin.Context) {
	relationshipID, err := httpx.ParamInt64(
		c,
		"relationshipID",
	)
	if err != nil {
		response.BadRequest(c, "invalid relationship id")
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	var req ExternalInvitationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	result, err := h.service.CreateExternalInvitation(
		c.Request.Context(),
		relationshipID,
		userID,
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.Created(c, result)
}

func (h *InvitationHandler) GetInvitationByToken(c *gin.Context) {

	token := strings.TrimSpace(
		c.Param("token"),
	)

	if token == "" {
		response.BadRequest(
			c,
			"invalid invitation token",
		)
		return
	}

	result, err := h.service.GetInvitationByToken(
		c.Request.Context(),
		token,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, result)
}

// POST /invitations/:token/accept
//
// Accepts an invitation and adds the authenticated user
// as a member of the relationship.
func (h *InvitationHandler) AcceptInvitation(c *gin.Context) {

	token := strings.TrimSpace(
		c.Param("token"),
	)

	if token == "" {
		response.BadRequest(
			c,
			"invalid invitation token",
		)
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.AcceptInvitation(
		c.Request.Context(),
		token,
		userID,
	); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}

// POST /invitations/:token/reject
//
// Rejects an invitation.
func (h *InvitationHandler) RejectInvitation(c *gin.Context) {

	token := strings.TrimSpace(
		c.Param("token"),
	)

	if token == "" {
		response.BadRequest(
			c,
			"invalid invitation token",
		)
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.RejectInvitation(
		c.Request.Context(),
		token,
		userID,
	); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}

// DELETE /invitations/:invitationId
//
// Cancels an invitation created by the authenticated user.
func (h *InvitationHandler) CancelInvitation(c *gin.Context) {

	invitationID, err := httpx.ParamInt64(
		c,
		"invitationId",
	)
	if err != nil {
		response.BadRequest(
			c,
			"invalid invitation id",
		)
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.service.CancelInvitation(
		c.Request.Context(),
		invitationID,
		userID,
	); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}
