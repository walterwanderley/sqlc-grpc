package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/flowchartsman/swaggerui"
	"github.com/fullstorydev/grpcui/standalone"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"

	"booktest/internal/server/middleware"
)

const (
	httpReadTimeout  = 15 * time.Second
	httpWriteTimeout = 15 * time.Second
	httpIdleTimeout  = 60 * time.Second

	startupTimeout = 2 * time.Minute
)

type RegisterServer func(srv *grpc.Server)

type RegisterHandler func(ctx context.Context, mux *runtime.ServeMux, cc *grpc.ClientConn) error

// Server represents a gRPC server
type Server struct {
	cfg Config
	log *zap.Logger

	grpcServer   *grpc.Server
	healthServer *health.Server

	register         RegisterServer
	registerHandlers []RegisterHandler
	openAPISpec      []byte
}

// New gRPC server
func New(cfg Config, log *zap.Logger, register RegisterServer, registerHandlers []RegisterHandler, openAPISpec []byte) *Server {
	return &Server{
		cfg:              cfg,
		log:              log,
		register:         register,
		registerHandlers: registerHandlers,
		openAPISpec:      openAPISpec,
	}
}

// ListenAndServe start the server
func (srv *Server) ListenAndServe() error {
	grpc_zap.ReplaceGrpcLoggerV2(srv.log)
	srv.grpcServer = grpc.NewServer(srv.cfg.grpcOpts(srv.log)...)
	reflection.Register(srv.grpcServer)
	srv.register(srv.grpcServer)

	srv.healthServer = health.NewServer()
	healthpb.RegisterHealthServer(srv.grpcServer, srv.healthServer)
	srv.healthServer.SetServingStatus("ww", healthpb.HealthCheckResponse_SERVING)

	var listen net.Listener
	dialOptions := []grpc.DialOption{grpc.WithBlock()}
	var schema string
	if srv.cfg.TLSEnabled() {
		schema = "https"
		tlsCert, err := tls.LoadX509KeyPair(srv.cfg.Cert, srv.cfg.Key)
		if err != nil {
			return fmt.Errorf("failed to parse certificate and key: %w", err)
		}
		tlsCert.Leaf, _ = x509.ParseCertificate(tlsCert.Certificate[0])
		tc := &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			MinVersion:   tls.VersionTLS12,
		}
		listen, err = tls.Listen("tcp", fmt.Sprintf(":%d", srv.cfg.Port), tc)
		if err != nil {
			return err
		}

		cp := x509.NewCertPool()
		cp.AddCert(tlsCert.Leaf)
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cp, "")))
	} else {
		schema = "http"
		var err error
		listen, err = net.Listen("tcp", fmt.Sprintf(":%d", srv.cfg.Port))
		if err != nil {
			return err
		}
		dialOptions = append(dialOptions, grpc.WithInsecure())
	}

	mux := cmux.New(listen)
	grpcListener := mux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := mux.Match(cmux.Any())

	go func() {
		if err := mux.Serve(); err != nil {
			srv.log.Error("failed to serve cmux", zap.Error(err))
		}
	}()

	if srv.cfg.PrometheusEnabled() {
		grpc_prometheus.Register(srv.grpcServer)
		go prometheusServer(srv.log, srv.cfg.PrometheusPort)
	}

	go func() {
		srv.log.Info("Server running", zap.String("addr", grpcListener.Addr().String()))
		if err := srv.grpcServer.Serve(grpcListener); err != nil {
			srv.log.Fatal("Failed to start gRPC Server", zap.Error(err))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), startupTimeout)
	defer cancel()

	sAddr := fmt.Sprintf("dns:///localhost:%d", srv.cfg.Port)
	cc, err := grpc.DialContext(
		ctx,
		sAddr,
		dialOptions...,
	)
	if err != nil {
		return err
	}
	defer cc.Close()

	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(annotator),
		runtime.WithForwardResponseOption(forwardResponse),
		runtime.WithOutgoingHeaderMatcher(outcomingHeaderMatcher),
	)

	for _, h := range srv.registerHandlers {
		if err := h(context.Background(), gwmux, cc); err != nil {
			return err
		}
	}

	httpMux := http.NewServeMux()

	if srv.cfg.EnableGrpcUI {
		grpcui, err := standalone.HandlerViaReflection(ctx, cc, sAddr)
		if err != nil {
			return err
		}

		httpMux.Handle("/grpcui/", http.StripPrefix("/grpcui", grpcui))
		srv.log.Info(fmt.Sprintf("Serving gRPC UI on %s://localhost:%d/grpcui", schema, srv.cfg.Port))
	}

	httpMux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(srv.openAPISpec)))
	srv.log.Info(fmt.Sprintf("Serving Swagger UI on %s://localhost:%d/swagger", schema, srv.cfg.Port))

	httpMux.Handle("/", gwmux)

	httpServer := &http.Server{
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
		IdleTimeout:  httpIdleTimeout,
		Handler:      httpMux,
	}

	if srv.cfg.EnableCors {
		srv.log.Info("Enable Cross-Origin Resource Sharing")
		httpServer.Handler = middleware.CORS(httpMux)
	}

	return httpServer.Serve(httpListener)
}

func prometheusServer(log *zap.Logger, port int) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
		IdleTimeout:  httpIdleTimeout,
		Handler:      mux,
	}
	log.Info("Metrics server running", zap.Int("port", port))
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("unable to start metrics server", zap.Error(err), zap.Int("port", port))
	}
}

// Shutdown the server
func (srv *Server) Shutdown() {
	srv.healthServer.Shutdown()
	srv.log.Info("Graceful stop")
	srv.grpcServer.GracefulStop()
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
