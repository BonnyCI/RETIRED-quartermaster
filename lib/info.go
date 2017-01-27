package lib

import "github.com/pschwartz/quartermaster/helpers"

var (
	// CommitHash contains the current Git revision. Use make to build to make
	// sure this gets set.
	CommitHash string

	// BuildDate contains the date of the current build.
	BuildDate string
)

var info *Info

// StandupbotInfo contains information about the current Standupbot environment
type Info struct {
	Version    string
	CommitHash string
	BuildDate  string
}

func init() {
	info = &Info{
		Version:    helpers.QuartermasterVersion(),
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
	}
}
