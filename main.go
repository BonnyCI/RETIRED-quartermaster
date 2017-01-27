package main

import (
	"os"
	"runtime"

	"github.com/pschwartz/quartermaster/cmd"
	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()

	if jww.LogCountForLevelsGreaterThanorEqualTo(jww.LevelError) > 0 {
		os.Exit(-1)
	}
}
