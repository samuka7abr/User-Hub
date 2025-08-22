package httpport

import (
	"context"
	"net/http"
	"strings"

	"userhub/internal/security"
	"userhub/internal/shared/httputil"
)

type ctxUserID struct{}

type Router struct {
	Mux *http.ServeMux
	H   *Handlers
	JWT *security.JWT
}

func NewRouter(h *Handlers, jwt *security.JWT) *Router {
	mux := http.NewServeMux()
	r := &Router{Mux: mux, H: h, JWT: jwt}
	r.routes()
	return r
}

func (r *Router) routes() {
	r.Mux.HandleFunc("GET /healthz", r.H.Health)
	r.Mux.HandleFunc("POST /v1/auth/signup", r.H.Signup)
	r.Mux.HandleFunc("POST /v1/auth/login", r.H.Login)
	r.Mux.HandleFunc("GET /v1/users/me", r.auth(r.H.Me))

	r.Mux.HandleFunc("PATCH /v1/users/me/profile", r.auth(r.H.UpdateMyProfile))
	r.Mux.HandleFunc("GET /v1/users/{id}/profile", r.H.GetProfile)
}

func (r *Router) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		h := req.Header.Get("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			httputil.Error(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		sub, err := r.JWT.Parse(parts[1])
		if err != nil {
			httputil.Error(w, http.StatusUnauthorized, "invalid token")
			return
		}
		ctx := context.WithValue(req.Context(), ctxUserID{}, sub)
		next.ServeHTTP(w, req.WithContext(ctx))
	}
}
