package authz

import (
	"context"
	"github.com/cultureamp/glamplify/env"
	"io"
	"net"
	"net/http"
	"time"
)

// Transport represents the mechanism to POST a request to an endpoint
type Transport interface {
	Post(ctx context.Context, url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

// HTTPConfig represents http config values
type HTTPConfig struct {
	ClientTimeout       time.Duration
	DialerTimeout       time.Duration
	TLSHandshakeTimeout time.Duration
}

// HTTPTransport contains the HTTPConfig and the network Client
type HTTPTransport struct {
	conf    *HTTPConfig
	network *http.Client
}

// NewHTTPTransport creates a new Transport
func NewHTTPTransport(configure ...func(*HTTPConfig)) Transport {
	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779

	c := env.GetInt(env.AuthzClientTimeoutEnv, 10000) // 10 secs
	clientTimeout := time.Duration(c) * time.Millisecond

	d := env.GetInt(env.AuthzDialerTimeoutEnv, 5000) // 5 secs
	dialerTimeout := time.Duration(d) * time.Millisecond

	t := env.GetInt(env.AuthzTLSTimeoutEnv, 5000) // 5 secs
	tlsTimeout := time.Duration(t) * time.Millisecond

	conf := &HTTPConfig{
		ClientTimeout:       clientTimeout,
		DialerTimeout:       dialerTimeout,
		TLSHandshakeTimeout: tlsTimeout,
	}

	for _, config := range configure {
		config(conf)
	}

	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: conf.DialerTimeout,
		}).DialContext,
		TLSHandshakeTimeout: conf.TLSHandshakeTimeout,
	}
	var netClient = &http.Client{
		Timeout:   conf.ClientTimeout,
		Transport: netTransport,
	}

	return &HTTPTransport{
		conf:    conf,
		network: netClient,
	}
}

// Post a request to the endpoint
func (client HTTPTransport) Post(ctx context.Context, url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	return client.network.Do(req)
}
