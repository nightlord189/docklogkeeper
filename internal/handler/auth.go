package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	sessionName = "default"
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
	ctx := c.Request.Context()

	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Ctx(ctx).Error().Msgf("auth: parse json error: %v", err)
		c.JSON(http.StatusBadRequest, GenericError("parse json error: "+err.Error()))
		return
	}

	if req.Username != defaultUser || req.Password != h.Config.Auth.Password {
		log.Ctx(ctx).Error().Msgf("auth: bad credentials: %s", req.Username)
		c.JSON(http.StatusUnauthorized, GenericError("bad credentials"))
		return
	}

	session := sessions.Default(c)
	session.Set(userKey, req.Username)

	if err := session.Save(); err != nil {
		log.Ctx(ctx).Error().Msgf("auth: save session error: %v", err)
		c.JSON(http.StatusInternalServerError, GenericError("save session error: "+err.Error()))
		return
	}

	log.Ctx(ctx).Info().Msgf("auth: authenticated user %s", req.Username)
	c.JSON(http.StatusOK, GenericError("authenticated"))
}
