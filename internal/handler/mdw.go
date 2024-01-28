package handler

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *Handler) CookieAuthMdw(c *gin.Context) {
	ctx := c.Request.Context()

	session := sessions.Default(c)
	user := session.Get(userKey)

	if user != defaultUser {
		log.Ctx(ctx).Error().Msgf("user from cookie is %s, invalid cookie, aborting", user)

		session.Delete(userKey)

		if err := session.Save(); err != nil {
			log.Ctx(ctx).Error().Msgf("CookieAuthMdw: save session error: %v", err)
		}

		c.JSON(http.StatusUnauthorized, GenericError("invalid cookie"))

		c.Abort()
		return
	}

	c.Next()
}
