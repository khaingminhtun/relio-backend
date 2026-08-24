package relationship

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.RouterGroup,
	relationshipHandler *RelationshipHandler,
	relationshipMemberHandler *RelationshipMemberHandler,
) {
	relationships := router.Group("/relationships")

	// Relationship CRUD
	relationships.POST("/", relationshipHandler.CreateRelationship)
	relationships.GET("/", relationshipHandler.List)
	relationships.GET("/:relationshipId", relationshipHandler.GetByID)
	relationships.PATCH("/:relationshipId", relationshipHandler.Update)
	relationships.DELETE("/:relationshipId", relationshipHandler.Delete)

	// Relationship Members
	members := relationships.Group("/:relationshipId/members")

	members.GET("/", relationshipMemberHandler.List)
	members.GET("/me", relationshipMemberHandler.GetMe)
	members.GET("/:memberId", relationshipMemberHandler.GetByID)
	members.PATCH("/:memberId", relationshipMemberHandler.UpdateMember)
	members.DELETE("/:memberId", relationshipMemberHandler.RemoveMember)
}
