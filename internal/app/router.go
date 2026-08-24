package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/relio-backend/internal/features/relationship"
	"github.com/khaingminhtun/relio-backend/internal/shared/middleware"

	"github.com/khaingminhtun/relio-backend/internal/features/auth"
	"github.com/khaingminhtun/relio-backend/internal/features/user"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/httperror"
)

func NewRouter(deps *Dependencies) *gin.Engine {
	router := gin.New()

	// Enable 405 Method Not Allowed handling.
	router.HandleMethodNotAllowed = true

	// Global middleware.
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.ErrorHandler())

	// 405 Method Not Allowed.
	router.NoMethod(func(c *gin.Context) {
		_ = c.Error(httperror.New(
			http.StatusMethodNotAllowed,
			"METHOD_NOT_ALLOWED",
			"the HTTP method used is not supported for this endpoint",
			nil,
		))
	})

	// 404 Route Not Found.
	router.NoRoute(func(c *gin.Context) {
		_ = c.Error(httperror.New(
			http.StatusNotFound,
			"ROUTE_NOT_FOUND",
			"the requested API endpoint does not exist",
			nil,
		))
	})

	api := router.Group("/api/v1")

	user.RegisterRoutes(
		api,
		deps.UserHandler,
	)

	auth.RegisterRoutes(
		api,
		deps.AuthHandler,
	)

	// ============================================================
	// Private routes
	// ============================================================

	authMiddleware := middleware.NewAuthMiddleware(
		deps.JWTManager,
	)

	private := api.Group("")
	private.Use(authMiddleware.RequireAuth())

	relationship.RegisterRoutes(
		private,
		deps.RelationshipHandler,
		deps.RelatonshipMemberHandler,
	)
	return router
}
