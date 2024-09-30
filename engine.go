package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/imports"

	"github.com/walterwanderley/sqlc-grpc/converter"
	"github.com/walterwanderley/sqlc-grpc/metadata"
	"github.com/walterwanderley/sqlc-grpc/templates"
)

func process(def *metadata.Definition, appendMode bool) error {
	return fs.WalkDir(templates.Files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println("ERROR ", err.Error())
			return err
		}

		newPath := strings.TrimSuffix(path, ".tmpl")

		if d.IsDir() {
			if strings.HasSuffix(newPath, "instrumentation") && (!def.DistributedTracing && !def.Metric) {
				return nil
			}

			if strings.HasSuffix(newPath, "trace") && !def.DistributedTracing {
				return nil
			}
			if strings.HasSuffix(newPath, "metric") && !def.Metric {
				return nil
			}
			if strings.HasSuffix(newPath, "litestream") && !(def.Database() == "sqlite" && def.Litestream) {
				return nil
			}

			if strings.HasSuffix(newPath, "litefs") && !(def.Database() == "sqlite" && def.LiteFS) {
				return nil
			}
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				err := os.MkdirAll(newPath, 0o750)
				if err != nil {
					return err
				}
			}
			return nil
		}

		if strings.HasSuffix(newPath, "templates.go") {
			return nil
		}

		log.Println(path, "...")

		in, err := templates.Files.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		if strings.HasSuffix(newPath, "service.proto") {
			dir := strings.TrimSuffix(newPath, "service.proto")
			tpl, err := io.ReadAll(in)
			if err != nil {
				return err
			}
			for _, pkg := range def.Packages {
				dest := filepath.Join(dir, converter.ToSnakeCase(pkg.Package), "v1")
				if _, err := os.Stat(dest); os.IsNotExist(err) {
					err := os.MkdirAll(dest, 0o750)
					if err != nil {
						return err
					}
				}
				destFile := filepath.Join(dest, (converter.ToSnakeCase(pkg.Package) + ".proto"))
				if appendMode && fileExists(destFile) {
					pkg.LoadOptions(destFile)
				}

				err = genFromTemplate(path, string(tpl), pkg, false, destFile)
				if err != nil {
					return err
				}
			}
			return nil
		}

		if strings.HasSuffix(newPath, "service.go") {
			tpl, err := io.ReadAll(in)
			if err != nil {
				return err
			}
			for _, pkg := range def.Packages {
				err = genFromTemplate(path, string(tpl), pkg, true, filepath.Join(pkg.SrcPath, "service.go"))
				if err != nil {
					return err
				}
			}
			return nil
		}

		if strings.HasSuffix(newPath, "service.factory.go") {
			tpl, err := io.ReadAll(in)
			if err != nil {
				return err
			}
			for _, pkg := range def.Packages {
				newPath := filepath.Join(pkg.SrcPath, "service.factory.go")
				if appendMode && fileExists(newPath) {
					continue
				}
				err = genFromTemplate(path, string(tpl), pkg, true, newPath)
				if err != nil {
					return err
				}
			}
			return nil
		}

		if strings.HasSuffix(newPath, "adapters.go") {
			tpl, err := io.ReadAll(in)
			if err != nil {
				return err
			}
			for _, pkg := range def.Packages {
				if len(pkg.OutputAdapters) > 0 || pkg.HasExecResult {
					err = genFromTemplate(path, string(tpl), pkg, true, filepath.Join(pkg.SrcPath, "adapters.go"))
					if err != nil {
						return err
					}
				}
			}
			return nil
		}

		if strings.HasSuffix(newPath, "tracing.go") && !def.DistributedTracing {
			return nil
		}

		if strings.HasSuffix(newPath, "metric.go") && !def.Metric {
			return nil
		}

		if strings.HasSuffix(newPath, "migration.go") && def.MigrationPath == "" {
			return nil
		}

		if strings.HasSuffix(newPath, "litestream.go") && !(def.Database() == "sqlite" && def.Litestream) {
			return nil
		}

		if (strings.HasSuffix(newPath, "litefs.go") || strings.HasSuffix(newPath, "forward.go")) && !(def.Database() == "sqlite" && def.LiteFS) {
			return nil
		}

		if strings.HasSuffix(path, ".tmpl") {
			tpl, err := io.ReadAll(in)
			if err != nil {
				return err
			}
			goCode := strings.HasSuffix(newPath, ".go")
			if goCode && appendMode && fileExists(newPath) && !strings.HasSuffix(newPath, "registry.go") {
				return nil
			}
			return genFromTemplate(path, string(tpl), def, goCode, newPath)
		}

		if appendMode && fileExists(newPath) {
			return nil
		}

		out, err := os.Create(newPath)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		return err
	})
}

func genFromTemplate(name, tmp string, data interface{}, goSource bool, outPath string) error {
	w, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer w.Close()

	var b bytes.Buffer

	t, err := template.New(name).Funcs(templates.Funcs).Parse(tmp)
	if err != nil {
		return err
	}
	err = t.Execute(&b, data)
	if err != nil {
		return fmt.Errorf("execute template error: %w", err)
	}

	var src []byte
	if goSource {
		src, err = imports.Process("", b.Bytes(), nil)
		if err != nil {
			return fmt.Errorf("organize imports error: %w", err)
		}
	} else {
		src = b.Bytes()
	}

	fmt.Fprintf(w, "%s", string(src))
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
