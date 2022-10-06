package labels

import (
	"encoding/json"
	"fmt"

	"github.com/buildpacks/libcnb"
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

	labelKey, _ := cr.Resolve("BP_APPLICATION_CONFIG_ENVIRONMENT_VARIABLE_LABEL_NAME")

	result := libcnb.NewBuildResult()

	extractor := NewTextPlaceHolderExtractorChain(
		&b.Logger,
		[]string{
			"BOOT-INF/classes/application.properties",
			"BOOT-INF/classes/application.yaml",
			"BOOT-INF/classes/application.yml",
			"WEB-INF/classes/application.properties",
			"WEB-INF/classes/application.yaml",
			"WEB-INF/classes/application.yml",
		},
	)

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
