package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	grpchandler "github.com/jian-hua-he/reference-app/internal/adapter/grpc/handler"
	"github.com/jian-hua-he/reference-app/internal/adapter/grpc/server"
	webhandler "github.com/jian-hua-he/reference-app/internal/adapter/web/handler"
	"github.com/jian-hua-he/reference-app/internal/adapter/web/router"
	"github.com/jian-hua-he/reference-app/internal/application/note"
	"github.com/jian-hua-he/reference-app/internal/repository/note/memory"
	"github.com/jian-hua-he/reference-app/pkg/uuid"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	repo := memory.NewRepo(uuid.NewUUID, time.Now)
	app := note.NewNoteApp(repo)

	// HTTP server
	httpPort := 8082
	e := echo.New()
	wh := webhandler.NewHandler(app)
	r := router.NewRouter(httpPort, wh, e)

	if err := r.SetUp(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to set up HTTP router")
		return
	}

	// gRPC server
	grpcPort := 50051
	gh := grpchandler.NewHandler(app)
	gs := server.NewServer(grpcPort, gh)

	// Start HTTP server
	go func() {
		log.Ctx(ctx).Info().Msgf("starting HTTP server on port %d", httpPort)
		if err := r.Start(); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("HTTP server stopped with error")
			stop()
		}
	}()

	// Start gRPC server
	go func() {
		log.Ctx(ctx).Info().Msgf("starting gRPC server on port %d", grpcPort)
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
