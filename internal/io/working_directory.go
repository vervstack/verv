package io

import (
	"os"
)

var workingDirectory string

//nolint:gochecknoinits // GetWd() has no error return; must resolve the working dir once before any command reads it, fail fast on error
func init() {
	var err error
	workingDirectory, err = os.Getwd()
	if err != nil {
		panic("cannot obtain working directory path:" + err.Error())
	}
}

func GetWd() string {
	return workingDirectory
}
