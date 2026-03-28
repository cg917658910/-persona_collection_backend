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

	r.Static(cfg.StaticMountPrefix, cfg.StaticLocalDir)
	assetURLBuilder := service.NewAssetURLBuilder(cfg.PublicBaseURL, cfg.StaticMountPrefix)

	var (
		catalogRepo repo.CatalogRepo
		adminRepo   repo.AdminRepo
	)

	if cfg.UseMock {
		catalogRepo = repo.NewMockCatalogRepo()
		adminRepo = repo.NewMockAdminRepo()
	} else {
		pool, err := store.NewPostgresPool(cfg.DatabaseURL)
		if err != nil {
			log.Printf("postgres init failed, fallback to mock: %v", err)
			catalogRepo = repo.NewMockCatalogRepo()
			adminRepo = repo.NewMockAdminRepo()
		} else {
			catalogRepo = repo.NewPostgresCatalogRepo(pool)
			adminRepo = repo.NewPostgresAdminRepo(pool)
		}
	}

	catalogSvc := service.NewCatalogService(catalogRepo, assetURLBuilder)
	catalogHandler := handler.NewCatalogHandler(catalogSvc, cfg.PublicBaseURL, cfg.FrontendBaseURL)
	r.GET("/share/character/:slug", catalogHandler.CharacterSharePage)
	r.GET("/share/relation/:slug", catalogHandler.RelationSharePage)
	adminSvc := service.NewAdminService(adminRepo, assetURLBuilder)
	adminHandler := handler.NewAdminHandler(adminSvc)
	adminWorkCreatorSvc := service.NewAdminWorkCreatorService(adminRepo, assetURLBuilder)
	adminWorkCreatorHandler := handler.NewAdminWorkCreatorHandler(adminWorkCreatorSvc)
	adminDictSvc := service.NewAdminDictService(adminRepo)
	adminDictHandler := handler.NewAdminDictHandler(adminDictSvc)
	adminPageSvc := service.NewAdminPageService(adminRepo)
	adminPageHandler := handler.NewAdminPageHandler(adminPageSvc)
	adminImportSvc := service.NewAdminImportService(adminRepo)
	adminImportHandler := handler.NewAdminImportHandler(adminImportSvc)
	adminAuthSvc := service.NewAdminAuthService(cfg)
	adminAuthHandler := handler.NewAdminAuthHandler(adminAuthSvc)
	adminResourceSvc := service.NewAdminResourceService(cfg, assetURLBuilder)
	adminResourceHandler := handler.NewAdminResourceHandler(adminResourceSvc)

	api := r.Group(cfg.APIPrefix)
	{
		// public catalog
		api.GET("/home", catalogHandler.Home)
		api.GET("/discover/random", catalogHandler.RandomCharacter)
		api.GET("/characters", catalogHandler.ListCharacters)
		api.GET("/characters/:slug", catalogHandler.GetCharacterDetail)
		api.GET("/relations", catalogHandler.ListRelationships)
		api.GET("/relations/:slug", catalogHandler.GetRelationshipDetail)
		api.GET("/relationships", catalogHandler.ListRelationships)
		api.GET("/relationships/:slug", catalogHandler.GetRelationshipDetail)
		api.GET("/works", catalogHandler.ListWorks)
		api.GET("/works/:slug", catalogHandler.GetWorkDetail)
		api.GET("/creators", catalogHandler.ListCreators)
		api.GET("/creators/:slug", catalogHandler.GetCreatorDetail)
		api.GET("/themes", catalogHandler.ListThemes)
		api.GET("/themes/:slug", catalogHandler.GetThemeDetail)
		api.GET("/songs", catalogHandler.ListSongs)
		api.GET("/search", catalogHandler.Search)

		// original admin auth endpoints
		adminPublic := api.Group("/admin")
		{
			auth := adminPublic.Group("/auth")
			{
				auth.POST("/login", adminAuthHandler.Login)
			}
		}

		adminProtected := api.Group("/admin")
		adminProtected.Use(middleware.AdminJWT(cfg.AdminJWTSecret))
		{
			adminProtected.GET("/auth/me", adminAuthHandler.Me)

			adminProtected.GET("/characters", adminHandler.ListCharacters)
			adminProtected.GET("/characters/page", adminPageHandler.Characters)
			adminProtected.GET("/characters/:ref", adminHandler.GetCharacter)
			adminProtected.POST("/characters", adminHandler.CreateCharacter)
			adminProtected.PATCH("/characters/:ref", adminHandler.UpdateCharacter)
			adminProtected.DELETE("/characters/:ref", adminHandler.DeleteCharacter)

			adminProtected.GET("/songs", adminHandler.ListSongs)
			adminProtected.GET("/songs/page", adminPageHandler.Songs)
			adminProtected.GET("/songs/:ref", adminHandler.GetSong)
			adminProtected.POST("/songs", adminHandler.CreateSong)
			adminProtected.PATCH("/songs/:ref", adminHandler.UpdateSong)
			adminProtected.DELETE("/songs/:ref", adminHandler.DeleteSong)

			adminProtected.GET("/relations", adminHandler.ListRelations)
			adminProtected.GET("/relations/page", adminPageHandler.Relations)
			adminProtected.GET("/relations/:ref", adminHandler.GetRelation)
			adminProtected.POST("/relations", adminHandler.CreateRelation)
			adminProtected.PATCH("/relations/:ref", adminHandler.UpdateRelation)
			adminProtected.DELETE("/relations/:ref", adminHandler.DeleteRelation)

			adminProtected.GET("/themes", adminHandler.ListThemes)
			adminProtected.GET("/themes/page", adminPageHandler.Themes)
			adminProtected.GET("/themes/:ref", adminHandler.GetTheme)
			adminProtected.POST("/themes", adminHandler.CreateTheme)
			adminProtected.PATCH("/themes/:ref", adminHandler.UpdateTheme)
			adminProtected.DELETE("/themes/:ref", adminHandler.DeleteTheme)

			adminProtected.GET("/resources/:type", adminResourceHandler.ListResources)
			adminProtected.GET("/resources/:type/page", adminResourceHandler.ListResources)
			adminProtected.POST("/resources/:type", adminResourceHandler.CreateResource)
			adminProtected.POST("/resources/:type/upload", adminResourceHandler.UploadResource)
			adminProtected.DELETE("/resources/:type/:ref", adminResourceHandler.DeleteResource)

			adminProtected.GET("/works", adminWorkCreatorHandler.ListWorks)
			adminProtected.GET("/works/page", adminPageHandler.Works)
			adminProtected.GET("/works/:ref", adminWorkCreatorHandler.GetWork)
			adminProtected.POST("/works", adminWorkCreatorHandler.CreateWork)
			adminProtected.PATCH("/works/:ref", adminWorkCreatorHandler.UpdateWork)
			adminProtected.DELETE("/works/:ref", adminWorkCreatorHandler.DeleteWork)

			adminProtected.GET("/creators", adminWorkCreatorHandler.ListCreators)
			adminProtected.GET("/creators/page", adminPageHandler.Creators)
			adminProtected.GET("/creators/:ref", adminWorkCreatorHandler.GetCreator)
			adminProtected.POST("/creators", adminWorkCreatorHandler.CreateCreator)
			adminProtected.PATCH("/creators/:ref", adminWorkCreatorHandler.UpdateCreator)
			adminProtected.DELETE("/creators/:ref", adminWorkCreatorHandler.DeleteCreator)

			adminProtected.GET("/dicts/:dictKey", adminDictHandler.List)
			adminProtected.GET("/dicts/:dictKey/page", adminPageHandler.Dicts)
			adminProtected.GET("/dicts/:dictKey/:ref", adminDictHandler.Get)
			adminProtected.POST("/dicts/:dictKey", adminDictHandler.Create)
			adminProtected.PATCH("/dicts/:dictKey/:ref", adminDictHandler.Update)
			adminProtected.DELETE("/dicts/:dictKey/:ref", adminDictHandler.Delete)

			adminProtected.POST("/imports/validate", adminImportHandler.Validate)
			adminProtected.POST("/imports/run", adminImportHandler.Run)
			adminProtected.POST("/relation-imports/validate", adminImportHandler.ValidateRelations)
			adminProtected.POST("/relation-imports/run", adminImportHandler.RunRelations)
		}

		// Vben compatible auth endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/login", adminAuthHandler.VbenLogin)
			auth.GET("/codes", adminAuthHandler.VbenPermissionCodes)
		}
		user := api.Group("/user")
		{
			user.GET("/info", adminAuthHandler.VbenUserInfo)
		}
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
