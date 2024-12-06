package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"personal-blog/internal/models"
	"runtime/debug"
	"sync"

	"personal-blog/internal/env"
	"personal-blog/internal/version"

	"github.com/lmittmann/tint"
)

func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL   string
	httpPort  int
	basicAuth struct {
		username       string
		hashedPassword string
	}
	cookie struct {
		secretKey string
	}
}

type application struct {
	config       config
	logger       *slog.Logger
	wg           sync.WaitGroup
	articleModel models.ArticleModel
}

func run(logger *slog.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:4444")
	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.basicAuth.username = env.GetString("BASIC_AUTH_USERNAME", "admin")
	cfg.basicAuth.hashedPassword = env.GetString("BASIC_AUTH_HASHED_PASSWORD", "$2a$12$A5VGgLlxr45NJ78Mqpz0vO2qGRSqtdKa78Jy9gpyExkI9FP07rLJW")
	cfg.cookie.secretKey = env.GetString("COOKIE_SECRET_KEY", "j5vxnte2jrfmxgxx2irl7awqnsxzkhgn")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	app := &application{
		config: cfg,
		logger: logger,
	}

	return app.serveHTTP()
}
