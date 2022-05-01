package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func RedirectServer() error {
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
	http.HandleFunc("/api/get_meetings", getMeetings(host))

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
		oauthReq := OAuthRequest{
			ClientID:     strings.TrimSpace(r.FormValue("client_id")),
			ClientSecret: strings.TrimSpace(r.FormValue("client_secret")),
			Scope:        strings.TrimSpace(r.FormValue("scope")),
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
		// log the complete path
		fmt.Println(r.URL.String())

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
		var oauthReq OAuthRequest
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

func getMeetings(host string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { // check where the cookie exists for auth response, if not redirect to auth page
		cookie, err := r.Cookie("WebexAPIClient")
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, "Complete the authentication flow."), http.StatusSeeOther)
			return
		}

		// get_meetings API call
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

		// render the meetings page, pretty print the meetings as json
		data, err := json.MarshalIndent(struct {
			Meetings []MeetingSeries
		}{
			Meetings: meetings,
		}, "", "\t")
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		t, _ := template.ParseFiles("./templates/get_meetings.html")
		t.Execute(w, string(data))
	}
}

func errorPage(w io.Writer, errorMsg string) error {
	tmpl, _ := template.ParseFiles("./templates/generic_page.html")
	return tmpl.Execute(w, GenericPage{
		Heading:          "Error",
		Message:          errorMsg,
		ShowHomeRedirect: true,
	})
}

// messagePage is the page that displays the messages also allowing for redirect to the
// apiCalls page
func messagePage(w io.Writer, message string, apiRedirect bool) error {
	tmpl, _ := template.ParseFiles("./templates/generic_page.html")
	return tmpl.Execute(w, GenericPage{
		Heading:          "Message",
		Message:          message,
		ShowAPIRedirect:  apiRedirect,
		ShowHomeRedirect: !apiRedirect,
	})
}

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

func decodeFromBase64(v interface{}, enc string) error {
	return json.NewDecoder(base64.NewDecoder(base64.StdEncoding, strings.NewReader(enc))).Decode(v)
}
