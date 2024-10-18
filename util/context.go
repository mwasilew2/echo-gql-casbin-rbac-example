package util

var CtxKeyAccountId = &contextKey{"accountID"}

type contextKey struct {
	name string
}
