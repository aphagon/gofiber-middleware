package otelfiber

import (
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// ContextKey
	// Optional. Default: "request".
	ServiceName string

	// ContextKey
	// Optional. Default: "otel-fiber".
	ContextKey string

	// TracerName
	// Optional. Default: github.com/aphagon/gofiber-middleware/otelfiber
	TracerName string

	// TracerProvider
	// Optional. Default: otel.GetTracerProvider().
	TracerProvider trace.TracerProvider

	// Propagators
	// Optional. Default: otel.GetTextMapPropagator().
	Propagators propagation.TextMapPropagator

	// TracerStartAttributes
	//
	// Optional. Default: []trace.SpanOption{
	// 	trace.WithSpanKind(trace.SpanKindServer),
	// 	trace.WithNewRoot(),
	// 	trace.WithRecord(),
	// }
	TracerStartAttributes []trace.SpanOption
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	ServiceName:    "request",
	ContextKey:     "otel-fiber",
	TracerName:     "github.com/aphagon/gofiber-middleware/otelfiber",
	TracerProvider: otel.GetTracerProvider(),
	Propagators:    otel.GetTextMapPropagator(),
	TracerStartAttributes: []trace.SpanOption{
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithNewRoot(),
		trace.WithRecord(),
	},
}

// helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Next == nil {
		cfg.Next = ConfigDefault.Next
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = ConfigDefault.ServiceName
	}

	if cfg.ContextKey == "" {
		cfg.ContextKey = ConfigDefault.ContextKey
	}

	if cfg.TracerName == "" {
		cfg.TracerName = ConfigDefault.TracerName
	}

	if cfg.TracerProvider == nil {
		cfg.TracerProvider = ConfigDefault.TracerProvider
	}

	if cfg.Propagators == nil {
		cfg.Propagators = ConfigDefault.Propagators
	}

	if len(cfg.TracerStartAttributes) == 0 {
		cfg.TracerStartAttributes = ConfigDefault.TracerStartAttributes
	}

	return cfg
}
