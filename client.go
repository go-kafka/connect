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
	userAgent      = "go-kafka/0.9 connect/" + Version
)

// A Client manages communication with the Kafka Connect REST API.
type Client struct {
	host *url.URL // Base host URL for API requests.

	// HTTP client used to communicate with the API. By default
	// http.DefaultClient will be used.
	HTTPClient *http.Client

	// User agent used when communicating with the Kafka Connect API.
	UserAgent string
}

// NewClient returns a new Kafka Connect API client that communicates with the
// optional host. If no host is given, DefaultHostURL (localhost) is used.
func NewClient(host ...string) *Client {
	var hostURL *url.URL
	var err error

	switch len(host) {
	case 0:
		hostURL, _ = url.Parse(DefaultHostURL)
	case 1:
		hostURL, err = url.Parse(host[0])
		if err != nil {
			panic(err.Error())
		}
	default:
		panic("only one host URL can be given")
	}

	return &Client{host: hostURL, UserAgent: userAgent}
}

func (c *Client) httpClient() *http.Client {
	if c.HTTPClient == nil {
		return http.DefaultClient
	}
	return c.HTTPClient
}

// Host returns the API root URL the Client is configured to talk to.
func (c *Client) Host() string {
	return c.host.String()
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

	url := c.host.ResolveReference(rel)

	var contentType string
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
		contentType = "application/json"
	}

	request, err := http.NewRequest(method, url.String(), buf)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if c.UserAgent != "" {
		request.Header.Set("User-Agent", c.UserAgent)
	}

	return request, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON-decoded and stored in the value pointed to by v, or returned as an
// error if an API or HTTP error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	response, err := c.httpClient().Do(req)
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

// Simple GET helper with no request body.
func (c *Client) get(path string, v interface{}) (*http.Response, error) {
	return c.doRequest("GET", path, nil, v)
}

func (c *Client) delete(path string) (*http.Response, error) {
	return c.doRequest("DELETE", path, nil, nil)
}

func (c *Client) doRequest(method, path string, body, v interface{}) (*http.Response, error) {
	request, err := c.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	return c.Do(request, v)
}

func buildError(req *http.Request, resp *http.Response) error {
	apiError := APIError{Response: resp}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		_ = json.Unmarshal(data, &apiError) // Fall back on general error below
	}

	// Possibly a general HTTP error, e.g. we're not even talking to a valid
	// Kafka Connect API host
	if apiError.Code == 0 {
		return fmt.Errorf("HTTP %v on %v %v", resp.Status, req.Method, req.URL)
	}
	return apiError
}
