package gosh

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Outputs struct {
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}

func checkAbsAndExist(path string) error {
	if !filepath.IsAbs(path) {
		return ErrPathIsNotAbs
	}
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return ErrPathNotFound
	}

	return nil
}

func getInterpreter(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	shebang, err := r.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	if strings.HasPrefix(shebang, "#!") {
		// found shebang
		return strings.TrimSpace(strings.TrimPrefix(shebang, "#!")), nil
	} else {
		// no shebang found
		return "", nil
	}
}

/*
 * file, interpreter, wd, if not empty, all should be abs
 */
func Run(file, interpreter, wd string, sync bool, arg []string) (Outputs, *exec.Cmd, error) {

	if err := checkAbsAndExist(file); err != nil {
		return Outputs{}, nil, err
	}

	if interpreter == "" {
		var err error
		interpreter, err = getInterpreter(file)
		if err != nil {
			return Outputs{}, nil, err
		}
	}
	if interpreter == "" {
		// fallback to default sh
		interpreter = "/bin/sh"
	}
	if err := checkAbsAndExist(interpreter); err != nil {
		return Outputs{}, nil, err
	}

	if wd != "" {
		if err := checkAbsAndExist(wd); err != nil {
			return Outputs{}, nil, err
		}
	}

	pwd, err := os.Getwd()
	if err == nil {
		if wd == "" {
			wd = pwd
		}
	}

	arg = append([]string{file}, arg...)
	cmd := exec.Command(interpreter, arg...)
	cmd.Dir = wd

	if sync {
		outputs, err := runSync(cmd)
		return outputs, nil, err
	} else {
		cmd, err := runAsync(cmd)
		return Outputs{}, cmd, err
	}
}

func runSync(cmd *exec.Cmd) (Outputs, error) {
	var ops Outputs
	ops.Stdout = new(bytes.Buffer)
	ops.Stderr = new(bytes.Buffer)

	cmd.Stdout = ops.Stdout
	cmd.Stderr = ops.Stderr
	err := cmd.Run()
	return ops, err
}

func runAsync(cmd *exec.Cmd) (*exec.Cmd, error) {
	err := cmd.Start()
	return cmd, err
}
