package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	AUTH_URL     = "https://webexapis.com/v1/authorize"
	TOKEN_URL    = "https://webexapis.com/v1/access_token"
	BASE_API_URL = "https://webexapis.com/v1"
)

// WebexAPIClient is a convenience wrapper that will be used to make API calls to the Webex API.
type WebexAPIClient struct {
	clientID     string
	clientSecret string
	redirectURI  string
	auth         AuthResponse
}

// When the user successfully authorizes the application, the OAuth code is retrieved from the redirect handler and
// used in creating the WebexAPIClient.
func NewWebexAPIClient(OAuthCode, clientID, clientSecret, redirectURI string) (*WebexAPIClient, error) {
	// using data form-urlencoded
	data := url.Values{}
	data.Add("grant_type", "authorization_code")
	data.Add("code", OAuthCode)
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)
	data.Add("redirect_uri", redirectURI)

	// retrive the access token using the OAuth code to verify the user's identity
	req, err := http.NewRequest(http.MethodPost, TOKEN_URL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// when the OAuth provide is invalid, the response will be a 401 error
	if resp.StatusCode == http.StatusUnauthorized {
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
		auth:         authResp,
	}, nil
}

// ListMeeting lists all meenting that are accessible to the client account.
func (c *WebexAPIClient) ListMeetings(tries int) ([]MeetingSeries, error) {
	if tries > 3 {
		return nil, fmt.Errorf("failed to get meetings")
	}

	req, err := http.NewRequest(http.MethodGet, BASE_API_URL+"/meetings", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	{

	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Authorization", "Bearer "+c.auth.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		{
			var meetings []MeetingSeries
			if err := json.NewDecoder(resp.Body).Decode(&meetings); err != nil {
				return nil, err
			}

			return meetings, nil
		}

	case http.StatusUnauthorized:
		if err = c.refreshToken(); err != nil {
			return nil, err
		}
		return c.ListMeetings(tries + 1)

	case http.StatusNoContent:
		return nil, nil

	default:
		return nil, fmt.Errorf("failed to get meetings")
	}
}

// When the access_token expires or is invalid, the refresh token is used to generate a new access token.
func (c *WebexAPIClient) refreshToken() error {
	data, err := json.Marshal(RefreshTokenRequest{
		GrantType:    "refresh_token",
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
		RefreshToken: c.auth.RefreshToken,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, TOKEN_URL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// when the refresh token is expired, the response will be a 400 error
	if resp.StatusCode == http.StatusBadRequest {
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
	c.auth = authResp
	return nil
}
