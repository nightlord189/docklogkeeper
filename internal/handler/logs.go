package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/nightlord189/docklogkeeper/internal/entity"
	"github.com/nightlord189/docklogkeeper/internal/log"
	"github.com/rs/zerolog"
	"net/http"
)

// GetLogs godoc
// @Tags log
// @Accept  json
// @Produce json
// @Param shortname path string true "container's short name"
// @Param chunk_number query int false "number of chunk"
// @Param offset query int false "offset in chunk"
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

	var req GetLogsRequest
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

	var resp log.GetLogsResponse

	switch req.Direction {
	case entity.DirFuture:
		resp, err = h.LogAdapter.GetLogsUpdate(ctx, log.GetLinesRequest{
			ShortName:   shortName,
			ChunkNumber: req.ChunkNumber,
			Offset:      req.Offset,
			Limit:       req.Limit,
		})
	case entity.DirPast:
		resp, err = h.LogAdapter.GetLogsNext(ctx, log.GetLinesRequest{
			ShortName:   shortName,
			ChunkNumber: req.ChunkNumber,
			Offset:      req.Offset,
			Limit:       req.Limit,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp)
}
