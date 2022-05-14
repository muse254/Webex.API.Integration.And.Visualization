package api

import (
	"errors"
	"reflect"
	"testing"

	"Webex.API.Integration.And.Visualization/types"
)

func TestEncodeDecodeFromBase64(t *testing.T) {
	data := WebexAPIClient{
		ClientID:     "clientID",
		ClientSecret: "clientSecret",
		RedirectURI:  "redirectURI",
		Auth: types.AuthResponse{
			AccessToken:  "accessToken",
			RefreshToken: "refreshToken",
			ExpiresIn:    100,
		},
	}

	// encode to base64
	encStr, err := encodeToBase64(data)
	if err != nil {
		t.Errorf("encodeToBase64 failed: %v", err)
	}

	// decode from base64
	var decData WebexAPIClient
	if err := decodeFromBase64(&decData, encStr); err != nil {
		t.Errorf("decodeFromBase64 failed: %v", err)
	}

	if !reflect.DeepEqual(data, decData) {
		t.Errorf("decodeFromBase64 failed: %v", err)
	}
}

func TestGetWebexAPIClientCookie(t *testing.T) {
	sample := "username=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;"
	tests := []struct {
		name    string
		arg     string
		want    string
		wantErr error
	}{
		{
			name:    "get non-existent value",
			arg:     "WebexAPIClient",
			want:    "",
			wantErr: errors.New(`"WebexAPIClient" not found in cookies value`),
		},
		{
			name:    `fetch "username" value`,
			arg:     "username",
			want:    "",
			wantErr: nil,
		},
		{
			name:    `fetch "expires" value`,
			arg:     "expires",
			want:    "Thu, 01 Jan 1970 00:00:00 UTC",
			wantErr: nil,
		},
		{
			name:    `fetch "path" value`,
			arg:     "path",
			want:    "/",
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getCookieValue(sample, test.arg)
			if err != err {
				if !errors.Is(err, test.wantErr) {
					t.Errorf("wantErr: %v but got err: %v", err, test.wantErr)
				}
			}

			if test.want != got {
				t.Errorf("want: %s but got: %s", test.want, got)
			}
		})
	}
}
