package gosh

import (
	"fmt"
)

var (
	ErrPathIsNotAbs = fmt.Errorf("given path is not abs.")
	ErrPathNotFound = fmt.Errorf("Script file not found.")
)
