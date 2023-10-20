package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type TailwindCompiler struct {
	OS             string
	Arch           string
	Version        string
	downloadURLVar string
}

func NewTailwindCompiler() *TailwindCompiler {
	o := runtime.GOOS
	a := runtime.GOARCH
	v := os.Getenv("TAILWIND_VERSION")
	tc := &TailwindCompiler{
		OS:             o,
		Arch:           a,
		Version:        v,
		downloadURLVar: fmt.Sprintf("TAILWIND_DOWNLOAD_URL_%s_%s", strings.ToUpper(o), strings.ToUpper(a)),
	}
	return tc
}

func (t *TailwindCompiler) downloadCompiler(projectName string, dest string) {
	dl := os.Getenv(t.downloadURLVar)
	err := downloadFile(dl, dest, projectName)
	if err != nil {
		panic(err)
	}
}
