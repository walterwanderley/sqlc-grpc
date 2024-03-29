// Code generated by sqlc-grpc (https://github.com/walterwanderley/sqlc-grpc).

package main

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/flowchartsman/swaggerui"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.0"
	"go.uber.org/automaxprocs/maxprocs"

	// database driver
	_ "github.com/jackc/pgx/v5/stdlib"

	"booktest/internal/server"
	"booktest/internal/server/instrumentation/trace"
)

//go:generate sqlc-grpc -m booktest -tracing -metric -append

const serviceName = "booktest"

var (
	dbURL string

	//go:embed api/apidocs.swagger.json
	openAPISpec []byte
)

func main() {
	cfg := server.Config{
		ServiceName: serviceName,
	}
	var dev bool
	flag.StringVar(&dbURL, "db", "", "The Database connection URL")
	flag.IntVar(&cfg.Port, "port", 5000, "The server port")
	flag.IntVar(&cfg.PrometheusPort, "prometheus-port", 0, "The metrics server port")
	flag.BoolVar(&cfg.EnableCors, "cors", false, "Enable CORS middleware")
	flag.BoolVar(&dev, "dev", false, "Set logger to development mode")
	flag.StringVar(&cfg.OtlpEndpoint, "otlp-endpoint", "", "The Open Telemetry Protocol Endpoint (example: localhost:4317)")

	flag.Parse()

	initLogger(dev)
	if err := run(cfg); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func run(cfg server.Config) error {
	_, err := maxprocs.Set()
	if err != nil {
		slog.Warn("startup", "error", err)
	}
	slog.Info("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	var db *sql.DB
	if cfg.TracingEnabled() {
		flush, err := trace.Init(context.Background(), serviceName, cfg.OtlpEndpoint)
		if err != nil {
			return err
		}
		defer flush()

		db, err = otelsql.Open("pgx", dbURL, otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
		))
		if err != nil {
			return err
		}

		err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
			semconv.DBSystemPostgreSQL,
		))
		if err != nil {
			return err
		}
	} else {

		db, err = sql.Open("pgx", dbURL)
		if err != nil {
			return err
		}
	}
	defer db.Close()

	srv := server.New(cfg, registerServer(db), registerHandlers(), httpHandlers)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-done
		slog.Warn("signal detected...", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	return srv.ListenAndServe()
}

func initLogger(dev bool) {
	var handler slog.Handler
	opts := slog.HandlerOptions{
		AddSource: true,
	}
	switch {
	case dev:
		handler = slog.NewTextHandler(os.Stderr, &opts)
	default:
		handler = slog.NewJSONHandler(os.Stderr, &opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func httpHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(openAPISpec)))

}
