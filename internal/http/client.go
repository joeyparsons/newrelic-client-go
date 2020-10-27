package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	retryablehttp "github.com/hashicorp/go-retryablehttp"

	"github.com/joeyparsons/newrelic-client-go/internal/version"
	"github.com/joeyparsons/newrelic-client-go/pkg/config"
	nrErrors "github.com/joeyparsons/newrelic-client-go/pkg/errors"
)

const (
	defaultNewRelicRequestingServiceHeader = "NewRelic-Requesting-Services"
	defaultServiceName                     = "newrelic-client-go"
	defaultTimeout                         = time.Second * 30
	defaultRetryMax                        = 3
)

var (
	defaultUserAgent = fmt.Sprintf("newrelic/%s/%s (https://github.com/newrelic/%s)", defaultServiceName, version.Version, defaultServiceName)
)

// Client represents a client for communicating with the New Relic APIs.
type Client struct {
	// client represents the underlying HTTP client.
	client *retryablehttp.Client

	// config is the HTTP client configuration.
	config config.Config

	// authStrategy allows us to use multiple authentication methods for API calls
	authStrategy RequestAuthorizer

	// compressor is used to compress the body of a request, and set the content-encoding header
	compressor RequestCompressor

	errorValue ErrorResponse
}

// NewClient is used to create a new instance of Client.
func NewClient(cfg config.Config) Client {
	c := http.Client{
		Timeout: defaultTimeout,
	}

	if cfg.Timeout != nil {
		c.Timeout = *cfg.Timeout
	}

	if cfg.HTTPTransport != nil {
		if transport, ok := (cfg.HTTPTransport).(*http.Transport); ok {
			c.Transport = transport
		}
	} else {
		c.Transport = http.DefaultTransport
	}

	if cfg.UserAgent == "" {
		cfg.UserAgent = defaultUserAgent
	}

	// Either set or append the library name
	if cfg.ServiceName == "" {
		cfg.ServiceName = defaultServiceName
	} else {
		cfg.ServiceName = fmt.Sprintf("%s|%s", cfg.ServiceName, defaultServiceName)
	}

	r := retryablehttp.NewClient()
	r.HTTPClient = &c
	r.RetryMax = defaultRetryMax
	r.CheckRetry = RetryPolicy

	// Disable logging in go-retryablehttp since we are logging requests directly here
	r.Logger = nil

	client := Client{
		authStrategy: &ClassicV2Authorizer{},
		client:       r,
		config:       cfg,
		errorValue:   &DefaultErrorResponse{},
	}

	switch cfg.Compression {
	case config.Compression.Gzip:
		client.compressor = &GzipCompressor{}
	default:
		client.compressor = &NoneCompressor{}
	}

	return client
}

// SetAuthStrategy is used to set the default auth strategy for this client
// which can be overridden per request
func (c *Client) SetAuthStrategy(da RequestAuthorizer) {
	c.authStrategy = da
}

// SetRequestCompressor is used to enable compression on the request using
// the RequestCompressor specified
func (c *Client) SetRequestCompressor(compressor RequestCompressor) {
	c.compressor = compressor
}

// SetErrorValue is used to unmarshal error body responses in JSON format.
func (c *Client) SetErrorValue(v ErrorResponse) *Client {
	c.errorValue = v
	return c
}

// Get represents an HTTP GET request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *Client) Get(
	url string,
	queryParams interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodGet, url, queryParams, nil, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Post represents an HTTP POST request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The reqBody argument will be marshaled to JSON from the type provided and included in the request body.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *Client) Post(
	url string,
	queryParams interface{},
	reqBody interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPost, url, queryParams, reqBody, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Put represents an HTTP PUT request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The reqBody argument will be marshaled to JSON from the type provided and included in the request body.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *Client) Put(
	url string,
	queryParams interface{},
	reqBody interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPut, url, queryParams, reqBody, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Delete represents an HTTP DELETE request to a New Relic API.
// The queryParams argument can be used to add query string parameters to the requested URL.
// The respBody argument will be unmarshaled from JSON in the response body to the type provided.
// If respBody is not nil and the response body cannot be unmarshaled to the type provided, an error will be returned.
func (c *Client) Delete(url string,
	queryParams interface{},
	respBody interface{},
) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodDelete, url, queryParams, nil, respBody)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// logNice removes newlines, tabs, and \" from the body of a nerdgraph request.
// This allows for easier debugging and testing the content straight from the
// log file.
func logNice(body string) string {
	var newBody string
	newBody = strings.ReplaceAll(body, "\n", " ")
	newBody = strings.ReplaceAll(newBody, "\t", " ")
	newBody = strings.ReplaceAll(newBody, "\\\"", `"`)
	re := regexp.MustCompile(` +`)
	newBody = re.ReplaceAllString(newBody, " ")

	return newBody
}

// obfuscate receives a string, and replaces everything after the first 8
// characters with an asterisk before returning the result.
func obfuscate(input string) string {
	result := make([]string, len(input))
	parts := strings.Split(input, "")

	for i, x := range parts {
		if i < 8 {
			result[i] = x
		} else {
			result[i] = "*"
		}
	}

	return strings.Join(result, "")
}

