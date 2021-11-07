package datadog

import (
	"context"
	"net/http"

	goahttp "goa.design/goa/v3/http"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// RequestDecoder wraps the request decoder and records a Datadog span
func RequestDecoder(decoderFunc func(*http.Request) goahttp.Decoder) func(*http.Request) goahttp.Decoder {
	return func(r *http.Request) goahttp.Decoder {
		return &tracingDecoder{
			req:     r,
			wrapped: decoderFunc(r),
		}
	}
}

type tracingDecoder struct {
	req     *http.Request
	wrapped goahttp.Decoder
}

func (dec *tracingDecoder) Decode(v interface{}) error {
	ctx := dec.req.Context()
	span, _ := ddtracer.StartSpanFromContext(ctx, "DecodeRequest")
	defer span.Finish()

	return dec.wrapped.Decode(v)
}

// ResponseEncoder wraps the request encoder and records a Datadog span
func ResponseEncoder(encoderFunc func(ctx context.Context, w http.ResponseWriter) goahttp.Encoder) func(ctx context.Context, w http.ResponseWriter) goahttp.Encoder {
	return func(ctx context.Context, w http.ResponseWriter) goahttp.Encoder {
		return &tracingEncoder{
			ctx:     ctx,
			wrapped: encoderFunc(ctx, w),
		}
	}
}

type tracingEncoder struct {
	ctx     context.Context
	wrapped goahttp.Encoder
}

func (dec *tracingEncoder) Encode(v interface{}) error {
	span, _ := ddtracer.StartSpanFromContext(dec.ctx, "EncodeRequest")
	defer span.Finish()

	return dec.wrapped.Encode(v)
}
