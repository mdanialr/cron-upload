package gdrive

import (
	"encoding/json"
	"github.com/mdanialr/cron-upload/internal/config"
	"github.com/mdanialr/cron-upload/internal/provider/gdrive/token"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestInitToken(t *testing.T) {
	// START PREPARE FAKE SERVER to MIMIC GOOGLE APIS SERVER
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			AccessToken string `json:"access_token"`
			TokenType   string `json:"token_type"`
			Expiry      string `json:"expiry"`
		}{AccessToken: "access", TokenType: "Bearer", Expiry: "0001-01-01T00:00:00Z"}

		js, _ := json.Marshal(&resp)

		w.WriteHeader(http.StatusOK)
		w.Write(js)
	}))
	// DONE PREPARE FAKE SERVER to MIMIC GOOGLE APIS SERVER
	cl := &http.Client{}

	t.Run("Should fail when using invalid or non exist auth.json file", func(t *testing.T) {
		mod := &config.Model{Provider: config.Provider{Name: "cloud", Auth: "/fake/dir/auth.json"}}

		err := InitToken(mod, cl)
		require.Error(t, err)
	})

	t.Run("Should fail when using valid json file but with mismatch structure", func(t *testing.T) {
		// START PREPARE FAKE JSON FILE
		fakeJsonFile := `
{
	"key": "value"
}
`
		fakeJsonPath := "/tmp/fake-auth.json"
		f, err := os.OpenFile(fakeJsonPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		require.NoError(t, err)
		defer f.Close()
		json.NewEncoder(f).Encode(fakeJsonFile)
		// DONE PREPARE FAKE JSON FILE

		mod := &config.Model{Provider: config.Provider{Name: "cloud", Auth: fakeJsonPath}}

		err = InitToken(mod, cl)
		require.Error(t, err)

		os.Remove(fakeJsonPath)
	})

	t.Run("Should fail when using invalid or fake server host", func(t *testing.T) {
		// START PREPARE FAKE JSON FILE
		fakeJsonFile := token.NewToken{
			RefreshToken: "token",
			ClientID:     "client",
			ClientSecret: "secret",
			TokenUrl:     "http://localhost",
		}
		fakeJsonPath := "/tmp/fake-auth.json"
		f, err := os.OpenFile(fakeJsonPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		require.NoError(t, err)
		defer f.Close()
		json.NewEncoder(f).Encode(fakeJsonFile)
		// DONE PREPARE FAKE JSON FILE

		mod := &config.Model{Provider: config.Provider{Name: "cloud", Auth: fakeJsonPath}}

		err = InitToken(mod, cl)
		require.Error(t, err)

		os.Remove(fakeJsonPath)
	})

	t.Run("Should fail when using not exist file path to save new token file", func(t *testing.T) {
		// START PREPARE FAKE JSON FILE
		fakeJsonFile := token.NewToken{
			RefreshToken: "token",
			ClientID:     "client",
			ClientSecret: "secret",
			TokenUrl:     server.URL,
		}
		fakeJsonPath := "/tmp/fake-auth.json"
		f, err := os.OpenFile(fakeJsonPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		require.NoError(t, err)
		defer f.Close()
		json.NewEncoder(f).Encode(fakeJsonFile)
		// DONE PREPARE FAKE JSON FILE

		mod := &config.Model{Provider: config.Provider{
			Name: "cloud", Auth: fakeJsonPath,
			Token: "/fake/path/token.json",
		}}

		err = InitToken(mod, server.Client())
		require.Error(t, err)

		os.Remove(fakeJsonPath)
	})

	t.Run("Should pass with valid values for server host, auth file & token file path", func(t *testing.T) {
		// START PREPARE FAKE JSON FILE
		fakeJsonFile := token.NewToken{
			RefreshToken: "token",
			ClientID:     "client",
			ClientSecret: "secret",
			TokenUrl:     server.URL,
		}
		fakeJsonPath := "/tmp/fake-auth.json"
		f, err := os.OpenFile(fakeJsonPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		require.NoError(t, err)
		defer f.Close()
		json.NewEncoder(f).Encode(fakeJsonFile)
		// DONE PREPARE FAKE JSON FILE

		mod := &config.Model{Provider: config.Provider{
			Name: "cloud", Auth: fakeJsonPath,
			Token: "/tmp/fake-token.json",
		}}

		err = InitToken(mod, server.Client())
		require.NoError(t, err)

		os.Remove(fakeJsonPath)
		os.Remove("/tmp/fake-token.json")
	})

	t.Cleanup(func() {
		server.Close()
	})
}
