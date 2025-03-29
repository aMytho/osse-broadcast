package server

import "net/http"

func validateUserToken(userID string, token string) bool {
	return true
}

func cors(h http.Handler, allowedOrigin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the origin matches the OSSE_HOST env var, we can allow the request.
		origin := r.Header.Get("origin")
		if origin == "http://"+allowedOrigin || origin == "https://"+allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		h.ServeHTTP(w, r)
	})
}
