package httpport

import (
	"encoding/json"
	"net/http"

	"userhub/internal/app"
	"userhub/internal/shared/httputil"
)

type Handlers struct{ Svc *app.Service }

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type signupReq struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

func (h *Handlers) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	u, err := h.Svc.Signup(req.Email, req.Password)
	if err != nil {
		switch err {
		case app.ErrInvalidInput:
			httputil.Error(w, http.StatusBadRequest, "invalid email or password")
		case app.ErrEmailExists:
			httputil.Error(w, http.StatusConflict, "email already exists")
		default:
			httputil.Error(w, http.StatusInternalServerError, "signup error")
		}
		return
	}
	httputil.WriteJSON(w, http.StatusCreated, map[string]any{
		"user": map[string]any{"id": u.ID, "email": u.Email, "created_at": u.CreatedAt},
	})
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	at, rt, err := h.Svc.Login(req.Email, req.Password)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"access_token": at, "refresh_token": rt})
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	uid, _ := r.Context().Value(ctxUserID{}).(string)
	u, ok := h.Svc.GetUser(uid)
	if !ok {
		httputil.Error(w, http.StatusNotFound, "user not found")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"id": u.ID, "email": u.Email, "created_at": u.CreatedAt})
}

type profileUpdateReq struct {
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

func (h *Handlers) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	uid, _ := r.Context().Value(ctxUserID{}).(string)
	var req profileUpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	p, err := h.Svc.UpdateProfile(uid, req.Name, req.Bio, req.AvatarURL)
	if err != nil {
		if err == app.ErrNotFound {
			httputil.Error(w, http.StatusNotFound, "user not found")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, "update error")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"profile": p})
}

func (h *Handlers) GetProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	p, err := h.Svc.GetProfile(id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, "not found")
		return
	}
	httputil.WriteJSON(w, http.StatusOK, map[string]any{"profile": p})
}
