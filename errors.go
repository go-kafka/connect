package connect

import (
	"fmt"
	"net/http"
)

// APIError holds information returned from a Kafka Connect API instance about
// why an API call failed.
type APIError struct {
	Code     int            `json:"error_code"`
	Message  string         `json:"message"`
	Response *http.Response // HTTP response that caused this error
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v (HTTP %d)", e.Message, e.Code)
}
