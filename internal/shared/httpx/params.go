package httpx

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
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
