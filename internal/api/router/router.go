package router

import (
	"log"
	"net/http"

	"pm-backend/internal/config"
	"pm-backend/internal/handler"
	"pm-backend/internal/middleware"
	"pm-backend/internal/repo"
	"pm-backend/internal/service"
	"pm-backend/internal/store"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer(cfg config.Config) *gin.Engine {
	gin.SetMode(modeFor(cfg.AppEnv))
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(cfg.AllowOrigins))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"app":    cfg.AppName,
			"mock":   cfg.UseMock,
		})
	})

	// static file serving for local assets
	r.Static(cfg.StaticMountPrefix, cfg.StaticLocalDir)

	assetURLBuilder := service.NewAssetURLBuilder(cfg.PublicBaseURL, cfg.StaticMountPrefix)

	var catalogRepo repo.CatalogRepo

	if cfg.UseMock {
		catalogRepo = repo.NewMockCatalogRepo()
	} else {
		pool, err := store.NewPostgresPool(cfg.DatabaseURL)
		if err != nil {
			log.Printf("postgres init failed, fallback to mock: %v", err)
			catalogRepo = repo.NewMockCatalogRepo()
		} else {
			catalogRepo = repo.NewPostgresCatalogRepo(pool)
		}
	}

	catalogSvc := service.NewCatalogService(catalogRepo, assetURLBuilder)
	catalogHandler := handler.NewCatalogHandler(catalogSvc)

	api := r.Group(cfg.APIPrefix)
	{
		api.GET("/home", catalogHandler.Home)
		api.GET("/discover/random", catalogHandler.RandomCharacter)
		api.GET("/characters", catalogHandler.ListCharacters)
		api.GET("/characters/:slug", catalogHandler.GetCharacterDetail)
		api.GET("/works", catalogHandler.ListWorks)
		api.GET("/works/:slug", catalogHandler.GetWorkDetail)
		api.GET("/creators", catalogHandler.ListCreators)
		api.GET("/creators/:slug", catalogHandler.GetCreatorDetail)
		api.GET("/themes", catalogHandler.ListThemes)
		api.GET("/themes/:slug", catalogHandler.GetThemeDetail)
		api.GET("/songs", catalogHandler.ListSongs)
	}

	return r
}

func modeFor(env string) string {
	switch env {
	case "prod", "production":
		return gin.ReleaseMode
	case "test":
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
