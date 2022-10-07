package labels

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/buildpacks/libcnb"
	"github.com/mattn/go-zglob"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

type EnvironmentVariable struct {
	Name         string `json:"name"`
	Required     bool   `json:"required"`
	DefaultValue string `json:"defaultValue"`
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	labelKey, _ := cr.Resolve("BP_APP_CONFIG_ENVIRONMENT_VARIABLE_LABEL_NAME")

	targetPatterns, _ := cr.Resolve("BP_APP_CONFIG_ENVIRONMENT_VARIABLE_TARGET_PATTERNS")

	targets := b.parseTargetPatterns(context.Application.Path, targetPatterns)
	candidates := b.findCandidates(targets)

	result := libcnb.NewBuildResult()

	extractor := NewTextPlaceHolderExtractorChain(&b.Logger, candidates)

	environmentVariables, err := extractor.Extract()
	if err != nil {
		panic(err)
	}

	label, err := b.toLabel(labelKey, environmentVariables)
	if err != nil {
		panic(err)
	}

	result.Labels = append(result.Labels, label)

	return result, nil

}

func (b Build) toLabel(labelKey string, environmentVariables []EnvironmentVariable) (libcnb.Label, error) {
	jsonBytes, err := json.Marshal(environmentVariables)

	if err != nil {
		return libcnb.Label{}, err
	}

	label := libcnb.Label{
		Key:   labelKey,
		Value: string(jsonBytes),
	}

	return label, nil
}

func (b Build) parseTargetPatterns(basePath string, targets string) []string {
	splitPattern := regexp.MustCompile(`, *`)
	split := splitPattern.Split(targets, -1)

	var result []string

	for _, v := range split {
		path := filepath.Join(basePath, v)
		result = append(result, path)
	}

	return result
}

func (b Build) findCandidates(targets []string) []string {

	var result []string

	for _, target := range targets {
		matched, err := zglob.Glob(target)

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
