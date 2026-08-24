package httpx

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/khaingminhtun/relio-backend/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/relio-backend/internal/shared/middleware"
)

var ErrInvalidParameter = errors.New("invalid parameter")

func ParamInt64(
	c *gin.Context,
	name string,
) (int64, error) {
	value := c.Param(name)

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrInvalidParameter
	}

	return id, nil
}

func QueryInt(
	c *gin.Context,
	name string,
	defaultValue int,
) (int, error) {
	value := c.Query(name)

	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}

func UserID(c *gin.Context) (int64, error) {
	value, exists := c.Get(middleware.UserIDKey)
	if !exists {
		return 0, apperror.New(
			apperror.CodeUnauthorized,
			"authenticated user not found",
			nil,
		)
	}

	userID, ok := value.(int64)
	if !ok || userID <= 0 {
		return 0, apperror.New(
			apperror.CodeUnauthorized,
			"invalid authenticated user",
			nil,
		)
	}

	return userID, nil
}
