package httpclient

import (
	"bytes"
	"context"
	"net/http"
)

// Client interface for transport
type Client interface {
	Send(addHdr map[string]string, body []byte, method string, url string) (*http.Response, error)
	SendWithContext(ctx context.Context, addHdr map[string]string, body []byte, method string, url string) (*http.Response, error)
}

// HTTPClient instance
type HTTPClient struct {
	Client *http.Client
}

// NewHTTPClient creates http client
func NewHTTPClient(cl *http.Client) *HTTPClient {
	return &HTTPClient{
		Client: cl,
	}
}

// Send an HTTP request using the Header, body and URL provided
func (cl *HTTPClient) Send(
	addHdr map[string]string,
	body []byte,
	method string,
	url string) (*http.Response, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	for key, value := range addHdr {
		req.Header.Set(key, value)
	}

	resp, err := cl.Client.Do(req)

	return resp, err
}

// SendWithContext send an HTTP request using the Header, body and URL provided
func (cl *HTTPClient) SendWithContext(ctx context.Context,
	addHdr map[string]string,
	body []byte,
	method string,
	url string) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	for key, value := range addHdr {
		req.Header.Set(key, value)
	}

	resp, err := cl.Client.Do(req)
	return resp, err
}
