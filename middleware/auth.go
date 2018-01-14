package middleware

import (
	"fmt"
	"net/http"

	"github.com/zacscodingclub/utube-tut/sessions"
)

func AuthRequired(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := sessions.Store.Get(r, "session")
		fmt.Sprintln(session.Values["username"])
		_, ok := session.Values["username"]
		if !ok {
			http.Redirect(w, r, "/login", 302)
			return
		}
		h.ServeHTTP(w, r)
	}
}
