package opa

import (
	"context"
	"io"
	"net/http"
)

type Transport interface {
	Post(ctx context.Context, url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type HttpTransport struct{}

func NewHttpTransport() Transport {
	return &HttpTransport{}
}

func (client HttpTransport) Post(ctx context.Context, url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return http.DefaultClient.Do(req)
}
