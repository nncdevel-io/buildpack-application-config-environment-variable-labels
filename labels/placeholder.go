package labels

import (
	"os"
	"regexp"
	"strings"

	"github.com/paketo-buildpacks/libpak/bard"
)

type PlaceholderExtractor interface {
	Extract() []EnvironmentVariable
}

type PlaceholderExtractorChain struct {
	Logger     *bard.Logger
	Extractors []PlaceholderExtractor
}

func NewTextPlaceHolderExtractorChain(logger *bard.Logger, paths []string) *PlaceholderExtractorChain {
	extractors := []PlaceholderExtractor{}
	for _, v := range paths {
		extractors = append(extractors, NewTextPlaceHolderExtractor(logger, v))
	}
	return NewPlaceholderExtractorChain(logger, extractors)
}

func NewPlaceholderExtractorChain(logger *bard.Logger, extractors []PlaceholderExtractor) *PlaceholderExtractorChain {
	return &PlaceholderExtractorChain{
		Logger:     logger,
		Extractors: extractors,
	}
}

func (p PlaceholderExtractorChain) Extract() []EnvironmentVariable {
	for _, extractor := range p.Extractors {
		r := extractor.Extract()
		if len(r) == 0 {
			continue
		}
		return r
	}
	return []EnvironmentVariable{}
}

type TextPlaceHolderExtractor struct {
	Logger         *bard.Logger
	TargetFilePath string
}

func NewTextPlaceHolderExtractor(logger *bard.Logger, path string) *TextPlaceHolderExtractor {
	return &TextPlaceHolderExtractor{
		Logger:         logger,
		TargetFilePath: path,
	}
}

func (p TextPlaceHolderExtractor) Extract() []EnvironmentVariable {

	p.Logger.Headerf("Target file : %s\n", p.TargetFilePath)

	environmentVariables := []EnvironmentVariable{}

	if targetFileIsNotExists(p.TargetFilePath) {
		p.Logger.Body("File does not exists.\n")
		return environmentVariables
	}

	propertiesBytes, err := os.ReadFile(p.TargetFilePath)
	if err != nil {
		p.Logger.Bodyf("Read target file failed. %v", err)
		return environmentVariables
	}

	environmentVariables = extractEnvironmentVariablePlaceholders(string(propertiesBytes), p.Logger)

	return environmentVariables

}

func extractEnvironmentVariablePlaceholders(input string, logger *bard.Logger) []EnvironmentVariable {
	environmentVariables := []EnvironmentVariable{}

	set := map[string]EnvironmentVariable{}

	placeholderRegExp := regexp.MustCompile(`(\$\{([^}]+)})`)
	matched := placeholderRegExp.FindAllStringSubmatch(input, -1)

	for _, v := range matched {
		inset := v[2]
		variable := ParsePlaceholder(inset)

		logger.Bodyf(`EnvironmentVariable: "%s" DefaultValue: "%s"`, variable.Name, variable.DefaultValue)

		if _, ok := set[variable.Name]; !ok {
			set[variable.Name] = variable
		}
	}

	for _, variable := range set {
		environmentVariables = append(environmentVariables, variable)
	}

	return environmentVariables
}

func targetFileIsNotExists(targetFilePath string) bool {
	_, err := os.Stat(targetFilePath)
	return err != nil && os.IsNotExist(err)
}

func ParsePlaceholder(placeholder string) EnvironmentVariable {
	split := strings.SplitN(placeholder, ":", 2)

	defaultValue := ""
	required := len(split) == 1

	if !required {
		placeholder = split[0]
		defaultValue = split[1]
	}

	return EnvironmentVariable{
		Name:         placeholder,
		Required:     required,
		DefaultValue: defaultValue,
	}
}
