//go:build ignore
// +build ignore

package main

import (
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Delete examples/bin and remake it.
	_ = os.RemoveAll("examples/bin")
	if err := os.MkdirAll("examples/bin", 0777); err != nil {
		panic(err)
	}

	// Loop compiling each item.
	files := strings.Split(os.Getenv("FILES"), " ")
	for _, file := range files {
		// Do preliminary checks.
		if file == "examples/bin" {
			// Ignore this one.
			continue
		}
		s, err := os.Stat(file)
		if err != nil {
			panic(err)
		}
		if !s.IsDir() {
			continue
		}

		// Create the output file path.
		out := strings.Replace(file, "examples", "examples/bin", 1)

		// Compile the folder.
		cmd := exec.Command("go", "build", "-o", "./"+out, "./"+file)
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err = cmd.Run(); err != nil {
			panic(err)
		}
	}
}
