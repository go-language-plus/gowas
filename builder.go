package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildWasm() ([]byte, error) {
	// org command: go build main.wasm
	args := []string{
		"build",
		"-o",
		filepath.Join(staticOutputDir, "main.wasm"),
	}

	if *flagTags != "" {
		args = append(args, "-flagTags", *flagTags)
	}

	// a build root directory could be put without flag, otherwise "." as default
	if len(flag.Args()) > 0 {
		args = append(args, flag.Args()[0])
	} else {
		args = append(args, ".")
	}

	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	if !hasGo111Module(cmd.Env) {
		cmd.Env = append(cmd.Env, "GO111MODULE=on")
	}
	cmd.Dir = "."

	// exec
	log.Print("go ", strings.Join(args, " "))
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error: %v\n", err)
		return stdoutStderr, err
	}
	if string(stdoutStderr) != "" {
		log.Printf("%s\n", stdoutStderr)
	}

	return stdoutStderr, nil
}

func hasGo111Module(env []string) bool {
	for _, e := range env {
		if strings.HasPrefix(e, "GO111MODULE=") {
			return true
		}
	}
	return false
}
