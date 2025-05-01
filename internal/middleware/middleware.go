package middleware

import (
	"context"
	"net/http"
	"strings"

	"fem-go-crud/internal/auth"
	"fem-go-crud/internal/store"
	"fem-go-crud/internal/utils"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

type contextKey string

const (
	UserContextKey contextKey = "user"
)

func NewUserMiddleware(us store.UserStore) *UserMiddleware {
	return &UserMiddleware{
		UserStore: us,
	}
}

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		// bad actor call
		panic("missing user in context")
	}

	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ") // "Bearer <token>"
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			_ = utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid authorization header"})
			return
		}

		token := headerParts[1]
		user, err := um.UserStore.GetUserFromToken(token, auth.TokenScopeAuth)
		if err != nil || user == nil {
			_ = utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid or expired token"})
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user == nil || user.IsAnonymous() {
			_ = utils.WriteJSONResponse(w, http.StatusUnauthorized, utils.Envelope{"error": "unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
