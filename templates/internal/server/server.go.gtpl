package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"{{ .GoModule}}/api"	
)
	
// Server represents a gRPC server
type Server struct {		
	cfg  Config
	service api.{{ .Package | UpperFirst}}Server		
	grpcServer *grpc.Server
}

// New gRPC server
func New(cfg Config, service api.{{ .Package | UpperFirst}}Server) *Server {
	return &Server{cfg: cfg, service: service}
}

// ListenAndServe start the server
func (srv *Server) ListenAndServe() error {
	srv.grpcServer = grpc.NewServer(srv.cfg.grpcOpts()...)
	api.Register{{ .Package | UpperFirst}}Server(srv.grpcServer, srv.service)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%s", srv.cfg.Port))
	if err != nil {
		return err
	}

	if srv.cfg.PrometheusEnabled() {
		grpc_prometheus.Register(srv.grpcServer)
		go prometheusServer(srv.cfg.PrometheusPort)
	}

	log.Printf("Server running on port %s...\n", srv.cfg.Port)
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard))
	return srv.grpcServer.Serve(conn)
}

func prometheusServer(port string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	httpServer := &http.Server{
		Addr: "0.0.0.0:" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux,
	}
	log.Printf("Metrics server running on port %s\n", port)
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("unable to start metrics server at port %s: %v", port, err)
	}
}

// Shutdown the server
func (srv *Server) Shutdown(ctx context.Context) error {
	if srv.grpcServer != nil {
		srv.grpcServer.GracefulStop()
	}		

	return nil
}