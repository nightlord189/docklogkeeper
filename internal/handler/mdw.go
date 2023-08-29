package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CookieAuthMdw(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userKey)
	if user != defaultUser {
		c.JSON(http.StatusUnauthorized, GenericError("invalid cookie"))
		c.Abort()
		return
	}
	c.Next()
}
