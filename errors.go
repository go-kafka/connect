package connect

// APIError holds information returned from a Kafka Connect API instance about
// why an API call failed.
type APIError struct {
	Code    int    `json:"error_code"`
	Message string `json:"message"`
}
