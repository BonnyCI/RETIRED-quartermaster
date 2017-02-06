package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/bonnyci/quartermaster/helpers"
	"github.com/bonnyci/quartermaster/lib"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of quartermaster",
	Long:  `Display version information pertaining to quartermaster`,
	RunE: func(cmd *cobra.Command, args []string) error {
		printQuartermasterVersion()
		return nil
	},
}

func printQuartermasterVersion() {
	if lib.BuildDate == "" {
		setBuildDate()
	} else {
		formatBuildDate()
	}

	if lib.CommitHash == "" {
		jww.FEEDBACK.Printf("Quartermaster v%s %s/%s BuildDate: %s\n", helpers.QuartermasterVersion(), runtime.GOOS, runtime.GOARCH, lib.BuildDate)
	} else {
		jww.FEEDBACK.Printf("Quartermaster v%s-%s %s/%s BuildDate: %s\n", helpers.QuartermasterVersion(), strings.ToUpper(lib.CommitHash), runtime.GOOS, runtime.GOARCH, lib.BuildDate)
	}
}

// setBuildDate checks the ModTime of the quartermaster executable and returns it as a
// formatted string.  This assumes that the executable name is quartermaster, if it does
// not exist, an empty string will be returned.  This is only called if the
// lib.BuildDate wasn't set during compile time.
//
// osext is used for cross-platform.
func setBuildDate() {
	fname, _ := osext.Executable()
	dir, err := filepath.Abs(filepath.Dir(fname))
	if err != nil {
		jww.ERROR.Println(err)
		return
	}
	fi, err := os.Lstat(filepath.Join(dir, filepath.Base(fname)))
	if err != nil {
		jww.ERROR.Println(err)
		return
	}
	t := fi.ModTime()
	lib.BuildDate = t.Format(time.RFC3339)
}

// formatBuildDate formats the lib.BuildDate according to the value in
// .Params.DateFormat, if it's set.
func formatBuildDate() {
	t, _ := time.Parse("2006-01-02T15:04:05-0700", lib.BuildDate)
	lib.BuildDate = t.Format(time.RFC3339)
}
