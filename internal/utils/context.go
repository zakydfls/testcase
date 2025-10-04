package utils

type contextKey string

const (
	UserIDContextKey   contextKey = "user_id"
	UsernameContextKey contextKey = "username"
	RoleContextKey     contextKey = "role"
	EmailContextKey    contextKey = "email"
	IsActiveContextKey contextKey = "is_active"
)

const (
	RequestIDContextKey contextKey = "request_id"
	TraceIDContextKey   contextKey = "trace_id"
	SessionIDContextKey contextKey = "session_id"
	IPAddressContextKey contextKey = "ip_address"
	UserAgentContextKey contextKey = "user_agent"
)
