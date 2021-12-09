package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/walterwanderley/sqlc-grpc/metadata"
)

var (
	module      string
	showVersion bool
	help        bool
)

func main() {
	flag.BoolVar(&help, "h", false, "Help for this program")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.StringVar(&module, "m", "my-project", "Go module name if there are no go.mod")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		fmt.Println("\nFor more information, please visit https://github.com/walterwanderley/sqlc-grpc")
		return
	}

	if showVersion {
		fmt.Println(version)
		return
	}

	cfg, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.Packages) == 0 {
		log.Fatal("no packages")
	}

	if m := moduleFromGoMod(); m != "" {
		fmt.Println("Using module path from go.mod:", m)
		module = m
	}

	def := metadata.Definition{
		GoModule: module,
		Packages: make([]*metadata.Package, 0),
	}

	for _, p := range cfg.Packages {
		pkg, err := metadata.ParsePackage(p.Path)
		if err != nil {
			log.Fatal("parser error:", err.Error())
		}
		pkg.GoModule = module
		pkg.Engine = p.Engine

		def.Packages = append(def.Packages, pkg)
	}
	sort.SliceStable(def.Packages, func(i, j int) bool {
		return strings.Compare(def.Packages[i].Package, def.Packages[j].Package) < 0
	})

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to get working directory:", err.Error())
	}

	err = process(&def, wd)
	if err != nil {
		log.Fatal("unable to process templates:", err.Error())
	}

	postProcess(&def, wd)
}

func moduleFromGoMod() string {
	f, err := os.Open("go.mod")
	if err != nil {
		return ""
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return ""
	}

	return modfile.ModulePath(b)
}

func postProcess(def *metadata.Definition, workingDirectory string) {
	fmt.Printf("Configuring project %s...\n", def.GoModule)
	if err := os.Chdir(workingDirectory); err != nil {
		panic(err)
	}
	execCommand("go mod init " + def.GoModule)
	execCommand("go mod tidy")
	execCommand("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway " +
		"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 " +
		"google.golang.org/protobuf/cmd/protoc-gen-go " +
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc " +
		"github.com/bufbuild/buf/cmd/buf")

	fmt.Println("Running compile.sh...")
	if err := os.Chdir(filepath.Join(workingDirectory, "proto")); err != nil {
		panic(err)
	}
	execCommand("buf mod update")
	for _, pkg := range def.Packages {
		newDir := filepath.Join(workingDirectory, "proto", pkg.Package)
		if _, err := os.Stat(newDir); os.IsNotExist(err) {
			err := os.MkdirAll(newDir, 0750)
			if err != nil {
				panic(err)
			}
		}
		if err := compileProto(pkg.Package); err != nil {
			fmt.Printf("Error on executing compile.sh for package %s: %v\n", pkg.Package, err)
		}
	}

	fmt.Println("Generating OpenAPIv2 specs...")
	execCommand("buf generate --template buf.doc.yaml")

	if err := os.Chdir(workingDirectory); err != nil {
		panic(err)
	}
	execCommand("go mod tidy")
	fmt.Println("Finished!")
}

func compileProto(pkg string) error {
	fmt.Printf("Compiling %s.proto...\n", pkg)
	return execCommand(fmt.Sprintf("buf generate --path %s.proto -o %s", pkg, pkg))
}

func execCommand(command string) error {
	line := strings.Split(command, " ")
	cmd := exec.Command(line[0], line[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("[error] %q: %w", command, err)
	}
	return nil
}
