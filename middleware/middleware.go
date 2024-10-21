package middleware

import (
	"context"

	"github.com/labstack/echo/v4"

	"github.com/mwasilew2/echo-gqlgen-casbin-rbac-example/util"
)

func AddEchoContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		echoCtx := context.WithValue(c.Request().Context(), util.CtxKeyEchoContext, c)
		c.SetRequest(c.Request().WithContext(echoCtx))
		return next(c)
	}
}
