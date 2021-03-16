package otelfiber

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/semconv"
	"go.opentelemetry.io/otel/trace"
)

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	tracer := cfg.TracerProvider.Tracer(
		cfg.TracerName,
		trace.WithInstrumentationVersion(contrib.SemVersion()),
	)

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// concat all span options, dynamic and static
		opts := concatSpanOptions(
			[]trace.SpanOption{
				trace.WithAttributes(semconv.HTTPMethodKey.String(c.Route().Method)),
				trace.WithAttributes(semconv.HTTPTargetKey.String(string(c.Request().RequestURI()))),
				trace.WithAttributes(semconv.HTTPRouteKey.String(c.Route().Path)),
				trace.WithAttributes(semconv.HTTPURLKey.String(c.OriginalURL())),
				trace.WithAttributes(semconv.HTTPUserAgentKey.String(string(c.Request().Header.UserAgent()))),
				trace.WithAttributes(semconv.HTTPRequestContentLengthKey.Int(c.Request().Header.ContentLength())),
				trace.WithAttributes(semconv.HTTPSchemeKey.String(c.Protocol())),
				trace.WithAttributes(semconv.HTTPServerNameKey.String(cfg.ServiceName)),
				trace.WithAttributes(semconv.NetHostIPKey.String(c.IP())),
				trace.WithAttributes(semconv.NetTransportTCP),
				trace.WithSpanKind(trace.SpanKindServer),
			},
			cfg.TracerStartAttributes,
		)

		spanName := c.Route().Path
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Route().Method)
		}

		ctx, span := tracer.Start(c.Context(), spanName, opts...)
		c.Locals(cfg.ContextKey, ctx)
		defer span.End()

		err := c.Next()
		if err != nil {
			span.SetAttributes(attribute.String("fiber.error", err.Error()))
		}

		statusCode := c.Response().StatusCode()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(statusCode)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(statusCode)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)

		return err
	}
}

func concatSpanOptions(sources ...[]trace.SpanOption) []trace.SpanOption {
	var spanOptions []trace.SpanOption
	for _, source := range sources {
		spanOptions = append(spanOptions, source...)
	}
	return spanOptions
}
