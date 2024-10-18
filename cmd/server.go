package main

import (
	"fmt"
	"log/slog"
	"net/http"

	_ "net/http/pprof"

	graphqlplayground "github.com/99designs/gqlgen/graphql/playground"
	echoprometheus "github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	slogecho "github.com/samber/slog-echo"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type serverCmd struct {
	// cli options
	HttpAddr   string `help:"address of the http server which the server should listen on" default:":8080"`
	DbAddr     string `help:"address of the database server" default:"127.0.0.1:5432"`
	DbPassword string `help:"password for the database server" default:"postgres"`
	SslMode    string `help:"ssl mode for the database connection" default:"disable"`

	// Dependencies
	logger *slog.Logger
	db     *gorm.DB
}

func dsn(dbAddr string, dbPassword string, sslMode string) string {
	return fmt.Sprintf("postgres://postgres:%s@%s/postgres?sslmode=%s", dbPassword, dbAddr, sslMode)
}

func (s *serverCmd) Run(cmdCtx *cmdContext) error {
	s.logger = cmdCtx.Logger.With("component", "serverCmd")
	s.logger.Info("Starting the server")

	var err error

	// Connect to the database
	psqlDsn := dsn(s.DbAddr, s.DbPassword, s.SslMode)
	db, err := gorm.Open(postgres.Open(psqlDsn))
	if err != nil {
		return errors.Wrap(err, "failed to initialize gorm")
	}
	err = db.AutoMigrate()
	if err != nil {
		return errors.Wrap(err, "failed to run migrations")
	}

	// initialize echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// echo middlewares
	slogEchoConfig := slogecho.Config{
		WithSpanID:    true, // OTEL
		WithTraceID:   true, // OTEL
		WithUserAgent: true,
	}
	e.Use(slogecho.NewWithConfig(s.logger.With("subcomponent", "echo"), slogEchoConfig))
	e.Use(middleware.Recover())
	e.Use(echoprometheus.NewMiddleware("echo"))

	// admin routes
	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	e.GET("/healthz", s.Healthz)

	// graphql
	playgroundHandler := graphqlplayground.Handler("GraphQL playground", "/query")

	// plain http routes
	e.GET("/ping", s.Ping)

	// graphql routes
	e.GET("/playground", echo.WrapHandler(playgroundHandler))
	//e.POST("/query", s.Query) // TODO: implement the graphql endpoint

	return e.Start(s.HttpAddr)
}

func (s *serverCmd) Healthz(c echo.Context) error {
	return c.String(200, "OK")
}

func (s *serverCmd) Ping(c echo.Context) error {
	return c.String(200, "pong")
}
