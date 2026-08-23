package relationship

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {

	relationship := router.Group("/relationship")

	relationship.POST("/", handler.CreateRelationship)

	relationship.GET("/", handler.List)

	relationship.GET("/:relationshipId", handler.GetByID)

	relationship.PATCH("/:relationshipId", handler.Update)

	relationship.DELETE("/:relationshipId", handler.Delete)

}
