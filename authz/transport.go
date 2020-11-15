package authz

import (
	"context"
	"github.com/cultureamp/glamplify/helper"
	"io"
	"net"
	"net/http"
	"time"
)

type Transport interface {
	Post(ctx context.Context, url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

type HttpConfig struct {
	ClientTimeout       time.Duration
	DialerTimeout       time.Duration
	TLSHandshakeTimeout time.Duration
}

type HttpTransport struct {
	conf    *HttpConfig
	network *http.Client
}

func NewHttpTransport(configure ...func(*HttpConfig)) Transport {
	// https://medium.com/@nate510/don-t-use-go-s-default-http-client-4804cb19f779

	c := helper.GetEnvInt(ClientTimeoutEnv, 10000) // 10 secs
	clientTimeout := time.Duration(c) * time.Millisecond

	d := helper.GetEnvInt(DialerTimeoutEnv, 5000) // 5 secs
	dialerTimeout := time.Duration(d) * time.Millisecond

	t := helper.GetEnvInt(TLSTimeoutEnv, 5000) // 5 secs
	tlsTimeout := time.Duration(t) * time.Millisecond

	conf := &HttpConfig{
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

	return &HttpTransport{
		conf:    conf,
		network: netClient,
	}
}

func (client HttpTransport) Post(ctx context.Context, url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	return client.network.Do(req)
}
