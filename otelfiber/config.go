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

	// LocalKeyName
	// Optional. Default: "otel-fiber".
	LocalKeyName string

	// SpanName is a template for span naming.
	// The scope is fiber context.
	SpanName string

	// TracerProvider
	// Optional. Default: otel.GetTracerProvider().
	TracerProvider trace.TracerProvider

	// Propagators
	// Optional. Default: otel.GetTextMapPropagator().
	Propagators propagation.TextMapPropagator

	// TracerStartAttributes
	//
	// Optional. Default: []trace.SpanStartOption{
	// 	trace.WithSpanKind(trace.SpanKindServer),
	// 	trace.WithNewRoot(),
	// }
	TracerStartAttributes []trace.SpanStartOption
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	SpanName:       "http/request",
	LocalKeyName:   "otel-fiber",
	TracerProvider: otel.GetTracerProvider(),
	Propagators:    otel.GetTextMapPropagator(),
	TracerStartAttributes: []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithNewRoot(),
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

	if cfg.SpanName == "" {
		cfg.SpanName = ConfigDefault.SpanName
	}

	if cfg.LocalKeyName == "" {
		cfg.LocalKeyName = ConfigDefault.LocalKeyName
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
