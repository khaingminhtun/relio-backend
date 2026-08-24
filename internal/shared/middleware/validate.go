package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/security"
)

const (
	UserIDKey   = "user_id"
	UserRoleKey = "user_role"
	ClaimsKey   = "claims"
)

type AuthMiddleware struct {
	jwtManager *security.JWTManager
}

func NewAuthMiddleware(
	jwtManager *security.JWTManager,
) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// ============================================================
		// Get Authorization header
		// ============================================================

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Error(apperror.New(
				apperror.CodeUnauthorized,
				"authorization header is required",
				nil,
			))
			c.Abort()
			return
		}

		// ============================================================
		// Extract Bearer token
		// ============================================================

		parts := strings.Fields(authHeader)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {
			c.Error(apperror.New(
				apperror.CodeUnauthorized,
				"invalid authorization header",
				nil,
			))
			c.Abort()
			return
		}

		tokenString := parts[1]

		if tokenString == "" {
			c.Error(apperror.New(
				apperror.CodeUnauthorized,
				"access token is required",
				nil,
			))
			c.Abort()
			return
		}

		// ============================================================
		// Validate access token
		// ============================================================

		claims, err := m.jwtManager.ValidateAccessToken(
			tokenString,
		)
		if err != nil {
			c.Error(apperror.New(
				apperror.CodeInvalidAccessToken,
				"invalid or expired access token",
				err,
			))
			c.Abort()
			return
		}

		// ============================================================
		// Validate user ID
		// ============================================================

		if claims.UserID <= 0 {
			c.Error(apperror.New(
				apperror.CodeInvalidAccessToken,
				"access token contains invalid user ID",
				nil,
			))
			c.Abort()
			return
		}

		// ============================================================
		// Store authenticated user information
		// ============================================================

		c.Set(UserIDKey, claims.UserID)
		c.Set(UserRoleKey, claims.Role)
		c.Set(ClaimsKey, claims)

		c.Next()
	}
}
