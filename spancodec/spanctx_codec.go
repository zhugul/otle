
package spancodec

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const otelEnvKey string = "OTEL_TRACE_SPAN_CONTEXT"


func EncodeSpanContextToEnv(ctx context.Context) (string, string) {
	//span := trace.SpanFromContext(ctx)
	//spanContext := span.SpanContext()
	//fmt.Println("Trace ID:" + spanContext.TraceID())
	//fmt.Println("Span ID:" + spanContext.SpanID())

	// Encode SpanContext into key-value text format
	textMapPropagator := propagation.TraceContext{}
	headers := propagation.MapCarrier{}
	textMapPropagator.Inject(ctx, headers)

	var parts []string
	for k, v := range headers {
		parts = append(parts, k+"="+v)
	}

	encodedContext := strings.Join(parts, ",")

	// Output the encoded SpanContext
	//key OTEL_TRACE_SPAN_CONTEXT, value is encodedContext
	key: = otelEnvKey
	fmt.Println(encodedContext)
	return key, encodedContext
}



func DecodeSpanContextFromEnv() context.Context {
	key = otelEnvKey
	encodedContext, ok := os.LookupEnv(key)
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
	ctx := propagator.Extract(context.Background(), carrier)
	return ctx
}
