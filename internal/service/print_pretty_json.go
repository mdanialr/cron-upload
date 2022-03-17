package service

import (
	"bytes"
	"encoding/json"
)

// PrettyJson print prettify json response.
func PrettyJson(in []byte) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, in, "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
