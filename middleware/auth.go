package middleware

import (
	"net/http"

	"github.com/zacscodingclub/utube-tut/sessions"
)

func AuthRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		user_id, ok := session.Values["user_id"]
		if !ok || user_id == nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
		h.ServeHTTP(w, r)
	}
}
