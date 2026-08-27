package relationship

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.RouterGroup,
	relationshipHandler *RelationshipHandler,
	relationshipMemberHandler *RelationshipMemberHandler,
	invitationHandler *InvitationHandler,
) {
	relationships := router.Group("/relationships")

	// Relationship CRUD
	relationships.POST("/", relationshipHandler.CreateRelationship)
	relationships.GET("/", relationshipHandler.List)
	relationships.GET("/:relationshipId", relationshipHandler.GetByID)
	relationships.PATCH("/:relationshipId", relationshipHandler.Update)
	relationships.DELETE("/:relationshipId", relationshipHandler.Delete)
	relationships.POST(
		"/:relationshipID/invitations",
		invitationHandler.CreateInvitation,
	)

	// Relationship Members
	members := relationships.Group("/:relationshipId/members")

	members.GET("/", relationshipMemberHandler.List)
	members.GET("/me", relationshipMemberHandler.GetMe)
	members.GET("/:memberId", relationshipMemberHandler.GetByID)
	members.PATCH("/:memberId", relationshipMemberHandler.UpdateMember)
	members.DELETE("/:memberId", relationshipMemberHandler.RemoveMember)

	invitations := router.Group("/invitations")

	invitations.POST(
		"/:token/accept",
		invitationHandler.AcceptInvitation,
	)

	invitations.POST(
		"/:token/reject",
		invitationHandler.RejectInvitation,
	)

	invitations.DELETE(
		"/:invitationId",
		invitationHandler.CancelInvitation,
	)

}

func RegisterPublicRoutes(
	router *gin.RouterGroup,
	invitationHandler *InvitationHandler,
) {
	invitations := router.Group("/invitations")

	invitations.GET(
		"/:token",
		invitationHandler.GetInvitationByToken,
	)
}
