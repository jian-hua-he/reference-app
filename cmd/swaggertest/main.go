package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/jian-hua-he/ddd_notes/internal/adapter/web"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := NewFakeApp()
	port := 8082
	e := echo.New()
	h := web.NewHandler(app)
	r := web.NewRouter(port, h, e)

	if err := r.SetUp(); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to set up router")
		return
	}

	log.Ctx(ctx).Info().Msgf("starting server on port %d", port)

	go func() {
		if err := r.Start(); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("server stopped with error")
			stop()
		}
	}()

	<-ctx.Done()

	log.Ctx(ctx).Info().Msg("shutting down server")

	if err := r.Shutdown(ctx); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to shut down server")
	}
}
