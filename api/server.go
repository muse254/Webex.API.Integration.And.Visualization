package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/template"

	"Webex.API.Integration.And.Visualization/persist"
	"Webex.API.Integration.And.Visualization/types"
)

// WebexApplicationServer is the server for the Webex Application.
func WebexApplicationServer(db *persist.Persist) error {
	// load the server's host
	host := os.Getenv("HOST")
	if host == "" {
		return fmt.Errorf("HOST environment variable is not set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.html")
	})
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		// The request will be like so: http://your-server.com/error?msg=<ErrorMsg>
		errorMsg := r.URL.Query().Get("msg")
		if errorMsg == "" {
			errorMsg = "Unknown error"
		}

		if err := errorPage(w, errorMsg); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		// The request will be like so: http://your-server.com/message?msg=<Msg>
		msg := r.URL.Query().Get("msg")
		apiRedirect := true
		if msg == "" {
			msg = "Unknown message"
			apiRedirect = false
		}
		if err := messagePage(w, msg, apiRedirect); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	http.HandleFunc("/init", init_flow(host))

	// "/auth" is called by Webex on redirect from the OAuth flow.
	http.HandleFunc("/auth", auth(host))

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		// check if cookie exists for API calls
		_, err := r.Cookie("WebexAPIClient")
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		// display all APIs calls page
		http.ServeFile(w, r, "./templates/api_calls.html")
	})
	http.HandleFunc("/get_meetings_page", getMeetings(host))
	http.HandleFunc("/get_analytics_page", analyticsVisualization(db, host))
	http.HandleFunc("/get_analytics_file", dowloadAnalyticsFile(db, host))
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	return http.ListenAndServe(":3000", nil)
}

// init_flow initializes the Oauth Flow for the application
func init_flow(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// parse the request form
		r.ParseForm()
		oauthReq := types.OAuthRequest{
			ClientID:     strings.TrimSpace(r.FormValue("client_id")),
			ClientSecret: strings.TrimSpace(r.FormValue("client_secret")),
			Scope:        "analytics:read_all meeting:schedules_read",
		}

		// create cookie for later reference
		oauthReqStr, err := encodeToBase64(oauthReq)
		if err != nil {
			// redirect to error page
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		cookie := &http.Cookie{
			Name:  "OAuthRequest",
			Value: oauthReqStr,
		}
		// The server is stateless and a dabatabase is not required  because the OAuthCode is valid for small period of time and client-bound.
		http.SetCookie(w, cookie)

		// redirect to Webex, calling the auth endpoint
		u, err := url.Parse(AUTH_URL)
		if err != nil {
			// redirect to error page
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}
		q := u.Query()
		q.Add("response_type", "code")
		q.Add("client_id", oauthReq.ClientID)
		q.Add("redirect_uri", fmt.Sprintf("%s/auth", host))
		q.Add("scope", oauthReq.Scope)
		q.Add("state", "some state")
		u.RawQuery = q.Encode()
		http.Redirect(w, r, u.String(), http.StatusSeeOther)
	}
}

// auth is the redirect URL that captures the OAuth code from the user's authentication.
// The request will be like so: http://your-server.com/auth?code=<OAuthCode>
func auth(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			// redirect to error page
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, "No OAuth code provided"), http.StatusSeeOther)
			return
		}

		// Retrive the OAuth request from the cookie
		cookie, err := r.Cookie("OAuthRequest")
		if err != nil {
			// redirect to error page
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		// Decode the OAuth code from the cookie
		var oauthReq types.OAuthRequest
		if err := decodeFromBase64(&oauthReq, cookie.Value); err != nil {
			// redirect to error page
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		// use the code to create a WebexAPIClient
		client, err := NewWebexAPIClient(code, oauthReq.ClientID, oauthReq.ClientSecret, fmt.Sprintf("%s/auth", host))
		if err != nil {
			// redirect to error page
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		// save the client as a cookie
		clientStr, err := encodeToBase64(client)
		if err != nil {
			// redirect to error page
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}
		cookie = &http.Cookie{
			Name:  "WebexAPIClient",
			Value: clientStr,
		}
		http.SetCookie(w, cookie)

		// redirect to message page with option for API redirect
		http.Redirect(w, r, fmt.Sprintf("%s/message?msg=%s", host, "Successfully authenticated"), http.StatusSeeOther)
	}
}

// getMeetings is the handler for the /get_meetings_page endpoint.
func getMeetings(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// check where the cookie exists from client, if not redirect to error page
		cookie, err := r.Cookie("WebexAPIClient")
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, "Complete the authentication flow."), http.StatusSeeOther)
			return
		}

		// get WebexAPIClient from cookie
		var client WebexAPIClient
		if err := decodeFromBase64(&client, cookie.Value); err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		meetings, err := client.ListMeetings(0)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		// render the page with data provided
		t, _ := template.ParseFiles("./templates/get_meetings.html")
		t.Execute(w, meetings)
	}
}

