package otelfiber

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	tracer := cfg.TracerProvider.Tracer(
		"github.com/aphagon/gofiber-middleware",
		trace.WithInstrumentationVersion(contrib.SemVersion()),
	)

	spanTmpl := template.Must(template.New("span").Parse(cfg.SpanName))

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// concat all span options, dynamic and static
		spanOptions := concatSpanOptions(
			[]trace.SpanStartOption{
				trace.WithAttributes(semconv.HTTPMethodKey.String(c.Method())),
				trace.WithAttributes(semconv.HTTPTargetKey.String(string(c.Request().RequestURI()))),
				trace.WithAttributes(semconv.HTTPRouteKey.String(c.Path())),
				trace.WithAttributes(semconv.HTTPURLKey.String(c.OriginalURL())),
				trace.WithAttributes(semconv.HTTPUserAgentKey.String(string(c.Request().Header.UserAgent()))),
				trace.WithAttributes(semconv.HTTPRequestContentLengthKey.Int(c.Request().Header.ContentLength())),
				trace.WithAttributes(semconv.HTTPSchemeKey.String(c.Protocol())),
				trace.WithAttributes(semconv.HTTPClientIPKey.String(c.Context().RemoteAddr().String())),
				trace.WithAttributes(semconv.HTTPHostKey.String(c.Hostname())),
				trace.WithAttributes(semconv.NetHostIPKey.String(c.IP())),
				trace.WithAttributes(semconv.NetTransportTCP),
				trace.WithSpanKind(trace.SpanKindServer),
			},
			cfg.TracerStartAttributes,
		)

		spanName := new(bytes.Buffer)
		err := spanTmpl.Execute(spanName, c)
		if err != nil {
			return fmt.Errorf("cannot execute span name template: %w", err)
		}

		ctx, span := tracer.Start(
			c.Context(),
			spanName.String(),
			spanOptions...,
		)

		c.Locals(cfg.LocalKeyName, ctx)
		defer span.End()

		err = c.Next()
		if err != nil {
			span.SetAttributes(attribute.String("fiber.error", err.Error()))
		}

		statusCode := c.Response().StatusCode()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(statusCode)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(statusCode)
		span.SetAttributes(attrs...)
		span.SetAttributes(semconv.HTTPResponseContentLengthKey.Int(c.Response().Header.ContentLength()))
		span.SetStatus(spanStatus, spanMessage)

		return err
	}
}

func concatSpanOptions(sources ...[]trace.SpanStartOption) []trace.SpanStartOption {
	var spanOptions []trace.SpanStartOption
	for _, source := range sources {
		spanOptions = append(spanOptions, source...)
	}
	return spanOptions
}
