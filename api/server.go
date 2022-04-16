package api

import (
	"net/http"
)

type Server struct {
	OAuthCode string
}

func (s *Server) RedirectServer() error {
	// This is the redirect URL that captures the OAuth code from the user's authentication.
	// The request will be like so: http://your-server.com/auth?code=<OAuthCode>
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		s.OAuthCode = r.URL.Query().Get("code")
		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(":3000", nil)
}
