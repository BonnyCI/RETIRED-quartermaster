package helpers

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// FilePathSeparator as defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

// NormalizeStandupbotFlags facilitates transitions of Standupbot command-line flags,
// e.g. --baseUrl to --baseURL
func NormalizeQuartermasterFlags(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "baseUrl":
		name = "baseURL"
		break
	}
	return pflag.NormalizedName(name)
}

func CreateLogger(logLevel string) error {
	var (
		logHandle       = ioutil.Discard
		stdoutThreshold jww.Threshold
		logThreshold    jww.Threshold
	)

	switch lvl := logLevel; lvl {
	case "info":
		stdoutThreshold = jww.LevelInfo
		logThreshold = jww.LevelInfo
	case "debug":
		stdoutThreshold = jww.LevelDebug
		logThreshold = jww.LevelDebug
	case "warn":
		stdoutThreshold = jww.LevelWarn
		logThreshold = jww.LevelWarn
	default:
		jww.DEBUG.Println("Setting log level to ", lvl)
		stdoutThreshold = jww.LevelError
		logThreshold = jww.LevelWarn
	}

	var err error
	if viper.IsSet("logFile") && viper.GetString("logFile") != "" {
		path := viper.GetString("logFile")
		logHandle, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			return newSystemError("Failed to open log file:", path, err)
		}
	} else {
		logHandle, err = ioutil.TempFile(os.TempDir(), "standupbot")
		if err != nil {
			return newSystemError(err)
		}
	}

	jww.SetStdoutThreshold(stdoutThreshold)
	jww.SetLogThreshold(logThreshold)
	jww.SetLogOutput(logHandle)
	jww.SetPrefix("")
	jww.SetFlags(log.Ldate | log.Ltime)

	return nil
}
