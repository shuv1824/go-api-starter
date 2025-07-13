package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
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
	// ctx := context.Background()

	cfg, err := InitConfig("./config.yaml")
	if err != nil {
		log.Fatalf("failed to initialize config: %v\n", err)
	}

	level := slog.LevelInfo
	if cfg.Mode == modeTypeDebug {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	err = router.Run(fmt.Sprintf(":%d", 8000))
	if err != nil {
		log.Fatalf("error starting api: %v\n", err)
	}

	logger.Info("server started", "port", 8000)
}
