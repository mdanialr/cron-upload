package token

import (
	"fmt"
	"strings"

	"google.golang.org/api/drive/v3"
)

// CheckTokenValidity just testing using random request to make sure connection to google apis is successful and
// make sure the token is still valid. Only return error if the error contain authError.
func CheckTokenValidity(srv *drive.Service) error {
	_, err := srv.About.Get().Do()
	if err != nil {
		if strings.Contains(err.Error(), "authError") {
			return fmt.Errorf("invalid token: %s\n", strings.TrimSpace(err.Error()))
		}
	}
	return nil
}
