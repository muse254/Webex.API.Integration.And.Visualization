package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"Webex.API.Integration.And.Visualization/persist"
	"Webex.API.Integration.And.Visualization/types"
)

const (
	AUTH_URL  = "https://webexapis.com/v1/authorize"
	TOKEN_URL = "https://webexapis.com/v1/access_token"
)

// WebexAPIClient is a convenience wrapper that will be used to make API calls to the Webex API.
// It holds the client_id, client_secret, redirect_uri, and access_token required for API calls.
// Values are also bound to the client and saved as a cookie.
type WebexAPIClient struct {
	ClientID     string             `json:"client_id"`
	ClientSecret string             `json:"client_secret"`
	RedirectURI  string             `json:"redirect_uri"`
	Auth         types.AuthResponse `json:"auth"`
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
		var errResponse types.HTTP4XXError
		if json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%s: %s", errResponse.Message, errResponse.Errors[0].Description)
	}

	// parse the response body
	var authResp types.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return nil, err
	}

	return &WebexAPIClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Auth:         authResp,
	}, nil
}

// ListMeeting lists all meetings that are accessible to the client account.
func (c *WebexAPIClient) ListMeetings(tries int) (*types.MeetingsList, error) {
	if tries > 3 {
		return nil, fmt.Errorf("failed to get meetings from API, StatusCode: StatusUnauthorized")
	}

	req, err := http.NewRequest(http.MethodGet, "https://webexapis.com/v1/meetings", nil)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	from := now.AddDate(0, 0, -30)

	req.URL.RawQuery = (url.Values{
		"to":          []string{now.Format(time.RFC3339)},
		"from":        []string{from.Format(time.RFC3339)},
		"meetingType": []string{"meeting"},
	}).Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		{
			var meetings types.MeetingsList
			if err := json.NewDecoder(resp.Body).Decode(&meetings); err != nil {
				return nil, err
			}

			return &meetings, nil
		}

	case http.StatusUnauthorized:
		if err = c.refreshToken(); err != nil {
			return nil, err
		}
		return c.ListMeetings(tries + 1)

	case http.StatusNoContent:
		return nil, nil

	default:
		return nil, fmt.Errorf("failed to get meetings from API, StatusCode: %s", resp.Status)
	}
}

// GetMeetingQualities gets the qualities of a meeting.
func (c *WebexAPIClient) GetMeetingQualities(db *persist.Persist, meetingID string, tries int) (*types.MeetingQualities, error) {
	if tries > 3 {
		return nil, fmt.Errorf("failed to get meeting quality from API, StatusCode: StatusUnauthorized")
	}

	req, err := http.NewRequest(http.MethodGet, "https://analytics.webexapis.com/v1/meeting/qualities", nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = (url.Values{
		"meetingId": []string{meetingID},
	}).Encode()
	req.Header.Add("Authorization", "Bearer "+c.Auth.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var meetingQualities types.MeetingQualities
		if err := json.Unmarshal(buf, &meetingQualities); err != nil {
			return nil, err
		}

		// if the meeting qualities fetch was a success and decoding was successful, persist to db
		if err := db.SaveAnalyticsData(meetingID, c.ClientID, string(buf)); err != nil {
			log.Printf("error on SaveAnalyticsData(): %s\n", err.Error())
		}

		return &meetingQualities, nil

	case http.StatusTooManyRequests: // this StatusCode is hit when the rate limit of 1 request per 5 minutes is hit
		// retrieve meeting qualities from persitance storage.
		data, err := db.RetriveAnalyticsData(c.ClientID, meetingID)
		if err != nil {
			if err == sql.ErrNoRows {
				// data was never persisted due to error in the parsing process
				return nil, fmt.Errorf("failed to get meeting qualities from API, StatusCode: %s, try after 5 min", resp.Status)
			}
			return nil, err
		}

		return data, nil

	case http.StatusUnauthorized:
		if err = c.refreshToken(); err != nil {
			return nil, err
		}
		return c.GetMeetingQualities(db, meetingID, tries+1)

	case http.StatusNoContent:
		return nil, nil

	default:
		return nil, fmt.Errorf("failed to get meeting qualities from API, StatusCode: %s", resp.Status)
	}
}

// When the access_token expires or is invalid, the refresh token is used to generate a new access token.
func (c *WebexAPIClient) refreshToken() error {
	data, err := json.Marshal(types.RefreshTokenRequest{
		GrantType:    "refresh_token",
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RefreshToken: c.Auth.RefreshToken,
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
		var errResponse types.HTTP4XXError
		if json.NewDecoder(resp.Body).Decode(&errResponse); err != nil {
			return err
		}

		return fmt.Errorf("%s: %s", errResponse.Message, errResponse.Errors[0].Description)
	}

	// parse the response body
	var authResp types.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	// update the client
	c.Auth = authResp
	return nil
}
