package note

import "github.com/gin-gonic/gin"

// RegisterRoutes registers the five Note CRUD routes under /notes.
// The router group passed in is expected to already carry the auth middleware.
func RegisterRoutes(
	router *gin.RouterGroup,
	noteHandler *Handler,
) {
	notes := router.Group("/notes")

	notes.POST("", noteHandler.Create)
	notes.GET("", noteHandler.List)
	notes.GET("/:noteId", noteHandler.Get)
	notes.PATCH("/:noteId", noteHandler.Update)
	notes.DELETE("/:noteId", noteHandler.Delete)
}
