package relationship

import (
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
