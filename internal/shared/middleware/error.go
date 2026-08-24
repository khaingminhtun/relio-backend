package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/httperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		httpErr := httperror.FromError(err)

		// Log the real internal error.
		log.Error().
			Err(err).
			Int("status", httpErr.Status).
			Str("code", httpErr.Code).
			Msg("request error")

		c.JSON(httpErr.Status, response.Response{
			Success: false,
			Error: &response.ErrorInfo{
				Code:    httpErr.Code,
				Message: httpErr.Message,
			},
		})
	}
}
