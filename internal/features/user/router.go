package user

import "github.com/gin-gonic/gin"

// RegisterRoutes registers all user module routes.
func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
	profileHandler *UserProfileHandler,
) {
	users := router.Group("/users")
	{
		users.POST("", handler.CreateUser)
		users.GET("/search", handler.SearchUsers)
		users.GET("", handler.ListUsers)
		users.GET("/:id", handler.GetUser)
		users.PATCH("/:id", handler.UpdateUser)
		users.DELETE("/:id", handler.DeleteUser)
	}

	profile := router.Group("/profile")
	profile.GET("/", profileHandler.Get)
	profile.PATCH("/", profileHandler.Update)
}
