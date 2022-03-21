package version

import (
	"embed"
	"fmt"
	"runtime/debug"
	"strings"
)

var (
	// Version semver for syncrepos (set through go build -ldflags)
	Version string
	// BuildUser is the user login who initiated the build (set through go build -ldflags)
	BuildUser string
	// GitTag is the result of git describe (set through go build -ldflags)
	GitTag string
	// BuildDate is the date of build (set through go build -ldflags)
	BuildDate string
)

// String displays all the version values
func String(verlevel int, versionFS embed.FS) string {
	res := ""
	// https://github.com/golang/go/issues/41191
	// https://stackoverflow.com/a/67357103/6309
	v, err := versionFS.ReadFile("version/version.txt")
	if err != nil {
		res = res + fmt.Sprintf("Unknown version: %+v\n", err)
	} else {
		res = res + strings.TrimSpace(string(v))
	}
	if verlevel >= 2 {
		res = res + "\n"
		// https://www.reddit.com/r/golang/comments/rxfs5i/go_118_debuginfo_why_not_include_the_current_git/
		info, ok := debug.ReadBuildInfo()
		if !ok {
			res = res + "Build info not found\n"
		} else {
			dbs := info.Settings
			vcs := get(dbs, "vcs")
			rev := get(dbs, "vcs.revision")
			dirty := ""
			if get(dbs, "vcs.modified") == "true" {
				dirty = " (dirty)"
			}
			date := get(dbs, "vcs.time")
			sum := ""
			if info.Main.Sum != "" {
				sum = " (" + info.Main.Sum + ")"
			}
			res = res + fmt.Sprintf("mod '%s'%s, %s", info.Main.Path, sum, info.GoVersion)
			if vcs != "" && rev != "" && date != "" {
				res = res + fmt.Sprintf("\nVCS %s revision %s%s, %s", vcs, rev, dirty, date)
			}
			/*
				op, err := json.MarshalIndent(info.Settings, "", " ")
				if err != nil {
					res = res + fmt.Sprintf("buildinfo error marshalling: %+v\n", err)
				}
				res = res + string(op)
			*/
		}
		//spew.Dump(info)
	}
	if verlevel >= 3 {
		if GitTag != "" && BuildUser != "" && Version != "" && BuildDate != "" {
			res = res + "\n"
			res = res + fmt.Sprintf("Git Tag   : %s\n", GitTag)
			res = res + fmt.Sprintf("Build User: %s\n", BuildUser)
			res = res + fmt.Sprintf("Version   : %s\n", Version)
			res = res + fmt.Sprintf("BuildDate : %s\n", BuildDate)
		}
	}
	return res
}

func get(dbs []debug.BuildSetting, key string) string {
	for _, bs := range dbs {
		if bs.Key == key {
			return bs.Value
		}
	}
	return ""
}
