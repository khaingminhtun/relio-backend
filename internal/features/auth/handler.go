package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/relio-backend/internal/shared/response"
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
// Register
// ============================================================

func (h *Handler) Register(c *gin.Context) {

	var req RegisterRequest

	// ----------------------------------------------------------
	// Bind request
	// ----------------------------------------------------------

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	// ----------------------------------------------------------
	// Register
	// ----------------------------------------------------------

	result, err := h.service.Register(
		c.Request.Context(),
		req,
	)

	if err != nil {
		// Let your centralized error middleware handle
		// application/domain errors.
		c.Error(err)
		return
	}

	// ----------------------------------------------------------
	// Response
	// ----------------------------------------------------------

	response.Accepted(c, result)
}

func (h *Handler) VerifyRegister(c *gin.Context) {

	var req VerifyRegisterRequest

	// ============================================================
	// Bind request
	// ============================================================

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}

	// ============================================================
	// Verify registration
	// ============================================================

	result, err := h.service.VerifyRegister(
		c.Request.Context(),
		req,
	)

	if err != nil {
		c.Error(err)
		return
	}

	// ============================================================
	// Response
	// ============================================================

	response.OK(c, result)
}

func (h *Handler) Authenticate(c *gin.Context) {
	var req LoginRequest

	// ========================================================
	// Bind request
	// ========================================================

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid login request")
		return
	}

	// ========================================================
	// Request metadata
	// ========================================================

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	// ========================================================
	// Authenticate
	// ========================================================

	result, err := h.service.Authenticate(
		c.Request.Context(),
		req,
		userAgent,
		ipAddress,
	)

	if err != nil {
		// Use your existing centralized error middleware.
		c.Error(err)
		return
	}

	// ========================================================
	// Response
	// ========================================================

	response.OK(c, result)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid request body")
		return
	}
	result, err := h.service.Refresh(
		c.Request.Context(),
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, result)
}
