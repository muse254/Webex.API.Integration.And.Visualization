package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	AUTH_URL     = "https://webexapis.com/v1/authorize"
	TOKEN_URL    = "https://webexapis.com/v1/access_token"
	BASE_API_URL = "https://analytics.webexapis.com/v1"
)

type WebexAPIClient struct {
	clientID     string
	clientSecret string
	redirectURI  string
	auth         *AuthResponse
}

// StartWebexAPIFlow inititiates the OAuth flow requesting the user for permissions in order to
// access the Webex API.
func StartWebexAPIFlow(clientID, clientSecret, scope, redirectURI string) error {
	// check the user credentails for authorization to interact with the API
	req, err := http.NewRequest(http.MethodGet, AUTH_URL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", clientID)
	q.Add("redirect_uri", redirectURI)
	q.Add("scope", scope)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// When the user successfully authorizes the application, the OAuth code is retrieved from the redirect handler and
// used in creating the WebexAPIClient.
func NewWebexAPIClient(OAuthCode, clientID, clientSecret, redirectURI string) (*WebexAPIClient, error) {
	// retrive the access token using the OAuth code to verify the user's identity
	req, err := http.NewRequest(http.MethodGet, TOKEN_URL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("code", OAuthCode)
	q.Add("redirect_uri", redirectURI)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// when the OAuth provide is invalid, the response will be a 401 error
	if resp.StatusCode != http.StatusUnauthorized {
		var errResponse HTTP4XXError
		if json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s", errResponse.Message, errResponse.Errors[0].Description)
	}

	// parse the response body
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &WebexAPIClient{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,
		auth:         &authResp,
	}, nil
}

func (c *WebexAPIClient) GetMeetings() error {
	req, err := http.NewRequest(http.MethodGet, BASE_API_URL+"/meetings", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	{

	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+c.auth.AccessToken)

	return nil
}

// When the access_token expires or is invalid, the refresh token is used to generate a new access token.
func (c *WebexAPIClient) refreshToken() error {
	req, err := http.NewRequest(http.MethodGet, TOKEN_URL, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("grant_type", "refresh_token")
	q.Add("client_id", c.clientID)
	q.Add("client_secret", c.clientSecret)
	q.Add("refresh_token", c.auth.RefreshToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// when the refresh token is expired, the response will be a 400 error
	if resp.StatusCode != http.StatusBadRequest {
		var errResponse HTTP4XXError
		if json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return err
		}

		return fmt.Errorf("%s: %s", errResponse.Message, errResponse.Errors[0].Description)
	}

	// parse the response body
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	// update the client
	c.auth = &authResp
	return nil
}
