package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"

	"{{ .GoModule}}/internal/server/middleware"
)

// Config represents the server configuration
type Config struct {
	ServiceName       string
	Port              string
	PrometheusPort    string
	JaegerAgent       string
}

// PrometheusEnabled check configuration
func (c Config) PrometheusEnabled() bool {
	return c.PrometheusPort != ""
}

// TracingEnabled check configuration
func (c Config) TracingEnabled() bool {
	return c.JaegerAgent != ""
}

// Validate the config
func (c Config) Validate() error {
	return nil
}

func (c Config) grpcOpts() []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0)

	interceptors := make([]grpc.UnaryServerInterceptor, 0)
	interceptors = append(interceptors, middleware.Logging)
	interceptors = append(interceptors, grpc_recovery.UnaryServerInterceptor())
	if c.PrometheusEnabled() {
		interceptors = append(interceptors, grpc_prometheus.UnaryServerInterceptor)
	}
	if c.TracingEnabled() {
		interceptors = append(interceptors, middleware.Tracing(c.ServiceName, c.JaegerAgent))
	}
	interceptors = append(interceptors, errorMapper)
	opts = append(opts, grpc_middleware.WithUnaryServerChain(interceptors...))

	return opts
}