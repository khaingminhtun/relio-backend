package user

import (
	"github.com/gin-gonic/gin"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/httpx"
	"github.com/khaingminhtun/relio-backend/internal/shared/response"
)

type UserProfileHandler struct {
	service UserProfileService
}

func NewUserProfileHandler(
	service UserProfileService,
) *UserProfileHandler {
	return &UserProfileHandler{
		service: service,
	}
}

// GET /profile
func (h *UserProfileHandler) Get(c *gin.Context) {
	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeUnauthorized,
			"authenticated user not found",
			err,
		))
		return
	}

	profile, err := h.service.Get(
		c.Request.Context(),
		userID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, profile)
}

// PATCH /profile
func (h *UserProfileHandler) Update(c *gin.Context) {
	var req UpdateUserProfileRequest

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
		c.Error(apperror.New(
			apperror.CodeUnauthorized,
			"authenticated user not found",
			err,
		))
		return
	}

	profile, err := h.service.Update(
		c.Request.Context(),
		userID,
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, profile)
}
