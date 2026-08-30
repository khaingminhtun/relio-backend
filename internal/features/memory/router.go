package memory

import "github.com/gin-gonic/gin"

// RegisterRoutes registers the five Memory CRUD routes under the
// /relationships/:relationshipId/memories path.
// The router group passed in is expected to already carry the auth middleware.
func RegisterRoutes(
	router *gin.RouterGroup,
	memoryHandler *MemoryHandler,
) {
	memories := router.Group("/relationships/:relationshipId/memories")

	memories.POST("/", memoryHandler.Create)
	memories.GET("/", memoryHandler.List)
	memories.GET("/:memoryId", memoryHandler.Get)
	memories.PATCH("/:memoryId", memoryHandler.Update)
	memories.DELETE("/:memoryId", memoryHandler.Delete)
}
