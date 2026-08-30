package note

import (
	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/httpx"
	"github.com/khaingminhtun/relio-backend/internal/shared/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// POST /notes
func (h *Handler) Create(c *gin.Context) {
	var req CreateNoteRequest

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

// GET /notes
func (h *Handler) List(c *gin.Context) {
	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	// --------------------------------------------------------
	// Parse optional query filters
	// --------------------------------------------------------

	filter := ListNotesFilter{}

	if mood := c.Query("mood"); mood != "" {
		filter.Mood = &mood
	}

	if pinned := c.Query("pinned"); pinned != "" {
		isPinned := pinned == "true"
		filter.IsPinned = &isPinned
	}

	if archived := c.Query("archived"); archived != "" {
		isArchived := archived == "true"
		filter.IsArchived = &isArchived
	}

	if search := c.Query("search"); search != "" {
		filter.Search = &search
	}

	notes, err := h.service.List(
		c.Request.Context(),
		userID,
		filter,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, notes)
}

// GET /notes/:noteId
func (h *Handler) Get(c *gin.Context) {
	noteID, err := httpx.ParamInt64(c, "noteId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid note id",
			err,
		))
		return
	}

	userID, err := httpx.UserID(c)
	if err != nil {
		c.Error(err)
		return
	}

	note, err := h.service.GetByID(
		c.Request.Context(),
		userID,
		noteID,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, note)
}

// PATCH /notes/:noteId
func (h *Handler) Update(c *gin.Context) {
	noteID, err := httpx.ParamInt64(c, "noteId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid note id",
			err,
		))
		return
	}

	var req UpdateNoteRequest

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
		noteID,
		req,
	)
	if err != nil {
		c.Error(err)
		return
	}

	response.OK(c, result)
}

// DELETE /notes/:noteId
func (h *Handler) Delete(c *gin.Context) {
	noteID, err := httpx.ParamInt64(c, "noteId")
	if err != nil {
		c.Error(apperror.New(
			apperror.CodeInvalidRequest,
			"invalid note id",
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
		noteID,
	); err != nil {
		c.Error(err)
		return
	}

	response.NoContent(c)
}
