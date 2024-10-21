package util

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const CookieKeySessionName = "session"

const SessionKeyAccountID = "account_id"
const SessionKeyUserID = "user_id"

var CtxKeyEchoContext = &contextKey{"echoContext"}

type contextKey struct {
	name string
}

func ExtractEchoContext(ctx context.Context) (echo.Context, error) {
	echoCtx := ctx.Value(CtxKeyEchoContext)
	if echoCtx == nil {
		return nil, errors.New("No echo context in context")
	}
	ec, ok := echoCtx.(echo.Context)
	if !ok {
		return nil, errors.New("Echo context in context is not an echo.Context")
	}
	return ec, nil
}
