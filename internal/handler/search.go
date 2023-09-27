package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nightlord189/docklogkeeper/internal/log"
	"github.com/rs/zerolog"
)

// SearchLogs godoc
// @Tags log
// @Accept  json
// @Produce json
// @Param shortname path string true "container's short name"
// @Param contains query string false "contains substring"
// @Success 200 {object} SearchLogsResponse
// @Failure 400 {object} GenericResponse
// @Failure 401 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/container/{shortname}/log/search [Get]
// @BasePath /
func (h *Handler) SearchLogs(c *gin.Context) {
	logger := zerolog.Ctx(c.Request.Context()).With().Str("handler_action", "SearchLogs").Logger()
	ctx := logger.WithContext(c.Request.Context())

	var req SearchLogsRequest
	err := c.ShouldBindQuery(&req) // <-check there
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericErrorf("binding query error: %v", err))
		return
	}

	shortName := c.Param("shortname")
	if shortName == "" {
		c.JSON(http.StatusBadRequest, GenericError("empty shortname"))
		return
	}

	logs := h.LogAdapter.SearchLines(ctx, shortName, log.SearchRequest{Contains: req.Contains})

	c.JSON(http.StatusOK, SearchLogsResponse{Records: logs})
}
