package main

import (
	"context"
	"os/signal"
	"syscall"

	dbpostgres "github.com/jian-hua-he/reference-app/internal/adapter/database/postgres"
	"github.com/jian-hua-he/reference-app/internal/adapter/database/postgres/migration"
	grpchandler "github.com/jian-hua-he/reference-app/internal/adapter/grpc/handler"
	"github.com/jian-hua-he/reference-app/internal/adapter/grpc/server"
	webhandler "github.com/jian-hua-he/reference-app/internal/adapter/web/handler"
	"github.com/jian-hua-he/reference-app/internal/adapter/web/router"
	"github.com/jian-hua-he/reference-app/internal/application/note"
	"github.com/jian-hua-he/reference-app/internal/config"
	notepostgres "github.com/jian-hua-he/reference-app/internal/repository/note/postgres"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	// Database
	db, err := dbpostgres.NewDB(dbpostgres.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := migration.Up(db); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Msg("failed to run database migrations")
	}

	repo := notepostgres.NewRepo(db)
	app := note.NewNoteApp(repo)

	// HTTP server
	e := echo.New()
	wh := webhandler.NewHandler(app)
	r := router.NewRouter(cfg.HTTP.Port, wh, e)

	if err := r.SetUp(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to set up HTTP router")
		return
	}

	// gRPC server
	gh := grpchandler.NewHandler(app)
	gs := server.NewServer(cfg.GRPC.Port, gh)

	// Start HTTP server
	go func() {
		log.Ctx(ctx).Info().Msgf("starting HTTP server on port %d", cfg.HTTP.Port)
		if err := r.Start(); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("HTTP server stopped with error")
			stop()
		}
	}()

	// Start gRPC server
	go func() {
		log.Ctx(ctx).Info().Msgf("starting gRPC server on port %d", cfg.GRPC.Port)
		if err := gs.Start(); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("gRPC server stopped with error")
			stop()
		}
	}()

	<-ctx.Done()

	log.Ctx(ctx).Info().Msg("shutting down servers")

	gs.Shutdown()

	if err := r.Shutdown(ctx); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to shut down HTTP server")
	}
}
