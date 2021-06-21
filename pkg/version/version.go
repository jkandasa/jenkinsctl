package version

import (
	"fmt"
	"runtime"
)

var (
	gitCommit string
	version   string
	buildDate string
)

// Version holds version data
type Version struct {
	Version       string `json:"version"`
	GitCommit     string `json:"gitCommit"`
	BuildDate     string `json:"buildDate"`
	GoLangVersion string `json:"goLangVersion"`
	Platform      string `json:"platform"`
}

// Get returns the Version object
func Get() Version {
	return Version{
		Version:       version,
		BuildDate:     buildDate,
		GitCommit:     gitCommit,
		GoLangVersion: runtime.Version(),
		Platform:      fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns the values as string
func (v Version) String() string {
	return fmt.Sprintf("{version:%s, buildDate:%s, gitCommit:%s, goLangVersion:%s, platform:%s}",
		v.Version, v.BuildDate, v.GitCommit, v.GoLangVersion, v.Platform)
}
