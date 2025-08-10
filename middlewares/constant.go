package middlewares

type contextKey string

const (
	usernameKey contextKey = "username"
	roleKey     contextKey = "role"
)
