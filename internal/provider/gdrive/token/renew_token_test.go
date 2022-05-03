package token

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewToken_RenewToken(t *testing.T) {
	nToken := NewToken{
		TokenUrl:     "/apis/token",
		ClientID:     "client",
		ClientSecret: "secret",
		RefreshToken: "asd",
	}

	// prepare fake server to mimic Google apis
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req.ParseForm()

		clID := req.PostForm.Get("client_id")
		clSecret := req.PostForm.Get("client_secret")
		rToken := req.PostForm.Get("refresh_token")
		gType := req.PostForm.Get("grant_type")

		assert.Equal(t, nToken.ClientID, clID)
		assert.Equal(t, nToken.ClientSecret, clSecret)
		assert.Equal(t, nToken.RefreshToken, rToken)
		assert.Equal(t, "refresh_token", gType)

		response := struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
		}{AccessToken: "token", TokenType: "bearer"}

		b, err := json.Marshal(&response)
		require.NoError(t, err, "failed to marshal json response")

		rw.WriteHeader(http.StatusOK)
		rw.Write(b)
	}))

	t.Run("Should error because using wrong or invalid server url", func(t *testing.T) {
		_, err := nToken.RenewToken(server.Client())
		require.Error(t, err)
	})

	t.Run("Should pass because using the right server url", func(t *testing.T) {
		nToken.TokenUrl = server.URL
		_, err := nToken.RenewToken(server.Client())
		require.NoError(t, err)
	})

	// prepare fake server to mimic Google apis
	server2 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("content-type", "applications/xml")
	}))

	t.Run("Should error because the fake server return non json response body", func(t *testing.T) {
		nToken.TokenUrl = server2.URL
		_, err := nToken.RenewToken(server2.Client())
		require.Error(t, err)
	})

	// cleanup test env
	t.Cleanup(func() {
		server.Close()
		server2.Close()
	})
}
