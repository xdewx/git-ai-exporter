package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

var Version = "dev"

func doUpdateFn() error {
	if Version == "dev" {
		return fmt.Errorf("version is dev (not built with ldflags). Use --update only on release binaries")
	}
	ver := strings.TrimPrefix(Version, "v")
	v, err := semver.Parse(ver)
	if err != nil {
		return fmt.Errorf("invalid version %q: %w", Version, err)
	}
	latest, err := selfupdate.UpdateSelf(v, "xdewx/git-ai-exporter")
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	if latest.Version.Equals(v) {
		fmt.Fprintln(os.Stderr, "Already up to date ("+Version+")")
	} else {
		fmt.Fprintf(os.Stderr, "Updated to %s\n", latest.Version)
	}
	return nil
}
