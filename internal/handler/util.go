package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getParamInt64(c *gin.Context, key string) (int64, error) {
	valueStr := c.Param(key)
	if valueStr == "" {
		return 0, fmt.Errorf("empty %s", key)
	}
	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return int64(valueInt), nil
}
