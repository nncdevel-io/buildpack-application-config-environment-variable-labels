package main

import (
	"os"

	"github.com/nncdevel-io/buildpack-application-config-environment-variable-labels/labels"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

func main() {
	libpak.Main(
		labels.Detect{},
		labels.Build{Logger: bard.NewLogger(os.Stdout)},
	)
}
