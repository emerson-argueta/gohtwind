package cmds

import (
	"fmt"
	"gohtwind/utils"
	"runtime"
	"strings"
)

type TailwindCompiler struct {
	OS             string
	Arch           string
	Version        string
	downloadURLVar string
	envMap         map[string]string
}

func NewTailwindCompiler(envMap map[string]string) *TailwindCompiler {
	o := runtime.GOOS
	a := runtime.GOARCH
	v := envMap["TAILWIND_VERSION"]
	tc := &TailwindCompiler{
		OS:             o,
		Arch:           a,
		Version:        v,
		downloadURLVar: fmt.Sprintf("TAILWIND_DOWNLOAD_URL_%s_%s", strings.ToUpper(o), strings.ToUpper(a)),
		envMap:         envMap,
	}
	return tc
}

func (t *TailwindCompiler) downloadCompiler(projectName string, dest string) {
	dl := t.envMap[t.downloadURLVar]
	err := utils.DownloadFile(dl, dest, projectName)
	if err != nil {
		panic(err)
	}
}
