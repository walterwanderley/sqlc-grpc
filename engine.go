package main

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/walterwanderley/sqlc-grpc/metadata"

	"golang.org/x/tools/imports"
)

//go:embed templates/*
var templates embed.FS

func process(def *metadata.Definition, outPath string) error {
	rootPath := "templates"
	return fs.WalkDir(templates, rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("ERROR ", err.Error())
			return err
		}

		newPath := strings.Replace(path, rootPath, outPath, 1)
		newPath = strings.TrimSuffix(newPath, ".tmpl")

		if d.IsDir() {
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				err := os.MkdirAll(newPath, 0777)
				if err != nil {
					return err
				}
			}
			return nil
		}

		fmt.Println(path, "...")

		in, err := templates.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		if strings.HasSuffix(newPath, "service.proto") {
			dir := strings.TrimSuffix(newPath, "service.proto")
			tpl, err := ioutil.ReadAll(in)
			if err != nil {
				return err
			}
			for _, pkg := range def.Packages {
				out, err := os.Create(filepath.Join(dir, (pkg.Package + ".proto")))
				if err != nil {
					return err
				}
				defer out.Close()

				err = genFromTemplate(path, string(tpl), pkg, false, out)
				if err != nil {
					return err
				}
			}
			return nil
		}

		if strings.HasSuffix(newPath, "service.go") {
			tpl, err := ioutil.ReadAll(in)
			if err != nil {
				return err
			}
			for _, pkg := range def.Packages {
				out, err := os.Create(filepath.Join(pkg.SrcPath, "service.go"))
				if err != nil {
					return err
				}
				defer out.Close()

				err = genFromTemplate(path, string(tpl), pkg, true, out)
				if err != nil {
					return err
				}
			}
			return nil
		}

		out, err := os.Create(newPath)
		if err != nil {
			return err
		}
		defer out.Close()

		if strings.HasSuffix(newPath, ".sh") {
			err = out.Chmod(os.ModePerm)
			if err != nil {
				return err
			}
		}

		if strings.HasSuffix(path, ".tmpl") {
			tpl, err := ioutil.ReadAll(in)
			if err != nil {
				return err
			}
			goCode := strings.HasSuffix(newPath, ".go")
			return genFromTemplate(path, string(tpl), def, goCode, out)
		}

		_, err = io.Copy(out, in)
		return err
	})
}

func genFromTemplate(name, tmp string, data interface{}, goSource bool, w io.Writer) error {

	var b bytes.Buffer

	funcMap := template.FuncMap{
		"UpperFirst": metadata.UpperFirstCharacter,
	}

	t, err := template.New(name).Funcs(funcMap).Parse(tmp)
	if err != nil {
		return err
	}
	err = t.Execute(&b, data)
	if err != nil {
		return fmt.Errorf("execute template error: %w", err)
	}

	var src []byte
	if goSource {
		src, err = format.Source(b.Bytes())
		if err != nil {
			fmt.Println(b.String())
			return fmt.Errorf("format source error: %w", err)
		}
		src, err = imports.Process("", src, nil)
		if err != nil {
			return fmt.Errorf("organize imports error: %w", err)
		}
	} else {
		src = b.Bytes()
	}

	fmt.Fprintf(w, "%s", string(src))
	return nil

}
