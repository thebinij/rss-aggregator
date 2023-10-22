package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts as API Kye from
// the headers of an HTTP request
// Ex: Authorization: ApiKey {insert apikey here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 || vals[0] != "ApiKey" {
		return "", errors.New("malformed auth header")
	}

	return vals[1], nil
}
