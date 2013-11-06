package gosh

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

const (
	shWithShebang    = "#!/bin/sh\necho \"Hello gosh\""
	shWithoutShebang = "echo \"Hello gosh\""

	shWithShebangArg    = "#!/bin/sh\necho $1"
	shWithoutShebangArg = "echo $1"
)

func generateScript(content string) (string, error) {
	f, err := ioutil.TempFile("", "gosh")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.WriteString(f, content)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func TestRunShWithoutArg(t *testing.T) {
	path, err := generateScript(shWithoutShebang)
	if err != nil {
		t.Fatalf("Generate sh script failed: %v", err)
	}
	defer os.Remove(path)

	ops, _, err := Run(path, "/bin/sh", "", true, []string{})
	if err != nil {
		t.Fatalf("Run returns error: %v", err)
	}
	if ops.Stdout.String() != "Hello gosh\n" {
		t.Fatalf("Script output wrong: %v", ops.Stdout.String())
	}
}

func TestRunShWithArg(t *testing.T) {
	path, err := generateScript(shWithoutShebangArg)
	if err != nil {
		t.Fatalf("Generate sh script failed: %v", err)
	}
	defer os.Remove(path)

	ops, _, err := Run(path, "/bin/sh", "", true, []string{"Hello gosh"})
	if err != nil {
		t.Fatalf("Run returns error: %v", err)
	}
	if ops.Stdout.String() != "Hello gosh\n" {
		t.Fatalf("Script output wrong: %v", ops.Stdout.String())
	}
}
