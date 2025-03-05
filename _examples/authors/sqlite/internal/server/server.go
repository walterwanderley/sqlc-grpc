// Code generated by sqlc-grpc (https://github.com/walterwanderley/sqlc-grpc).

package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	oteltrace "go.opentelemetry.io/otel/trace"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"

	"authors/internal/server/instrumentation/metric"
	"authors/internal/server/instrumentation/trace"
	"authors/internal/server/middleware"
)

const (
	httpReadTimeout  = 15 * time.Second
	httpWriteTimeout = 15 * time.Second
	httpIdleTimeout  = 60 * time.Second
)

type HttpMiddlewareType func(h http.Handler) http.Handler

type RegisterServer func(srv *grpc.Server)

type RegisterHandlerFromEndpoint func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

type RegisterHttpHandler func(mux *http.ServeMux)

// Server represents a gRPC server
type Server struct {
	cfg Config

	grpcServer   *grpc.Server
	healthServer *health.Server
	httpServer   *http.Server

	register             RegisterServer
	registerHandlers     []RegisterHandlerFromEndpoint
	registerHttpHandlers RegisterHttpHandler
}

// New gRPC server
func New(cfg Config, register RegisterServer, registerHandlers []RegisterHandlerFromEndpoint, registerHttpHandler RegisterHttpHandler) *Server {
	return &Server{
		cfg:                  cfg,
		register:             register,
		registerHandlers:     registerHandlers,
		registerHttpHandlers: registerHttpHandler,
	}
}

// ListenAndServe start the server
func (srv *Server) ListenAndServe() error {
	grpcInterceptors := srv.cfg.grpcInterceptors()

	srvMetrics := grpc_prometheus.NewServerMetrics(
		grpc_prometheus.WithServerHandlingTimeHistogram(
			grpc_prometheus.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	if srv.cfg.PrometheusEnabled() {
		prometheus.MustRegister(srvMetrics)
		exemplarFromContext := func(ctx context.Context) prometheus.Labels {
			if span := oteltrace.SpanContextFromContext(ctx); span.IsSampled() {
				return prometheus.Labels{"traceID": span.TraceID().String()}
			}
			return nil
		}
		grpcInterceptors = append(grpcInterceptors, srvMetrics.UnaryServerInterceptor(grpc_prometheus.WithExemplarFromContext(exemplarFromContext)))
	}

	grpcOpts := make([]grpc.ServerOption, 0)
	if srv.cfg.TracingEnabled() {
		grpcOpts = append(grpcOpts, trace.ServerOption())
	}
	grpcOpts = append(grpcOpts, grpc.ChainUnaryInterceptor(grpcInterceptors...))

	srv.grpcServer = grpc.NewServer(grpcOpts...)
	reflection.Register(srv.grpcServer)
	srv.register(srv.grpcServer)

	if srv.cfg.PrometheusEnabled() {
		srvMetrics.InitializeMetrics(srv.grpcServer)
		err := metric.Init(srv.cfg.PrometheusPort, srv.cfg.ServiceName)
		if err != nil {
			return err
		}
	}

	srv.healthServer = health.NewServer()
	healthpb.RegisterHealthServer(srv.grpcServer, srv.healthServer)
	srv.healthServer.SetServingStatus(srv.cfg.ServiceName, healthpb.HealthCheckResponse_SERVING)

	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(annotator),
		runtime.WithForwardResponseOption(forwardResponse),
		runtime.WithOutgoingHeaderMatcher(outcomingHeaderMatcher),
	)
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	sAddr := fmt.Sprintf("dns:///localhost:%d", srv.cfg.Port)
	for _, h := range srv.registerHandlers {
		if err := h(context.Background(), gwmux, sAddr, dialOptions); err != nil {
			return err
		}
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", gwmux)

	if srv.registerHttpHandlers != nil {
		srv.registerHttpHandlers(httpMux)
	}

	srv.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.cfg.Port),
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
		IdleTimeout:  httpIdleTimeout,
		Handler:      grpcHandlerFunc(srv.grpcServer, httpMux),
	}

	if srv.cfg.EnableCors {
		slog.Info("Enable Cross-Origin Resource Sharing")
		srv.httpServer.Handler = middleware.CORS(srv.httpServer.Handler)
	}

	for _, mid := range srv.cfg.Middlewares {
		srv.httpServer.Handler = mid(srv.httpServer.Handler)
	}

	slog.Info("Server is running...", "port", srv.cfg.Port)
	return srv.httpServer.ListenAndServe()
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			if r.URL.Path == "/" {
				http.Redirect(w, r, "/swagger/", http.StatusFound)
				return
			}
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

// Shutdown the server
func (srv *Server) Shutdown(ctx context.Context) {
	srv.healthServer.Shutdown()
	slog.Info("Graceful stop")
	srv.grpcServer.GracefulStop()
	if err := srv.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Shutdown error", "error", err)
	}
}

func annotator(ctx context.Context, req *http.Request) metadata.MD {
	return metadata.New(map[string]string{"requestURI": req.Host + req.URL.RequestURI()})
}

func forwardResponse(ctx context.Context, w http.ResponseWriter, message proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		w.WriteHeader(code)
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
	}

	return nil
}

func outcomingHeaderMatcher(header string) (string, bool) {
	switch header {
	case "location", "authorization", "access-control-expose-headers":
		return header, true
	default:
		return header, false
	}
}
