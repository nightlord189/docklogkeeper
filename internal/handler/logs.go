package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GetLogs godoc
// @Tags log
// @Accept  json
// @Produce json
// @Param shortname path string true "container's short name"
// @Param direction query string true "future or past"
// @Param cursor query int false "cursor"
// @Param limit query int false "limit of result lines"
// @Success 200 {object} log.GetLogsResponse
// @Failure 400 {object} GenericResponse
// @Failure 401 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/container/{shortname}/log [Get]
// @BasePath /
func (h *Handler) GetLogs(c *gin.Context) {
	logger := zerolog.Ctx(c.Request.Context()).With().Str("handler_action", "GetLogs").Logger()
	ctx := logger.WithContext(c.Request.Context())

	var req entity.GetLogsRequest
	err := c.ShouldBindQuery(&req) // <-check there
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericErrorf("binding query error: %v", err))
		return
	}

	if err := req.IsValid(); err != nil {
		c.JSON(http.StatusBadRequest, GenericError(err.Error()))
		return
	}

	shortName := c.Param("shortname")
	if shortName == "" {
		c.JSON(http.StatusBadRequest, GenericError("empty shortname"))
		return
	}

	req.ShortName = shortName

	resp, err := h.LogAdapter.GetLogs(req)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("get logs error")
		c.JSON(http.StatusInternalServerError, GenericError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp)
}
