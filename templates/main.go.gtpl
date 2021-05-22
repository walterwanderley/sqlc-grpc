package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"{{ .GoModule}}/internal/application"
	database "{{ .GoModule}}/{{.SrcPath}}"
	"{{ .GoModule}}/internal/server"

	// database driver
	_ {{if eq .Database "mysql"}}"github.com/go-sql-driver/mysql"{{else}}"github.com/jackc/pgx/v4/stdlib"{{end}}
)

const serviceName = "{{ .Package}}"

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

	queries := database.New(db)
	srv := server.New(cfg, application.NewService(queries))

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
