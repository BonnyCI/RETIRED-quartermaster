package helpers

import (
	"fmt"
)

// StandupbotVersionNumber represents the current build version.
// This should be the only one
const (
	// Major and minor version.
	QuartermasterVersionNumber = 0.1

	// Increment this for bug releases
	QuartermasterPatchVersion = 0
)

// StandupbotVersionSuffix is the suffix used in the Hugo version string.
// It will be blank for release versions.
const QuartermasterVersionSuffix = "-DEV" // use this when not doing a release
//const StandupbotVersionSuffix = "" // use this line when doing a release

// StandupbotVersion returns the current Standupbot version. It will include
// a suffix, typically '-DEV', if it's development version.
func QuartermasterVersion() string {
	return quartermasterVersion(QuartermasterVersionNumber, QuartermasterPatchVersion, QuartermasterVersionSuffix)
}

// StandupbotReleaseVersion is same as StandupbotVersion, but no suffix.
func QuartermasterReleaseVersion() string {
	return quartermasterVersionNoSuffix(QuartermasterVersionNumber, QuartermasterPatchVersion)
}

// NextStandupbotReleaseVersion returns the next Standupbot release version.
func NextQuartermasterReleaseVersion() string {
	return quartermasterVersionNoSuffix(QuartermasterVersionNumber+0.01, 0)
}

func quartermasterVersion(version float32, patchVersion int, suffix string) string {
	if patchVersion > 0 {
		return fmt.Sprintf("%.2g.%d%s", version, patchVersion, suffix)
	}
	return fmt.Sprintf("%.2g%s", version, suffix)
}

func quartermasterVersionNoSuffix(version float32, patchVersion int) string {
	if patchVersion > 0 {
		return fmt.Sprintf("%.2g.%d", version, patchVersion)
	}
	return fmt.Sprintf("%.2g", version)
}
