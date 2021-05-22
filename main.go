package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/walterwanderley/sqlc-grpc/metadata"
)

var module string
var help bool

func main() {
	flag.BoolVar(&help, "h", false, "Help for this program")
	flag.StringVar(&module, "m", "my-project", "Go module name if there are no go.mod")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		return
	}

	cfg, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.Packages) != 1 {
		log.Fatal("multiple packages aren't supported yet")
	}

	if m := moduleFromGoMod(); m != "" {
		fmt.Println("Using module path from go.mod:", m)
		module = m
	}

	def, err := metadata.ParseDefinition(cfg.Packages[0].Path, cfg.Packages[0].Engine, module)
	if err != nil {
		log.Fatal("parser error:", err.Error())
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to get working directory:", err.Error())
	}

	err = process(def, wd)
	if err != nil {
		log.Fatal("unable to process templates:", err.Error())
	}

	postProcess(module)
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

func postProcess(module string) {
	fmt.Printf("Configuring project %s...\n", module)
	execCommand("go mod init " + module)
	execCommand("go mod tidy")

	fmt.Println("Running compile.sh...")
	if err := os.Chdir("api"); err != nil {
		panic(err)
	}
	fmt.Println("Generating protocol buffer...")
	err := execCommand("protoc -I. -Ideps --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative service.proto")
	if err != nil {
		fmt.Println("error calling protoc:", err.Error())
	}

	if err := os.Chdir("../"); err != nil {
		panic(err)
	}

	execCommand("go mod tidy")

	fmt.Println("Finished!")
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
