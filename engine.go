package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/emicklei/proto"
	"golang.org/x/tools/imports"

	"github.com/walterwanderley/sqlc-grpc/metadata"
)

//go:embed templates/*
var templates embed.FS

func process(def *metadata.Definition, outPath string, appendMode bool) error {
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
				err := os.MkdirAll(newPath, 0750)
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
				dest := filepath.Join(dir, metadata.ToSnakeCase(pkg.Package), "v1")
				if _, err := os.Stat(dest); os.IsNotExist(err) {
					err := os.MkdirAll(dest, 0750)
					if err != nil {
						return err
					}
				}
				destFile := filepath.Join(dest, (metadata.ToSnakeCase(pkg.Package) + ".proto"))
				if appendMode && fileExists(destFile) {
					loadOptions(destFile, pkg)
				}

				err = genFromTemplate(path, string(tpl), pkg, false, destFile)
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
				err = genFromTemplate(path, string(tpl), pkg, true, filepath.Join(pkg.SrcPath, "service.go"))
				if err != nil {
					return err
				}
			}
			return nil
		}

		if strings.HasSuffix(newPath, "service.factory.go") {
			tpl, err := ioutil.ReadAll(in)
			if err != nil {
				return err
			}
			for _, pkg := range def.Packages {
				newPath := filepath.Join(pkg.SrcPath, "service.factory.go")
				if appendMode && fileExists(newPath) {
					return nil
				}
				err = genFromTemplate(path, string(tpl), pkg, true, newPath)
				if err != nil {
					return err
				}
			}
			return nil
		}

		if strings.HasSuffix(newPath, "adapters.go") {
			tpl, err := ioutil.ReadAll(in)
			if err != nil {
				return err
			}
			for _, pkg := range def.Packages {
				if len(pkg.OutputAdapters) > 0 {
					err = genFromTemplate(path, string(tpl), pkg, true, filepath.Join(pkg.SrcPath, "adapters.go"))
					if err != nil {
						return err
					}
				}
			}
			return nil
		}

		if strings.HasSuffix(path, ".tmpl") {
			tpl, err := ioutil.ReadAll(in)
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

func loadOptions(protoFile string, pkg *metadata.Package) {
	f, err := os.Open(protoFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer f.Close()
	parser := proto.NewParser(f)
	def, err := parser.Parse()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var hasPackageOptions bool
	proto.Walk(def, proto.WithOption(func(opt *proto.Option) {
		if hasPackageOptions {
			return
		}
		if opt.Name == "(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger)" {
			res := make([]string, 0)
			res = append(res, fmt.Sprintf("option %s = {", opt.Name))
			res = append(res, printProtoLiteral(opt.Constant.OrderedMap, 1)...)
			res = append(res, "};")
			pkg.CustomOpenAPIOptions = res
			hasPackageOptions = true
		}
	}))

	proto.Walk(def, proto.WithRPC(func(rpc *proto.RPC) {
		res := make([]string, 0)

		for _, e := range rpc.Elements {
			opt, ok := e.(*proto.Option)
			if !ok {
				continue
			}
			res = append(res, fmt.Sprintf("option %s = {", opt.Name))
			for _, item := range opt.Constant.OrderedMap {
				if item.IsString {
					res = append(res, fmt.Sprintf("    %s: \"%s\"", item.Name, item.Source))
				}
			}
			res = append(res, "};")
		}

		for _, s := range pkg.Services {
			if s.Name == rpc.Name {
				s.CustomHttpOptions = res
				break
			}
		}

	}))

}

func printProtoLiteral(literal proto.LiteralMap, deep int) []string {
	res := make([]string, 0)
	layout := fmt.Sprintf("%%-%ds", deep*4)
	prefix := fmt.Sprintf(layout, "")
	for _, item := range literal {
		if item.IsString {
			res = append(res, fmt.Sprintf("%s%s: \"%s\"", prefix, item.Name, item.Source))
		} else {
			res = append(res, fmt.Sprintf("%s%s: {", prefix, item.Name))
			res = append(res, printProtoLiteral(item.OrderedMap, deep+1)...)
			res = append(res, fmt.Sprintf("%s};", prefix))
		}

	}
	return res
}

func genFromTemplate(name, tmp string, data interface{}, goSource bool, outPath string) error {

	w, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer w.Close()

	var b bytes.Buffer

	funcMap := template.FuncMap{
		"UpperFirst": metadata.UpperFirstCharacter,
		"SnakeCase":  metadata.ToSnakeCase,
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
