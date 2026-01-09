package web

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jian-hua-he/ddd_notes/internal/adapter/web/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoswagger "github.com/swaggo/echo-swagger"
)

const (
	UrlPathNoteList = "/notes"
)

type Router struct {
	httpPort int
	handler  *Handler
	echo     *echo.Echo
}

func NewRouter(httpPort int, handler *Handler, e *echo.Echo) *Router {
	return &Router{
		httpPort: httpPort,
		handler:  handler,
		echo:     e,
	}
}

// Setup
//
// @Title App API
// @Version 1.0
// @Description This is a sample server for a App application.
// @BasePath /app
func (r *Router) SetUp() error {
	r.echo.HideBanner = true

	r.echo.Use(
		middleware.Recover(),
		middleware.CORS(),
		middleware.RequestID(),
		middleware.ContextTimeout(60*time.Second),
	)

	group := r.echo.Group("app")

	v1 := group.Group("/v1")
	v1.GET("/swagger/*", echoswagger.WrapHandler)
	v1.GET(UrlPathNoteList, r.handler.ListNotes)

	return nil
}

func (r *Router) Start() error {
	return r.echo.Start(fmt.Sprintf(":%d", r.httpPort))
}

func (r *Router) Shutdown(ctx context.Context) error {
	return r.echo.Shutdown(ctx)
}
