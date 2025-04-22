package main

import (
	"os"

	"github.com/nncdevel-io/buildpack-application-config-environment-variable-labels/labels"
	"github.com/paketo-buildpacks/libpak/v2"
	"github.com/paketo-buildpacks/libpak/v2/log"
)

func main() {
	detector := labels.Detect{}
	builder := labels.Build{Logger: log.NewPaketoLogger(os.Stdout)}
	libpak.BuildpackMain(
		detector.Detect,
		builder.Build,
	)
}
