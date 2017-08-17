package server

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/hectane/hectane/db"
)

const (
	userNone   = "none"
	userNormal = "normal"
	userAdmin  = "admin"

	sessionName   = "session"
	sessionUserID = "user_id"

	contextUser   = "user"
	contextParams = "params"
)

// post ensures that the request method is POST.
func (s *Server) post(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h.ServeHTTP(w, r)
	}
}

// auth ensures that the user is logged in and adds the current user to the
// context.
func (s *Server) auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u *db.User
		session, _ := s.sessions.Get(r, sessionName)
		v, ok := session.Values[sessionUserID]
		if !ok {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		u, err := db.FindUser(db.DefaultToken, "ID", v)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), contextUser, u)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

// json attempts to parse the request body into a struct of the provided type.
func (s *Server) json(h http.HandlerFunc, v interface{}) http.HandlerFunc {
	t := reflect.TypeOf(v)
	return func(w http.ResponseWriter, r *http.Request) {
		v := reflect.New(t).Interface()
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), contextParams, v)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
