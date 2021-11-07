package datadog

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"time"

	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type recordContextKey string

const (
	key = recordContextKey("tracerecord")
)

type clientTraceRecord struct {
	dnsStart             time.Time
	dnsDone              time.Time
	connectDone          time.Time
	gotConn              time.Time
	gotFirstResponseByte time.Time
	tlsHandshakeStart    time.Time
	tlsHandshakeDone     time.Time
	tlsConnectionState   tls.ConnectionState
}

func WithHTTPClientTrace(ctx context.Context) context.Context {
	record := &clientTraceRecord{}

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { record.dnsStart = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { record.dnsDone = time.Now() },
		ConnectStart: func(_, _ string) {
			if record.dnsDone.IsZero() { // connecting directly via IP (no DNS involved)
				record.dnsStart = time.Now()
				record.dnsDone = record.dnsStart
			}
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				return
			}
			record.connectDone = time.Now()
		},
		GotConn: func(_ httptrace.GotConnInfo) {
			record.gotConn = time.Now()
			if record.connectDone.IsZero() { // cached connection
				record.connectDone = record.gotConn
				record.dnsStart = record.gotConn
				record.dnsDone = record.gotConn
			}
		},
		GotFirstResponseByte: func() { record.gotFirstResponseByte = time.Now() },
		TLSHandshakeStart:    func() { record.tlsHandshakeStart = time.Now() },
		TLSHandshakeDone: func(state tls.ConnectionState, _ error) {
			record.tlsHandshakeDone = time.Now()
			record.tlsConnectionState = state
		},
	}

	// store the data in the context; the trace is used by net/http
	// and the record is pulled out in the roundtripper.
	ctx = httptrace.WithClientTrace(ctx, trace)
	ctx = context.WithValue(ctx, key, record)

	return ctx
}

func withClientTraceRecorder(wrapped *http.Client) *http.Client {
	if wrapped.Transport == nil {
		wrapped.Transport = http.DefaultTransport
	}

	// first wrap with http client tracing
	wrapped.Transport = &clientTraceRoundTripper{
		base: wrapped.Transport,
	}

	return wrapped
}

type clientTraceRoundTripper struct {
	base http.RoundTripper
}

func (rt *clientTraceRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {

	res, err = rt.base.RoundTrip(req)

	ctx := req.Context()

	record, ok := ctx.Value(key).(*clientTraceRecord)
	if !ok {
		return
	}

	span, ok := ddtracer.SpanFromContext(ctx)
	if !ok {
		return
	}

	// connection stage timings
	span.SetTag("httpclient.stage.dnslookup", (record.dnsDone.Sub(record.dnsStart)).Milliseconds())
	span.SetTag("httpclient.stage.tcpconnection", (record.connectDone.Sub(record.dnsDone)).Milliseconds())
	span.SetTag("httpclient.stage.tlshandshake", (record.tlsHandshakeDone.Sub(record.tlsHandshakeStart)).Milliseconds())
	span.SetTag("httpclient.stage.serverprocessing", (record.gotFirstResponseByte.Sub(record.gotConn)).Milliseconds())

	// composite measurements
	span.SetTag("httpclient.connect", (record.connectDone.Sub(record.dnsStart)).Milliseconds())
	span.SetTag("httpclient.pretransfer", (record.gotConn.Sub(record.dnsStart)).Milliseconds())
	span.SetTag("httpclient.starttransfer", (record.gotFirstResponseByte.Sub(record.dnsStart)).Milliseconds())

	// HTTP/TLS protocol deets
	span.SetTag("httpclient.tlsversion", tlsVersionString(record.tlsConnectionState.Version))
	span.SetTag("httpclient.ciphersuite", cipherSuiteString(record.tlsConnectionState.CipherSuite))
	span.SetTag("httpclient.protocol", record.tlsConnectionState.NegotiatedProtocol)

	return
}

func tlsVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS10"
	case tls.VersionTLS11:
		return "TLS11"
	case tls.VersionTLS12:
		return "TLS12"
	case tls.VersionTLS13:
		return "TLS13"
	case tls.VersionSSL30:
		return "SSL30"
	default:
		return fmt.Sprintf("Unknown TLS version: 0x%04x", version)
	}
}

func cipherSuiteString(suite uint16) string {
	switch suite {
	case tls.TLS_RSA_WITH_RC4_128_SHA:
		return "TLS_RSA_WITH_RC4_128_SHA"
	case tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:
		return "TLS_RSA_WITH_3DES_EDE_CBC_SHA"
	case tls.TLS_RSA_WITH_AES_128_CBC_SHA:
		return "TLS_RSA_WITH_AES_128_CBC_SHA"
	case tls.TLS_RSA_WITH_AES_256_CBC_SHA:
		return "TLS_RSA_WITH_AES_256_CBC_SHA"
	case tls.TLS_RSA_WITH_AES_128_CBC_SHA256:
		return "TLS_RSA_WITH_AES_128_CBC_SHA256"
	case tls.TLS_RSA_WITH_AES_128_GCM_SHA256:
		return "TLS_RSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_RSA_WITH_AES_256_GCM_SHA384:
		return "TLS_RSA_WITH_AES_256_GCM_SHA384"
	case tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA:
		return "TLS_ECDHE_ECDSA_WITH_RC4_128_SHA"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:
		return "TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA"
	case tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA:
		return "TLS_ECDHE_RSA_WITH_RC4_128_SHA"
	case tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA:
		return "TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:
		return "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA"
	case tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:
		return "TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:
		return "TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:
		return "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:
		return "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256"
	case tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:
		return "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
	case tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:
		return "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384"
	case tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256:
		return "TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256"
	case tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256:
		return "TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256"
	case tls.TLS_AES_128_GCM_SHA256:
		return "TLS_AES_128_GCM_SHA256"
	case tls.TLS_AES_256_GCM_SHA384:
		return "TLS_AES_256_GCM_SHA384"
	case tls.TLS_CHACHA20_POLY1305_SHA256:
		return "TLS_CHACHA20_POLY1305_SHA256"
	case tls.TLS_FALLBACK_SCSV:
		return "TLS_FALLBACK_SCSV"
	default:
		return fmt.Sprintf("Unknown cipher suite: 0x%04x", suite)
	}
}
