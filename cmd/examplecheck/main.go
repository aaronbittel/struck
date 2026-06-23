package main

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

type example struct {
	name string
	args []string
	want string
	ran  bool
}

var examples = []*example{
	&example{
		name: "basic.go",
		args: []string{"--name", "Bob", "--age", "42", "-v", "123.456"},
		want: "{Name:Bob Age:42 Verbose:true Input:123.456}",
	},
}

func (e *example) equals(actual string) bool {
	return strings.TrimSpace(e.want) == strings.TrimSpace(actual)
}

const pattern = "./examples/*.go"

func main() {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatal(err)
	}

outer:
	for _, path := range paths {
		for _, example := range examples {
			if example.name != filepath.Base(path) {
				continue
			}
			testExample(path, example)
			example.ran = true
			continue outer
		}

		log.Fatalf("examplecheck: %q has no example entry", path)
	}

	for _, example := range examples {
		if !example.ran {
			log.Fatalf("examplecheck: there is no example for %q", example.name)
		}
	}

	fmt.Println("examplecheck: OK")
}

func testExample(path string, example *example) {
	args := []string{"run", path}
	args = append(args, example.args...)
	cmd := exec.Command("go", args...)

	var (
		stdout strings.Builder
		stderr strings.Builder
	)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("command failed: go %v\nerr: %v\nstderr:\n%s", args, err, stderr.String())
	}

	got := stdout.String()

	if !example.equals(got) {
		log.Fatalf("Example FAIL: %q\nExpected:\n%s\nGot:\n%s\n", path, example.want, got)
	}

	s := stderr.String()
	if s != "" {
		log.Fatalf("stderr:\n%s\n", s)
	}
}
