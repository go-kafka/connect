package connect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// StatusUnprocessableEntity is the status code returned when sending a
	// request with invalid fields.
	StatusUnprocessableEntity = 422
)

const (
	// DefaultHostURL is the default HTTP host used for connecting to a Kafka
	// Connect REST API.
	DefaultHostURL = "http://localhost:8083/"
	userAgent      = "go-kafka/0.9 connect/" + VERSION
)

// A Client manages communication with the Kafka Connect REST API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests. Defaults to http://localhost:8083/. BaseURL
	// should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Kafka Connect API.
	UserAgent string
}

// NewClient returns a new Kafka Connect API client. If a nil httpClient is
// provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	baseURL, _ := url.Parse(DefaultHostURL)

	client := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	return client
}

// NewRequest creates an API request. A relative URL can be provided in path,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON-encoded and included as the
// request body.
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	url := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		request.Header.Set("User-Agent", c.UserAgent)
	}

	return request, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON-decoded and stored in the value pointed to by v, or returned as an
// error if an API or HTTP error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		return response, buildError(req, response)
	}

	if v != nil {
		err = json.NewDecoder(response.Body).Decode(v)
		if err == io.EOF {
			err = nil // ignore EOF, empty response body
		}
	}

	return response, err
}

func buildError(req *http.Request, resp *http.Response) error {
	apiError := &APIError{Response: resp}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, apiError)
	}

	// Possibly a general HTTP error, e.g. we're not even talking to a valid
	// Kafka Connect API host
	if apiError.Code == 0 {
		return fmt.Errorf("HTTP %v on %v %v", resp.Status, req.Method, req.URL)
	}
	return apiError
}
