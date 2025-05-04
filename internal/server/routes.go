package server

import (
	"Auth/internal/database"
	"Auth/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/markbates/goth/gothic"
)

const REFRESH_POLL_INTERVAL = 15

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api", func(api chi.Router) {
		api.Use(s.AuthMiddleware)

		api.Get("/me", s.meHandler)
		api.Get("/users", s.usersHandler)
		api.Put("/users", s.putUsersHandler)
		api.Post("/users", s.postUsersHandler)
		api.Delete("/users/{id}", s.deleteUsersHandler)
	})

	r.Route("/health", func(api chi.Router) {
		api.Use(s.AuthMiddleware)

		r.Get("/", s.healthHandler)
	})

	r.Get("/auth/{provider}/callback", s.getAuthCallbackFunction)
	r.Get("/auth/{provider}", gothic.BeginAuthHandler)
	r.Get("/auth/logout", s.logoutHandler)

	return r
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

func (s *Server) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Println("OAuth error:", err)
		fmt.Fprintln(w, r)
		return
	}

	dbService := database.New()

	isValid, isAdmin, err := dbService.CheckUser(user.Email)
	if err != nil {
		http.Error(w, "Error checking user", http.StatusUnauthorized)
		return
	}

	session, _ := s.sessionStore.Get(r, "auth-session")
	session.Values["user_id"] = user.UserID
	session.Values["user_email"] = user.Email
	session.Values["is_valid"] = isValid
	session.Values["is_admin"] = isAdmin
	session.Values["cached_at"] = time.Now()
	session.Save(r, w)

	http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
}

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.sessionStore.Get(r, "auth-session")

		cachedAt, _ := session.Values["cached_at"].(time.Time)
		email, _ := session.Values["cached_at"].(string)

		if time.Since(cachedAt) > REFRESH_POLL_INTERVAL*time.Minute {
			dbService := database.New()

			isValid, isAdmin, err := dbService.CheckUser(email)
			if err != nil {
				http.Error(w, "Error checking user", http.StatusUnauthorized)
				return
			}
			session.Values["is_valid"] = isValid
			session.Values["is_admin"] = isAdmin
			session.Values["cached_at"] = time.Now()
			session.Save(r, w)
		}

		if !session.Values["is_valid"].(bool) {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		if session.Values["user_id"] == nil {
			http.Error(w, "Access denied", http.StatusForbidden)
			// http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessionStore.Get(r, "auth-session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) meHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessionStore.Get(r, "auth-session")
	userID, ok := session.Values["user_id"].(string)
	userEmail, _ := session.Values["user_email"].(string)
	isValid, _ := session.Values["is_valid"].(bool)
	isAdmin, _ := session.Values["is_admin"].(bool)

	if !ok || userID == "" {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user_id":    userID,
		"user_email": userEmail,
		"is_valid":   isValid,
		"is_admin":   isAdmin,
	})
}

func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := s.db.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (s *Server) putUsersHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.db.EditUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) postUsersHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err = s.db.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *Server) deleteUsersHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := s.db.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
