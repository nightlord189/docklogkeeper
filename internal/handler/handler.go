package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/nightlord189/docklogkeeper/docs" // for swagger docs
	"github.com/nightlord189/docklogkeeper/internal/config"
	"github.com/nightlord189/docklogkeeper/internal/log"
	"github.com/nightlord189/docklogkeeper/internal/repo"
	"github.com/nightlord189/docklogkeeper/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Config     config.Config
	Repo       *repo.Repo
	Usecase    *usecase.Usecase
	LogAdapter *log.Adapter
}

func New(cfg config.Config, repoInst *repo.Repo, ucInst *usecase.Usecase, lgAdapter *log.Adapter) *Handler {
	return &Handler{Config: cfg, Repo: repoInst, Usecase: ucInst, LogAdapter: lgAdapter}
}

//nolint:funlen
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
	store.Options(sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   7 * 86400,
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
	})

	router.Use(sessions.Sessions(sessionName, store))

	router.GET("/swagger/*any", func(context *gin.Context) {
		ginSwagger.WrapHandler(swaggerFiles.Handler)(context)
	})

	router.OPTIONS("/", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})

	router.POST("/api/auth", h.Auth)
	router.GET("/api/container", h.CookieAuthMdw, h.GetContainers)
	router.GET("/api/container/:shortname/log", h.CookieAuthMdw, h.GetLogs)
	router.GET("/api/container/:shortname/log/search", h.CookieAuthMdw, h.SearchLogs)

	router.GET("/api/trigger", h.CookieAuthMdw, h.GetTriggers)
	router.POST("/api/trigger", h.CookieAuthMdw, h.CreateTrigger)
	router.PUT("/api/trigger/:id", h.CookieAuthMdw, h.UpdateTrigger)
	router.DELETE("/api/trigger/:id", h.CookieAuthMdw, h.DeleteTrigger)

	htmlPages := []string{
		"static/web/auth.html",
		"static/web/logs.html",
		"static/web/triggers.html",
		"static/web/trigger_edit.html",
	}
	router.LoadHTMLFiles(htmlPages...)

	router.Static("/js", "static/web/js")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "auth.html", TemplateData{Analytics: h.Config.Analytics})
	})

	router.GET("/logs", func(c *gin.Context) {
		c.HTML(http.StatusOK, "logs.html", TemplateData{Analytics: h.Config.Analytics})
	})

	router.GET("/triggers", func(c *gin.Context) {
		c.HTML(http.StatusOK, "triggers.html", TemplateData{Analytics: h.Config.Analytics})
	})

	router.GET("/trigger/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "trigger_edit.html", TemplateData{Analytics: h.Config.Analytics})
	})

	router.GET("/trigger/:id/edit", func(c *gin.Context) {
		c.HTML(http.StatusOK, "trigger_edit.html", TemplateData{Analytics: h.Config.Analytics})
	})

	return router.Run(fmt.Sprintf(":%d", h.Config.HTTP.Port))
}
