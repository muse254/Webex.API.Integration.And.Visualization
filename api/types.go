package api

// AuthResponse is returned on successful authorization. The access token is to be used in susbsequent requests.
type AuthResponse struct {
	AccessToken           string `json:"access_token"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

// HTTP4XXError is returned when the access token is expired or invalid.
// To recover, a new access token must be generated using the refresh token.
type HTTP4XXError struct {
	Message string `json:"message"`
	Errors  []struct {
		Description string `json:"description"`
	} `json:"errors"`
	TrackingID string `json:"trackingId"`
}

type OAuthRequest struct {
	ClientID     string
	ClientSecret string
	Scope        string
}

type GenericPage struct {
	Heading          string
	Message          string
	ShowHomeRedirect bool
	ShowAPIRedirect  bool
}
