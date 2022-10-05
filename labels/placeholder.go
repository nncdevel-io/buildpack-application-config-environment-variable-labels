package labels

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/paketo-buildpacks/libpak/bard"
)

type PlaceholderExtractor interface {
	Extract() ([]EnvironmentVariable, error)
}

type PlaceholderExtractorChain struct {
	Logger     *bard.Logger
	Extractors []PlaceholderExtractor
}

func NewTextPlaceHolderExtractorChain(logger *bard.Logger, paths []string) *PlaceholderExtractorChain {
	logger.Header("Candidate files")
	var extractors []PlaceholderExtractor
	for _, v := range paths {
		logger.Body(v)
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

func (p PlaceholderExtractorChain) Extract() ([]EnvironmentVariable, error) {
	for _, extractor := range p.Extractors {
		r, err := extractor.Extract()
		if len(r) == 0 || err != nil {
			continue
		}
		return r, err
	}
	return []EnvironmentVariable{}, nil
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

func (p TextPlaceHolderExtractor) Extract() ([]EnvironmentVariable, error) {

	p.Logger.Headerf("Target file : %s\n", p.TargetFilePath)

	if p.targetFileIsNotExists() {
		p.Logger.Body("File does not exists.\n")
		return []EnvironmentVariable{}, nil
	}

	propertiesBytes, err := ioutil.ReadFile(p.TargetFilePath)
	if err != nil {
		p.Logger.Bodyf("Read target file failed. %v", err)
		return []EnvironmentVariable{}, nil
	}

	environmentVariables := []EnvironmentVariable{}

	placeholderRegExp := regexp.MustCompile(`(\$\{([^}]+)})`)
	matched := placeholderRegExp.FindAllStringSubmatch(string(propertiesBytes), -1)

	for _, v := range matched {
		inset := v[2]
		variable := ParsePlaceholder(inset)

		p.Logger.Bodyf(`EnvironmentVariable: "%s" DefaultValue: "%s"`, variable.Name, variable.DefaultValue)
		environmentVariables = append(environmentVariables, variable)
	}

	return environmentVariables, nil

}

func (p TextPlaceHolderExtractor) targetFileIsNotExists() bool {
	_, err := os.Stat(p.TargetFilePath)
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
