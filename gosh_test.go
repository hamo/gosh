package gosh

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	shWithShebang    = "#!/bin/sh\necho \"Hello gosh\""
	shWithoutShebang = "echo \"Hello gosh\""

	shWithShebangArg    = "#!/bin/sh\necho $1"
	shWithoutShebangArg = "echo $1"

	perlWithShebang    = "#!/usr/bin/perl\nprint \"Hello gosh\\n\";"
	perlWithoutShebang = "print \"Hello gosh\\n\";"
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

func TestRunShWithShebang(t *testing.T) {
	path, err := generateScript(shWithShebang)
	if err != nil {
		t.Fatalf("Generate sh script failed: %v", err)
	}
	defer os.Remove(path)

	ops, _, err := Run(path, "", "", true, []string{})
	if err != nil {
		t.Fatalf("Run returns error: %v", err)
	}
	if ops.Stdout.String() != "Hello gosh\n" {
		t.Fatalf("Script output wrong: %v", ops.Stdout.String())
	}
}

func TestRunPerlWithoutArg(t *testing.T) {
	path, err := generateScript(perlWithoutShebang)
	if err != nil {
		t.Fatalf("Generate sh script failed: %v", err)
	}
	defer os.Remove(path)

	ops, _, err := Run(path, "/usr/bin/perl", "", true, []string{})
	if err != nil {
		t.Fatalf("Run returns error: %v", err)
	}
	if ops.Stdout.String() != "Hello gosh\n" {
		t.Fatalf("Script output wrong: %v", ops.Stdout.String())
	}
}

func TestRunPerlWithShebang(t *testing.T) {
	path, err := generateScript(perlWithShebang)
	if err != nil {
		t.Fatalf("Generate sh script failed: %v", err)
	}
	defer os.Remove(path)

	ops, _, err := Run(path, "", "", true, []string{})
	if err != nil {
		t.Fatalf("Run returns error: %v", err)
	}
	if ops.Stdout.String() != "Hello gosh\n" {
		t.Fatalf("Script output wrong: %v", ops.Stdout.String())
	}
}

func TestRunShUnderWd(t *testing.T) {
	tmpdir := filepath.Join(os.TempDir(), "goshdir")
	err := os.Mkdir(tmpdir, 0755)
	if err != nil {
		t.Fatalf("Make temp work dir failed: %v", err)
	}
	defer os.Remove(tmpdir)

	contentFile := filepath.Join(tmpdir, "content")
	f, err := os.Create(contentFile)
	if err != nil {
		t.Fatalf("Create content file failed: %v", err)
	}
	_, err = f.WriteString("Hello gosh\n")
	if err != nil {
		t.Fatalf("Write content to content file failed: %v", err)
	}
	defer os.Remove(contentFile)
	defer f.Close()

	path, err := generateScript("cat content")
	if err != nil {
		t.Fatalf("Generate sh script failed: %v", err)
	}
	defer os.Remove(path)

	ops, _, err := Run(path, "/bin/sh", tmpdir, true, []string{})
	if err != nil {
		t.Fatalf("Run returns error: %v", err)
	}
	if ops.Stdout.String() != "Hello gosh\n" {
		t.Fatalf("Script output wrong: %v", ops.Stdout.String())
	}

}
