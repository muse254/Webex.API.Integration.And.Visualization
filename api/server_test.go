package api

import (
	"reflect"
	"testing"
)

func TestEncodeDecodeFromBase64(t *testing.T) {
	data := WebexAPIClient{
		ClientID:     "clientID",
		ClientSecret: "clientSecret",
		RedirectURI:  "redirectURI",
		Auth: AuthResponse{
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
