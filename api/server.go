package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
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
		// The request will be like so: http://your-server.com/message?msg=<ErrorMsg>
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
	http.HandleFunc("/api", apiHandler)

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

		// request to Webex API to get the OAuth code
		if err := StartWebexAPIFlow(oauthReq.ClientID, oauthReq.ClientSecret, oauthReq.Scope, fmt.Sprintf("%s/auth", host)); err != nil {
			// redirect to error page
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, err.Error()), http.StatusSeeOther)
			return
		}

		// redirect to WebexAPI calls page

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
			http.Redirect(w, r, fmt.Sprintf("%s/error?msg=%s", host, "No OAuth code provided"), http.StatusSeeOther)
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
	}
}

// TODO: implement stop.
func apiHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func errorPage(w io.Writer, errorMsg string) error {
	tmpl, _ := template.ParseFiles("./templates/generic_template.html")
	return tmpl.Execute(w, GenericPage{
		Heading:          "Error",
		Message:          errorMsg,
		ShowHomeRedirect: true,
	})
}

// messagePage is the page that displays the messages also allowing for redirect to the
// apiCalls page
func messagePage(w io.Writer, message string, apiRedirect bool) error {
	tmpl, _ := template.ParseFiles("./templates/generic_template.html")
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
