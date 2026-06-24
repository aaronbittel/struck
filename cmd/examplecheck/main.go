package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type example struct {
	dirname string
	args    []string
	want    string
	ran     bool
}

var examples = []*example{
	&example{
		dirname: "basic",
		args:    []string{"--name", "Bob", "--age", "42", "-v", "123.456"},
		want:    "{Name:Bob Age:42 Verbose:true Input:123.456}",
	},
	&example{
		dirname: "slices",
		args:    []string{"-v", "true", "--verbose", "t", "--verbose", "1", "--name", "Alice", "-n", "Charlie"},
		want:    "🎉 A most magnificent welcome to Bob, Alice, Charlie! Thank you for gracing this humble program with your presence. 🎉",
	},
}

func (e *example) equals(actual string) bool {
	return strings.TrimSpace(e.want) == strings.TrimSpace(actual)
}

const exampleDir = "./examples/"

func main() {
	entries, err := os.ReadDir(exampleDir)
	if err != nil {
		log.Fatal(err)
	}

outer:
	for _, dir := range entries {
		for _, example := range examples {
			if example.dirname != dir.Name() {
				continue
			}
			testExample(filepath.Join(exampleDir, dir.Name()), example)
			example.ran = true
			continue outer
		}

		log.Fatalf("examplecheck: %q has no example entry", dir)
	}

	for _, example := range examples {
		if !example.ran {
			log.Fatalf("examplecheck: there is no example for %q", example.dirname)
		}
	}

	fmt.Println("examplecheck: OK")
}

func testExample(path string, example *example) {
	args := []string{"run", "./" + path}
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
