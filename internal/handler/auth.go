package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	userKey     = "user_id"
	defaultUser = "admin"
)

// Auth godoc
// @Tags auth
// @Accept  json
// @Produce json
// @Param data body AuthRequest true "Input model"
// @Success 200 {object} GenericResponse
// @Failure 400 {object} GenericResponse
// @Failure 401 {object} GenericResponse
// @Failure 500 {object} GenericResponse
// @Router /api/auth [Post]
// @BasePath /
func (h *Handler) Auth(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenericError("parse json error: "+err.Error()))
		return
	}

	if req.Username != defaultUser || req.Password != h.Config.Auth.Password {
		c.JSON(http.StatusUnauthorized, GenericError("bad credentials"))
		return
	}

	session := sessions.Default(c)
	session.Set(userKey, defaultUser)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, GenericError("save session error: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, GenericError("authenticated"))
}
