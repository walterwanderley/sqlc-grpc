package main

import (
	"context"
	"database/sql"
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
	log.Printf("Starting %s services...\n", serviceName)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	cfg := server.Config{
		ServiceName:       serviceName,
		Port:              port,
		PrometheusPort:    os.Getenv("PROMETHEUS_PORT"),
		JaegerAgent:       os.Getenv("JAEGER_AGENT"),
	}
	var err error
	if err = cfg.Validate(); err != nil {
		log.Fatal(err)
	}	

	dbURL := os.Getenv("DB_URL")
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
