package labels

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/buildpacks/libcnb/v2"
	"github.com/mattn/go-zglob"
	"github.com/paketo-buildpacks/libpak/v2"
	"github.com/paketo-buildpacks/libpak/v2/log"
)

var (
	jsonMarshal           = json.Marshal
	glob                  = zglob.Glob
	configurationResolver = libpak.NewConfigurationResolver
)

type Build struct {
	Logger log.Logger
}

type EnvironmentVariable struct {
	Name         string `json:"name"`
	Required     bool   `json:"required"`
	DefaultValue string `json:"defaultValue"`
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack.Info.Name, context.Buildpack.Info.Version, context.Buildpack.Info.Homepage)

	meta, err := libpak.NewBuildModuleMetadata(context.Buildpack.Metadata)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to build module metadata\n%w", err)
	}

	cr, err := configurationResolver(meta)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	labelKey, _ := cr.Resolve("BP_APP_CONFIG_ENVIRONMENT_VARIABLE_LABEL_NAME")
	targetPatterns, _ := cr.Resolve("BP_APP_CONFIG_ENVIRONMENT_VARIABLE_TARGET_PATTERNS")

	targets := parseTargetPatterns(context.ApplicationPath, targetPatterns)
	candidates := b.findCandidates(targets)

	result := libcnb.NewBuildResult()

	if len(candidates) == 0 {
		return result, nil
	}

	extractor := NewTextPlaceHolderExtractorChain(b.Logger, candidates)

	environmentVariables := extractor.Extract()

	label, err := toLabel(labelKey, environmentVariables)
	if err != nil {
		return result, err
	}

	result.Labels = append(result.Labels, label)

	return result, nil

}

func (b Build) findCandidates(targets []string) []string {

	result := []string{}

	for _, target := range targets {
		matched, err := glob(target)

		b.Logger.Debugf("pattern: %s", target)

		if err != nil {
			b.Logger.Debugf("Glob failed. %s", err)
		}

		result = append(result, matched...)
	}

	b.Logger.Header("Candidate files")

	for _, v := range result {
		b.Logger.Body(v)
	}

	return result

}

func toLabel(labelKey string, environmentVariables []EnvironmentVariable) (libcnb.Label, error) {

	vars := []EnvironmentVariable{}

	if environmentVariables != nil {
		vars = environmentVariables
	}

	jsonBytes, err := jsonMarshal(vars)

	if err != nil {
		return libcnb.Label{
			Key:   labelKey,
			Value: "[]", // empty array
		}, err
	}

	return libcnb.Label{
		Key:   labelKey,
		Value: string(jsonBytes),
	}, nil
}

func parseTargetPatterns(basePath string, targets string) []string {
	split := strings.Split(targets, ",")

	result := []string{}

	for _, v := range split {
		path := filepath.Join(basePath, strings.TrimSpace(v))
		result = append(result, path)
	}

	return result
}
