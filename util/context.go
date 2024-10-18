package util

const SessionKeyCookieKey = "session"
const SessionKeyAccountID = "account_id"
const SessionKeyUserID = "user_id"

var CtxKeyEchoContext = &contextKey{"echoContext"}

type contextKey struct {
	name string
}
