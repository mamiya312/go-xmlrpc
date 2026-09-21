package xmlrpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// Client calls an XML-RPC endpoint over HTTP.
type Client struct {
	endpoint   string
	httpClient *http.Client
}

// NewClient creates a client that uses http.DefaultClient.
func NewClient(endpoint string) *Client {
	return NewClientWithHTTPClient(endpoint, http.DefaultClient)
}

// NewClientWithHTTPClient creates a client with a custom HTTP client.
func NewClientWithHTTPClient(endpoint string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

// Call invokes method with args and decodes the response parameters into replies.
// Each reply must be a non-nil pointer to a value of the expected XML-RPC type.
func (c *Client) Call(ctx context.Context, method string, args []any, replies ...any) error {
	if c == nil || c.httpClient == nil || c.endpoint == "" {
		return fmt.Errorf("xmlrpc: client is not configured")
	}
	if ctx == nil {
		return fmt.Errorf("xmlrpc: nil context")
	}

	requestXML, err := ToRequestXML(method, args...)
	if err != nil {
		return fmt.Errorf("xmlrpc: encode request: %w", err)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint,
		bytes.NewReader([]byte(requestXML)),
	)
	if err != nil {
		return fmt.Errorf("xmlrpc: create request: %w", err)
	}
	request.Header.Set("Content-Type", "text/xml; charset=utf-8")
	request.Header.Set("Accept", "text/xml")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("xmlrpc: call %q: %w", method, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("xmlrpc: read response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("xmlrpc: call %q returned HTTP status %s: %s",
			method, response.Status, body)
	}

	reader, err := NewResponseReader(bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("xmlrpc: decode response: %w", err)
	}
	if err := reader.Decode(replies...); err != nil {
		return fmt.Errorf("xmlrpc: decode response values: %w", err)
	}
	return nil
}
