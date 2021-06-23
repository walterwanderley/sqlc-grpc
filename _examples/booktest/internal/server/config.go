package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Config represents the server configuration
type Config struct {
	ServiceName     string
	Port            int
	PrometheusPort  int
	JaegerCollector string
	Cert            string
	Key             string
	EnableCors      bool
	EnableGrpcUI    bool
}

// PrometheusEnabled check configuration
func (c Config) PrometheusEnabled() bool {
	return c.PrometheusPort > 0
}

// TLSEnabled check configuration
func (c Config) TLSEnabled() bool {
	return c.Cert != "" && c.Key != ""
}

// TracingEnabled check configuration
func (c Config) TracingEnabled() bool {
	return c.JaegerCollector != ""
}

func (c Config) grpcOpts(log *zap.Logger) []grpc.ServerOption {
	interceptors := make([]grpc.UnaryServerInterceptor, 0)
	interceptors = append(interceptors, grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)))
	interceptors = append(interceptors, grpc_zap.UnaryServerInterceptor(log))
	interceptors = append(interceptors, grpc_recovery.UnaryServerInterceptor())
	if c.PrometheusEnabled() {
		interceptors = append(interceptors, grpc_prometheus.UnaryServerInterceptor)
	}
	if c.TracingEnabled() {
		interceptors = append(interceptors, otelgrpc.UnaryServerInterceptor())
	}
	interceptors = append(interceptors, errorMapper)

	opts := make([]grpc.ServerOption, 0)
	opts = append(opts, grpc_middleware.WithUnaryServerChain(interceptors...))
	return opts
}