func analyticsVisualization(db *persist.Persist, host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// id := r.URL.Query().Get("id")
		// if id == "" {
		// 	http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, "No meeting ID provided"), http.StatusSeeOther)
		// 	return
		// }

		// dp := r.URL.Query().Get("dp")
		// if dp == "" {
		// 	dp = "audio_in"
		// }

		// qualities, errUrl := analyticsCommonfetch(r, db, id, host)
		// if errUrl != "" {
		// 	http.Redirect(w, r, errUrl, http.StatusSeeOther)
		// 	return
		// }

		// chartData, err := types.GetVisualData(qualities, dp)
		// if err != nil {
		// 	http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
		// 	return
		// }

		// data, _ := json.Marshal(chartData)
		// t, err := template.ParseFiles("./templates/analytics_visualization.html")
		// if err != nil {
		// 	http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
		// 	return
		// }

		// if err = t.Execute(w, struct {
		// 	StrData string
		// 	Actual  *types.VisualData
		// }{
		// 	StrData: string(data),
		// 	Actual:  chartData,
		// }); err != nil {
		// 	http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
		// 	return
		// }

		w.Write([]byte("Work in Progress"))
	}
}

func dowloadAnalyticsFile(db *persist.Persist, host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, "No meeting ID provided"), http.StatusSeeOther)
			return
		}

		qualities, errUrl := analyticsCommonfetch(r, db, id, host)
		if errUrl != "" {
			http.Redirect(w, r, errUrl, http.StatusSeeOther)
			return
		}

		// transform to visual data
		visualData, err := types.GetAllVisualData(qualities)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, "Internal Error"), http.StatusSeeOther)
			return
		}

		// file data
		analytics := struct {
			Analytics []types.VisualData `json:"analytics"`
		}{visualData}

		// pretty print the qualities as json
		data, err := json.Marshal(analytics)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		// write the data as a binary stream to client that will be donloaded as file
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=analytics_%s.json", qualities.MeetingID))
		w.Header().Add("Content-Type", "application/octet-stream")
		w.Write(data)
	}
}

func analyticsCommonfetch(r *http.Request, db *persist.Persist, id, host string) (*types.MeetingQualities, string) {
	// check where the cookie exists from client, if not redirect to error page
	cookie, err := r.Cookie("WebexAPIClient")
	if err != nil {
		return nil, fmt.Sprintf("%s/error?msg=%s", host, "Complete the authentication flow.")
	}

	// get WebexAPIClient from cookie
	var client WebexAPIClient
	if err := decodeFromBase64(&client, cookie.Value); err != nil {
		return nil, fmt.Sprintf("%s/error?msg=%s", host, err.Error())
	}

	// fetch analytics data
	qualities, err := client.GetMeetingQualities(db, id, 0)
	if err != nil {
		return nil, fmt.Sprintf("%s/error?msg=%s", host, err.Error())
	}

	qualities.MeetingID = id
	return qualities, ""
}

// errorPage is the error page that is displayed when an error occurs.
func errorPage(w io.Writer, errorMsg string) error {
	tmpl, _ := template.ParseFiles("./templates/generic_page.html")
	return tmpl.Execute(w, types.GenericPage{
		Heading:          "Error",
		Message:          errorMsg,
		ShowHomeRedirect: true,
	})
}

// messagePage is the page that displays the messages also allowing for redirect to the
// apiCalls page
func messagePage(w io.Writer, message string, apiRedirect bool) error {
	tmpl, _ := template.ParseFiles("./templates/generic_page.html")
	return tmpl.Execute(w, types.GenericPage{
		Heading:          "Message",
		Message:          message,
		ShowAPIRedirect:  apiRedirect,
		ShowHomeRedirect: !apiRedirect,
	})
}

// encodeToBase64 encodes a non-pointer type to a base64 string
func encodeToBase64(v interface{}) (string, error) {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	err := json.NewEncoder(encoder).Encode(v)
	if err != nil {
		return "", err
	}
	encoder.Close()
	return buf.String(), nil
}

// decodeFromBase64 decodes a base64 string to a non-pointer type
func decodeFromBase64(v interface{}, enc string) error {
	return json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(enc))).Decode(v)
}

// getWebexAPIClientCookie parses the cookies string and retrives the "cookieName" key cookie value
func getCookieValue(cookies, cookieName string) (string, error) {
	// Sample cookies string: "username=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
	for _, cookie := range strings.Split(cookies, ";") {
		// if prefix is "cookieName", split at "=" & return RHS
		cookie = strings.TrimSpace(cookie)
		if strings.HasPrefix(cookie, cookieName) {
			return strings.Split(cookie, "=")[1], nil
		}
	}

	return "", fmt.Errorf(`"%s" not found in cookies value`, cookieName)
}
