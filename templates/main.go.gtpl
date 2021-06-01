package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	{{range .Packages}}app_{{.Package}} "{{ .GoModule}}/{{.SrcPath}}"
	{{end}}	"{{ .GoModule}}/internal/server"

	// database driver
	_ {{if eq .Database "mysql"}}"github.com/go-sql-driver/mysql"{{else}}"github.com/jackc/pgx/v4/stdlib"{{end}}
)

const serviceName = "{{ .GoModule}}"

func main() {
	var cfg server.Config
	var dbURL string
	flag.StringVar(&dbURL, "db", "", "The Database connection URL")
	flag.IntVar(&cfg.Port, "port", 5000, "The server port")
	flag.IntVar(&cfg.PrometheusPort, "prometheusPort", 0, "The metrics server port")
	flag.StringVar(&cfg.JaegerAgent, "jaegerAgent", "", "The Jaeger Tracing agent URL")
	flag.StringVar(&cfg.Cert, "cert", "", "The path to the server certificate file in PEM format")
	flag.StringVar(&cfg.Key, "key", "", "The path to the server private key in PEM format")
	flag.Parse()

	log.Printf("Starting %s services...\n", serviceName)
	db, err := sql.Open("{{if eq .Database "mysql"}}mysql{{else}}pgx{{end}}", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(cfg,
	{{range .Packages}}app_{{.Package}}.NewService(app_{{.Package}}.New(db)),
	{{end}})

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-done
		log.Print("Signal detected...")
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Shutdown: %v", err)
	}
}	