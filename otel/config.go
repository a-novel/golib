package otel

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var AppName string

type Config interface {
	Init() error
	GetPropagators() (propagation.TextMapPropagator, error)
	GetTraceProvider() (trace.TracerProvider, error)
	GetLogger() (log.LoggerProvider, error)
	Flush()
}

func SetAppName(name string) {
	AppName = name
}

func InitOtel(config Config) error {
	err := config.Init()
	if err != nil {
		return fmt.Errorf("initialize otel: %w", err)
	}

	tracePropagator, err := config.GetPropagators()
	if err != nil {
		return fmt.Errorf("get trace propagators: %w", err)
	}

	traceProvider, err := config.GetTraceProvider()
	if err != nil {
		return fmt.Errorf("get trace provider: %w", err)
	}

	logger, err := config.GetLogger()
	if err != nil {
		return fmt.Errorf("get logger: %w", err)
	}

	otel.SetTextMapPropagator(tracePropagator)
	otel.SetTracerProvider(traceProvider)
	global.SetLoggerProvider(logger)

	return nil
}
