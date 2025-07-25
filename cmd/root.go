package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shuv1824/go-api-starter/internal/common/auth"
	"github.com/shuv1824/go-api-starter/internal/common/middleware"
	"github.com/shuv1824/go-api-starter/internal/config"
	userHandlers "github.com/shuv1824/go-api-starter/internal/domains/user/handlers"
	userDomain "github.com/shuv1824/go-api-starter/internal/domains/user/infra"
	"github.com/shuv1824/go-api-starter/internal/migration"
	"github.com/shuv1824/go-api-starter/pkg/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-api-starter",
	Short: "A Gin-based REST API with JWT authentication",
	Long:  `A production-ready REST API template built with Gin, JWT authentication, GORM, and PostgreSQL.`,
	Run:   rootRun,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func rootRun(cmd *cobra.Command, args []string) {
	cfg, err := config.InitConfig("./config.yaml")
	if err != nil {
		log.Fatalf("failed to initialize config: %v\n", err)
	}

	gin.SetMode(string(cfg.Mode))

	// level := slog.LevelInfo
	// if cfg.Mode == config.ModeTypeDebug {
	// 	level = slog.LevelDebug
	// }
	//
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: level,
	// }))

	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to database: %v\n", err)
	}

	err = migration.MigrateUp(db, cfg.Database.Type)
	if err != nil {
		log.Fatalf("database migration error: %v\n", err)
	}

	// Initialize services
	jwtService := auth.NewService(cfg.Secret, time.Hour, time.Hour*24)
	userRepo := userDomain.NewRepository(db)
	userService := userDomain.NewService(userRepo, jwtService)
	userHandlers := userHandlers.NewHandler(userService)

	router := gin.Default()

	// Add middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(gin.Recovery())

	// Healthcheck
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Public routes
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/register", userHandlers.Register)
		auth.POST("/login", userHandlers.Login)
	}

	// Authenticated routes
	authenticated := router.Group("/api/v1")
	{
		authenticated.Use(middleware.AuthMiddleware(jwtService))
		authenticated.GET("/profile", userHandlers.GetProfile)
	}

	err = router.Run(fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("error starting api: %v\n", err)
	}
}
