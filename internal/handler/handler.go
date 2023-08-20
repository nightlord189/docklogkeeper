package handler

import (
	"fmt"
	"github.com/nightlord189/docklogkeeper/internal/config"
	docker2 "github.com/nightlord189/docklogkeeper/internal/docker"
	"io"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Config config.HTTPConfig
	Docker *docker2.Adapter
}

func New(cfg config.HTTPConfig, dock *docker2.Adapter) *Handler {
	return &Handler{Config: cfg, Docker: dock}
}

func (h *Handler) Run() error {
	gin.SetMode(h.Config.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	ginLoggerConfig := gin.LoggerConfig{}

	if h.Config.GinMode == "release" {
		ginLoggerConfig.Output = io.Discard
	}

	router.Use(gin.LoggerWithConfig(ginLoggerConfig))

	corsMiddleware := cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin", "Access-Control-Allow-Origin", "Content-Type", "Content-Length",
			"Accept-Encoding", "Authorization", "X-CSRF-Token", "X-Request-FolderName", "X-Forwarded-For", "Origin", "Referer",
		},
		ExposeHeaders:    []string{"Content-Length", "Content-Range", "X-Request-FolderName", "X-Forwarded-For", "Origin", "Referer"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	})

	router.Use(corsMiddleware)

	router.OPTIONS("/", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})
	return router.Run(fmt.Sprintf(":%d", h.Config.Port))
}
