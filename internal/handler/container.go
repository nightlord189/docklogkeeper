package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetContainers godoc
// @Tags container
// @Accept  json
// @Produce json
// @Success 200 {object} GetContainersResponse
// @Failure 403 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/container [Get]
// @BasePath /
func (h *Handler) GetContainers(c *gin.Context) {
	containers, err := h.Usecase.GetContainers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericError(err.Error()))
	}
	c.JSON(http.StatusOK, GetContainersResponse{Containers: containers})
}
