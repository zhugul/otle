
package sctxcodec

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)



func encodeSpanContextToEnv(ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()
	fmt.Println(spanContext.TraceID())
	fmt.Println(spanContext.SpanID())

	// Encode SpanContext into key-value text format
	textMapPropagator := propagation.TraceContext{}
	headers := propagation.MapCarrier{}
	textMapPropagator.Inject(ctx, headers)

	var parts []string
	for k, v := range headers {
		parts = append(parts, k+"="+v)
	}

	encodedContext := strings.Join(parts, ",")
	fmt.Println("encodedContext" + encodedContext)

	// Put the encoded SpanContext into environment variables
	os.Setenv("OTEL_TRACE_SPAN_CONTEXT", encodedContext)
}



func decodeSpanContextFromEnv() context.Context {
	encodedContext, ok := os.LookupEnv("OTEL_TRACE_SPAN_CONTEXT")
	if !ok {
		fmt.Println("Error:Not found OTEL_TRACE_SPAN_CONTEXT.")
		return nil
		//encodedContext = "traceparent=00-fa984ea3aed3d4055c9f0217d5985867-25928de94ca993a1-01"
	}
	propagator := propagation.TraceContext{}
	headers := strings.Split(encodedContext, ",")
	var headerMap map[string]string
	for _, header := range headers {
		parts := strings.SplitN(header, "=", 2)
		if len(parts) == 2 {
			if headerMap == nil {
				headerMap = make(map[string]string)
			}
			headerMap[parts[0]] = parts[1]
		}
	}
	carrier := propagation.MapCarrier(headerMap)
	fmt.Println(carrier)
	ctx := propagator.Extract(context.Background(), carrier)
	return ctx
}
