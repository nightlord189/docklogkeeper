package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/nightlord189/docklogkeeper/internal/config"
	docker2 "github.com/nightlord189/docklogkeeper/internal/docker"
	"io"
	"time"

	_ "github.com/nightlord189/docklogkeeper/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Config config.Config
	Docker *docker2.Adapter
}

func New(cfg config.Config, dock *docker2.Adapter) *Handler {
	return &Handler{Config: cfg, Docker: dock}
}

func (h *Handler) Run() error {
	gin.SetMode(h.Config.HTTP.GinMode)
	router := gin.New()
	router.Use(gin.Recovery())
	ginLoggerConfig := gin.LoggerConfig{}

	if h.Config.HTTP.GinMode == "release" {
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

	store := cookie.NewStore([]byte(h.Config.Auth.Secret))
	router.Use(sessions.Sessions("defaultsession", store))

	router.GET("/swagger/*any", func(context *gin.Context) {
		ginSwagger.WrapHandler(swaggerFiles.Handler)(context)
	})

	router.OPTIONS("/", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})

	router.POST("/api/auth", h.Auth)
	
	return router.Run(fmt.Sprintf(":%d", h.Config.HTTP.Port))
}