func logCleanHeaderMarshalJSON(header http.Header) ([]byte, error) {
	h := http.Header{}

	for k, values := range header {
		if _, ok := h[k]; ok {
			h[k] = make([]string, len(values))
		}

		switch k {
		case "Api-Key", "X-Api-Key", "X-Insert-Key":
			newValues := []string{}
			for _, v := range values {
				newValues = append(newValues, obfuscate(v))
			}

			if len(newValues) > 0 {
				h[k] = newValues
			} else {
				h[k] = values
			}
		default:
			h[k] = values
		}
	}

	return json.Marshal(h)
}

// Do initiates an HTTP request as configured by the passed Request struct.
func (c *Client) Do(req *Request) (*http.Response, error) {
	var resp *http.Response
	var errorValue ErrorResponse
	var body []byte

	logger := c.config.GetLogger()
	logger.Debug("performing request", "method", req.method, "url", req.url)

	for i := 0; ; i++ {
		var shouldRetry bool
		var err error

		errorValue = req.errorValue.New()
		resp, body, shouldRetry, err = c.innerDo(req, errorValue, i)
		if err != nil {
			return nil, err
		}

		if !shouldRetry {
			break
		}
	}

	if !isResponseSuccess(resp) {
		return nil, nrErrors.NewUnexpectedStatusCode(resp.StatusCode, errorValue.Error())
	}

	if errorValue.IsNotFound() {
		return nil, nrErrors.NewNotFound("resource not found")
	}

	if errorValue.IsTimeout() {
		return nil, nrErrors.NewTimeout("server timeout")
	}

	if errorValue.Error() != "" {
		return nil, errors.New(errorValue.Error())
	}

	if req.value == nil {
		return resp, nil
	}

	jsonErr := json.Unmarshal(body, req.value)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return resp, nil
}

func (c *Client) innerDo(req *Request, errorValue ErrorResponse, i int) (*http.Response, []byte, bool, error) {
	shouldRetry := false

	r, err := req.makeRequest()
	if err != nil {
		return nil, nil, false, err
	}

	logger := c.config.GetLogger()

	logHeaders, err := logCleanHeaderMarshalJSON(r.Header)
	if err != nil {
		return nil, nil, false, err
	}

	if req.reqBody != nil {
		switch reflect.TypeOf(req.reqBody).String() {
		case "*http.graphQLRequest":
			x := req.reqBody.(*graphQLRequest)

			logVariables, marshalErr := json.Marshal(x.Variables)
			if marshalErr != nil {
				return nil, nil, false, marshalErr
			}

			logger.Trace("request details",
				"headers", logNice(string(logHeaders)),
				"query", logNice(x.Query),
				"variables", string(logVariables),
			)
		case "string":
			logger.Trace("request details", "headers", string(logHeaders), "body", logNice(req.reqBody.(string)))
		}
	} else {
		logger.Trace("request details", "headers", string(logHeaders))
	}

	if i > 0 {
		logger.Debug(fmt.Sprintf("retrying request (attempt %d)", i), "method", req.method, "url", r.URL)
	}

	resp, retryErr := c.client.Do(r)
	if retryErr != nil {
		return resp, nil, false, retryErr
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return resp, nil, false, &nrErrors.NotFound{}
	}

	body, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		return resp, body, false, readErr
	}

	logHeaders, err = json.Marshal(resp.Header)
	if err != nil {
		return resp, body, false, err
	}

	logger.Trace("request completed", "method", req.method, "url", r.URL, "status_code", resp.StatusCode, "headers", string(logHeaders), "body", string(body))

	_ = json.Unmarshal(body, &errorValue)

	if errorValue.IsTimeout() {
		shouldRetry = true
	}

	if !shouldRetry {
		return resp, body, false, nil
	}

	remain := c.client.RetryMax - i
	if remain <= 0 {
		logger.Debug(fmt.Sprintf("giving up after %d attempts", c.client.RetryMax), "method", req.method, "url", r.URL)
		return resp, body, false, nil
	}

	wait := c.client.Backoff(c.client.RetryWaitMin, c.client.RetryWaitMax, i, resp)

	time.Sleep(wait)

	return resp, body, true, nil
}

// Ensures the response status code falls within the
// status codes that are commonly considered successful.
func isResponseSuccess(resp *http.Response) bool {
	statusCode := resp.StatusCode

	return statusCode >= http.StatusOK && statusCode <= 299
}

// NerdGraphQuery runs a Nerdgraph query.
func (c *Client) NerdGraphQuery(query string, vars map[string]interface{}, respBody interface{}) error {
	req, err := c.NewNerdGraphRequest(query, vars, respBody)
	if err != nil {
		return err
	}

	_, err = c.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// NewNerdGraphRequest runs a Nerdgraph request object.
func (c *Client) NewNerdGraphRequest(query string, vars map[string]interface{}, respBody interface{}) (*Request, error) {
	graphqlReqBody := &graphQLRequest{
		Query:     query,
		Variables: vars,
	}

	graphqlRespBody := &graphQLResponse{
		Data: respBody,
	}

	req, err := c.NewRequest(http.MethodPost, c.config.Region().NerdGraphURL(), nil, graphqlReqBody, graphqlRespBody)
	if err != nil {
		return nil, err
	}

	req.SetAuthStrategy(&NerdGraphAuthorizer{})
	req.SetErrorValue(&GraphQLErrorResponse{})

	return req, nil
}
