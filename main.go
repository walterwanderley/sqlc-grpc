package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/walterwanderley/sqlc-grpc/metadata"
)

var (
	module             string
	ignoreQueries      string
	migrationPath      string
	liteFS             bool
	distributedTracing bool
	appendMode         bool
	showVersion        bool
	help               bool
)

func main() {
	flag.BoolVar(&help, "h", false, "Help for this program")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.BoolVar(&appendMode, "append", false, "Enable append mode. Don't rewrite editable files")
	flag.StringVar(&module, "m", "my-project", "Go module name if there are no go.mod")
	flag.StringVar(&ignoreQueries, "i", "", "Comma separated list (regex) of queries to ignore")
	flag.StringVar(&migrationPath, "migration-path", "", "Path to migration directory")
	flag.BoolVar(&liteFS, "litefs", false, "Enable support to LiteFS")
	flag.BoolVar(&distributedTracing, "tracing", false, "Enable support to distributed tracing with Jeager")
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

	if migrationPath != "" {
		fi, err := os.Stat(migrationPath)
		if err != nil {
			log.Fatal("invalid -migration-path: ", err.Error())
		}
		if !fi.IsDir() {
			log.Fatal("-migration-path must be a directory")
		}
	}

	cfg, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.Packages) == 0 {
		log.Fatal("no packages")
	}

	queriesToIgnore := make([]*regexp.Regexp, 0)
	for _, queryName := range strings.Split(ignoreQueries, ",") {
		s := strings.TrimSpace(queryName)
		if s == "" {
			continue
		}
		queriesToIgnore = append(queriesToIgnore, regexp.MustCompile(s))
	}

	if m := moduleFromGoMod(); m != "" {
		log.Println("Using module path from go.mod:", m)
		module = m
	}

	args := strings.Join(os.Args, " ")
	if !strings.Contains(args, " -append") {
		args += " -append"
	}

	def := metadata.Definition{
		Args:               args,
		GoModule:           module,
		MigrationPath:      migrationPath,
		Packages:           make([]*metadata.Package, 0),
		LiteFS:             liteFS,
		DistributedTracing: distributedTracing,
	}

	for _, p := range cfg.Packages {
		pkg, err := metadata.ParsePackage(metadata.PackageOpts{
			Path:               p.Path,
			EmitInterface:      p.EmitInterface,
			EmitParamsPointers: p.EmitParamsStructPointers,
			EmitResultPointers: p.EmitResultStructPointers,
			EmitDbArgument:     p.EmitMethodsWithDBArgument,
		}, queriesToIgnore)
		if err != nil {
			log.Fatal("parser error:", err.Error())
		}
		pkg.GoModule = module
		pkg.Engine = p.Engine

		if len(pkg.Services) == 0 {
			log.Println("No services on package", pkg.Package)
			continue
		}

		def.Packages = append(def.Packages, pkg)
	}
	sort.SliceStable(def.Packages, func(i, j int) bool {
		return strings.Compare(def.Packages[i].Package, def.Packages[j].Package) < 0
	})

	if len(def.Packages) == 0 {
		log.Fatal("No services found, verify the -i parameter")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to get working directory:", err.Error())
	}

	err = process(&def, wd, appendMode)
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
	log.Printf("Configuring project %s...\n", def.GoModule)
	execCommand("go mod init " + def.GoModule)
	execCommand("go mod tidy")
	execCommand("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway")
	execCommand("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2")
	execCommand("go install google.golang.org/protobuf/cmd/protoc-gen-go")
	execCommand("go install google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	execCommand("go install github.com/bufbuild/buf/cmd/buf")
	log.Println("Compiling protocol buffers...")
	if err := os.Chdir("proto"); err != nil {
		panic(err)
	}
	execCommand("buf mod update")
	if err := os.Chdir(workingDirectory); err != nil {
		panic(err)
	}
	execCommand("buf generate")
	execCommand("buf format -w")
	execCommand("go mod tidy")
	log.Println("Finished!")
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
