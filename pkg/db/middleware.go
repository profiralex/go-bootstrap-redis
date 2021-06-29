package db

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/go-redis/redis/v8"
	"net/http"
)

type redisClientContextKeyType string

const redisClientContextKey redisClientContextKeyType = "redis_client"

func SessionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := AddRedisClientToContext(r.Context())
		if err != nil {
			msg := fmt.Errorf("failed to add redis client to context: %w", err).Error()
			_ = render.Render(w, r, newSessionMiddlewareError(msg))
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type SessionMiddlewareError struct {
	message string
}

func (s SessionMiddlewareError) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusInternalServerError)
	return nil
}

func (s SessionMiddlewareError) Error() string {
	return s.message
}

func newSessionMiddlewareError(message string) *SessionMiddlewareError {
	return &SessionMiddlewareError{message: message}
}

func AddRedisClientToContext(ctx context.Context) (context.Context, error) {
	ctx = context.WithValue(ctx, redisClientContextKey, getRedisClient())
	return ctx, nil
}

func GetRedisClientFromContext(ctx context.Context) (*redis.Client, error) {
	sess, ok := ctx.Value(redisClientContextKey).(*redis.Client)
	if !ok {
		return nil, fmt.Errorf("missing *redis.Client in context")
	}

	return sess, nil
}
