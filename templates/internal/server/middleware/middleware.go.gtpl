package middleware

import (
	"context"		
	"io"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	grpc_tracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

func Logging(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	if err != nil {
		grpclog.Errorf("Service: %s | Input: %v | Duration: %s | Error: %v\n",
			info.FullMethod,
			req,
			time.Since(start),
			err)
	} else {
		grpclog.Infof("Service: %s | Input: %v | Duration: %s\n",
			info.FullMethod,
			req,
			time.Since(start))
	}

	return h, err
}

func Tracing(serviceName string, agentHost string) grpc.UnaryServerInterceptor {
	t, _, err := newTracer(serviceName, agentHost)
	if err != nil {
		grpclog.Fatalf("unable to start tracer at %s: %v", agentHost, err)
	}
	return grpc_tracing.UnaryServerInterceptor(grpc_tracing.WithTracer(t))
}

func newTracer(serviceName string, agentHost string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentHost,
		},
		ServiceName: serviceName,
	}

	tracer, closer, err = cfg.NewTracer(
		jaegercfg.Logger(jaeger.StdLogger),
	)
	if err != nil {
		return
	}

	opentracing.SetGlobalTracer(tracer)
	return
}