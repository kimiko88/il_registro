package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the application semantic version string.
	// Can be overridden at build time via -ldflags "-X registro-backend/pkg/version.Version=X.Y.Z".
	Version = "1.0.0-beta"

	// GitCommit is the short commit hash of the build.
	// Can be overridden at build time via -ldflags "-X registro-backend/pkg/version.GitCommit=abc1234".
	GitCommit = "dev"

	// BuildDate is the ISO timestamp when the binary was built.
	// Can be overridden at build time via -ldflags "-X registro-backend/pkg/version.BuildDate=2026-09-11T12:00:00Z".
	BuildDate = "unknown"
)

// Info contains complete version and runtime build metadata.
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_commit,omitempty"`
	BuildDate string `json:"build_date,omitempty"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

// Get returns the structured version information.
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
}

// String returns a human-readable representation of the version.
func String() string {
	if GitCommit != "dev" && GitCommit != "" {
		return fmt.Sprintf("v%s (%s, %s)", Version, GitCommit, BuildDate)
	}
	return "v" + Version
}
