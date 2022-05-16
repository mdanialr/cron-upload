package token

import (
	"fmt"
	"strings"

	"google.golang.org/api/drive/v3"
)

// CheckRefreshValidity just testing using random request to make sure connection to google apis is successful and
// make sure the refresh token is still valid. Only return error if the error contain oauth2.
func CheckRefreshValidity(srv *drive.Service) error {
	_, err := srv.About.Get().Do()
	if err != nil {
		if strings.Contains(err.Error(), "oauth2") {
			return fmt.Errorf("expired refresh token: %s\n", strings.TrimSpace(err.Error()))
		}
	}
	return nil
}
