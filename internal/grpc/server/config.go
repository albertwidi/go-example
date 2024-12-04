package server

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	tracenoop "go.opentelemetry.io/otel/trace/noop"
)

type Config struct {
	Address      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Trace        *TraceConfig
	Meter        *MeterConfig
}

func (c Config) Validate() error {
	if c.Trace == nil {
		c.Trace = &TraceConfig{
			// Use a noop trace provider as the default trace configuration.
			Tracer: tracenoop.NewTracerProvider().Tracer("noop"),
		}
	}
	if err := c.Trace.Validate(); err != nil {
		return err
	}
	return nil
}

type TraceConfig struct {
	Tracer            trace.Tracer
	DefaultAttributes []attribute.KeyValue
}

func (t *TraceConfig) Validate() error {
	return nil
}

type MeterConfig struct {
}
