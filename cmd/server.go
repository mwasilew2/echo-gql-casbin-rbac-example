package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "net/http/pprof"

	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	graphqlplayground "github.com/99designs/gqlgen/graphql/playground"
	"github.com/casbin/casbin/v2"
	"github.com/gorilla/sessions"
	echoprometheus "github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	slogecho "github.com/samber/slog-echo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/graph/model"
	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/middleware"
	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/service/authorization"
	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/util"

	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/graph"
)

type serverCmd struct {
	// cli options
	HttpAddr                 string `help:"address of the http server which the server should listen on" default:":8080"`
	DbAddr                   string `help:"address of the database server" default:"127.0.0.1:5432"`
	DbPassword               string `help:"password for the database server" default:"postgres"`
	SslMode                  string `help:"ssl mode for the database connection" default:"disable"`
	CookieStoreSigningKey    string `help:"secret key to use for signing cookies" default:"changemechangemechangemechangeme"`
	CookieStoreEncryptionKey string `help:"secret key to use for encrypting cookies" default:"changemechangemechangemechangeme"`

	// Dependencies
	logger               *slog.Logger
	db                   *gorm.DB
	store                *sessions.CookieStore
	ulidManager          *util.UlidManager
	authorizationService authorization.Authorization
}

func dsn(dbAddr string, dbPassword string, sslMode string) string {
	return fmt.Sprintf("postgres://postgres:%s@%s/postgres?sslmode=%s", dbPassword, dbAddr, sslMode)
}

func (s *serverCmd) Run(cmdCtx *cmdContext) error {
	s.logger = cmdCtx.Logger.With("component", "serverCmd")
	s.logger.Info(fmt.Sprintf("starting server on %s", s.HttpAddr))

	var err error

	// Connect to the database
	psqlDsn := dsn(s.DbAddr, s.DbPassword, s.SslMode)
	s.db, err = gorm.Open(postgres.Open(psqlDsn))
	if err != nil {
		return errors.Wrap(err, "failed to initialize gorm")
	}

	// Cookie store
	s.store = sessions.NewCookieStore([]byte(s.CookieStoreSigningKey), []byte(s.CookieStoreEncryptionKey))

	// ULID manager
	s.ulidManager = util.NewUlidManager()

	// Authorization service
	casbinEnforcer, err := casbin.NewEnforcer("rbac_with_domains_model.conf", "rbac_with_domains_policy.csv")
	// TODO: use gorm for storing policies: https://github.com/casbin/gorm-adapter
	// TODO: expose casbin policy creation through an API: https://casbin.org/docs/rbac-api/#addrolesforuser
	// TODO: use a different casbin models (the current one is extremely simple): https://github.com/casbin/casbin/tree/master/examples
	// TODO: use group membership from SSO: https://github.com/casbin/casbin/issues/929
	// TODO: test hierarchy of roles within the same account, but different groups: https://github.com/casbin/casbin/issues/493
	// TODO: leverage Go type system for referencing resources
	if err != nil {
		return errors.Wrap(err, "failed to initialize casbin enforcer")
	}
	authorizationService := authorization.NewCasbinAuthorizationService(casbinEnforcer)
	s.authorizationService = authorizationService

	// graphql
	graphResolver := graph.NewResolver(s.db, s.logger, s.ulidManager, s.authorizationService)
	graphqlHandler := graphqlhandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graphResolver}))
	playgroundHandler := graphqlplayground.Handler("GraphQL playground", "/query")

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
	e.Use(echomiddleware.Recover())
	e.Use(echoprometheus.NewMiddleware("echo"))
	e.Use(echomiddleware.Secure()) // protects from xss, clickjacking, etc
	e.Use(session.Middleware(s.store))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get(util.CookieKeySessionName, c)
			if err != nil {
				s.logger.Error("failed to get session", "err", err)
				return echo.NewHTTPError(http.StatusBadRequest, "failed to get session")
			}
			s.logger.Debug("checking session", "userid", sess.Values[util.SessionKeyUserID])
			return next(c)
		}
	})
	e.Use(echomiddleware.RateLimiter(echomiddleware.NewRateLimiterMemoryStore(rate.Limit(20)))) // per second, per IP  //TODO: make it per username in session, not per IP or per account
	e.Use(echomiddleware.BodyLimit("2M"))
	e.Use(middleware.AddEchoContext)

	// admin routes
	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	e.GET("/healthz", s.Healthz)

	// http routes
	e.GET("/ping", s.Ping)
	e.GET("/favicon.ico", echo.NotFoundHandler)
	e.GET("/login", s.Login)
	e.POST("/signup", s.Signup)

	// graphql routes
	e.GET("/playground", echo.WrapHandler(playgroundHandler))
	e.POST("/query", echo.WrapHandler(graphqlHandler))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(s.HttpAddr); err != nil && err != http.ErrServerClosed {
			s.logger.Error("shutting down the server", "error", err)
		}
	}()

	<-ctx.Done()
	s.logger.Info("captured signal, gracefully shutting down the server with timeout")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return e.Shutdown(ctx)
}

func (s *serverCmd) Healthz(c echo.Context) error {
	return c.String(200, "OK")
}

func (s *serverCmd) Ping(c echo.Context) error {
	return c.String(200, "pong")
}

func (s *serverCmd) Signup(c echo.Context) error {
	// parse input
	inputUsername := c.FormValue("username")
	inputPassword := c.FormValue("password")
	inputAccountName := c.FormValue("accountname")
	s.logger.Debug("signup request", "username", inputUsername)
	if inputUsername == "" || inputPassword == "" || inputAccountName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username, password and accountname are required")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("bcrypt hash error", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}

	// create user and account in the database
	newAccount := &model.Account{
		Ulid: s.ulidManager.NewULID().String(),
		Name: inputAccountName,
	}
	err = s.db.Create(&newAccount).Error
	if err != nil {
		s.logger.Error("failed to create account", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	user := &model.User{
		Ulid:     s.ulidManager.NewULID().String(),
		Username: inputUsername,
		Password: string(hashedPassword),
		Account:  newAccount,
	}
	err = s.db.Create(&user).Error
	if err != nil {
		s.logger.Error("failed to create user", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}
	s.logger.Debug("created user", "username", user.Username)

	return nil
}

const OneHourSeconds = 3600   // 1 hour
const OneWeekSeconds = 604800 // 1 week

func (s *serverCmd) Login(c echo.Context) error {
	// parse input
	inputUsername := c.FormValue("username")
	inputPassword := c.FormValue("password")
	s.logger.Debug("login request", "username", inputUsername)
	if inputUsername == "" || inputPassword == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "username and password are required")
	}

	// find user in the database
	u := &model.User{}
	err := s.db.
		Preload("Account").
		Where("username = ?", inputUsername).
		First(u).Error
	if err != nil {
		s.logger.Debug("failed to find user", "err", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
	}

	// compare password hashes
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inputPassword))
	if err != nil {
		// Unauthorized
		s.logger.Debug("bcrypt compare error", "err", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid username or password")
	}

	// create session
	session, _ := s.store.Get(c.Request(), util.CookieKeySessionName) // this func returns an error when a session is created
	session.Options = &sessions.Options{
		MaxAge: OneWeekSeconds,
	}

	if u.Account != nil {
		session.Values[util.SessionKeyAccountID] = u.Account.Ulid
	}
	session.Values[util.SessionKeyUserID] = u.ID

	s.logger.Debug("found user", "username", u.ID, "account", u.Account.ID)

	err = session.Save(c.Request(), c.Response())
	if err != nil {
		s.logger.Error("failed to save session", "err", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save session")
	}

	return nil
}
