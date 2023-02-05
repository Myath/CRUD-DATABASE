package handler

import (
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	sessionManager *scs.SessionManager
	decoder        *form.Decoder
	db             *sqlx.DB
}

const(
	adminLoginPath = "/adminLogin"
)

func NewHandler(sm *scs.SessionManager, fromdecoder *form.Decoder, db *sqlx.DB) *chi.Mux {
	h := &Handler{
		sessionManager: sm,
		decoder:        fromdecoder,
		db:             db,
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(VerbMethod)

	r.Group(func(r chi.Router) {
		r.Use(h.AuthenticationForLogin)
		r.Get(adminLoginPath, h.AdminLogin)
		r.Post(adminLoginPath, h.AdminLoginProcess)
	})

	r.Get("/adminlogout", h.AdminLogOut)

	r.Route("/student", func(r chi.Router) {

		r.Use(h.Authentication)

		r.Get("/list", h.StudentsList)

		r.Get("/create", h.CreateStudent)

		r.Post("/store", h.StudentStore)

		r.Get("/{id:[0-9]+}/edit", h.StudentEdit)

		r.Put("/{id:[0-9]+}/update", h.StudentUpdate)

		r.Get("/{id:[0-9]+}/delete", h.DeleteStudent)
	})

	return r
}

func VerbMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch strings.ToLower(r.PostFormValue("_method")) {
			case "put":
				r.Method = http.MethodPut
			case "patch":
				r.Method = http.MethodPatch
			case "delete":
				r.Method = http.MethodDelete
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h Handler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := h.sessionManager.GetString(r.Context(), "username")
		if username == "" {
			http.Redirect(w, r, "/adminLogin", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h Handler) AuthenticationForLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := h.sessionManager.GetString(r.Context(), "username")
		if username != "" {
			http.Redirect(w, r, "/student/list", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

