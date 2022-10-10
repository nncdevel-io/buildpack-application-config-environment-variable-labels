package labels

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		ctx    libcnb.BuildContext
	)

	it.Before(func() {

		ctx = libcnb.BuildContext{
			Buildpack: libcnb.Buildpack{
				Metadata: map[string]interface{}{
					"configurations": []map[string]interface{}{
						{
							"build":   true,
							"default": "test-key",
							"name":    "BP_APP_CONFIG_ENVIRONMENT_VARIABLE_LABEL_NAME",
						},
						{
							"build":   true,
							"default": "test-1, test-2",
							"name":    "BP_APP_CONFIG_ENVIRONMENT_VARIABLE_TARGET_PATTERNS",
						},
					},
				},
				Info: libcnb.BuildpackInfo{
					Name: "test",
				},
			},
		}

		ctx.Application.Path = "/workspace"
	})

	when("Build", func() {

		it("no candidates", func() {
			build := Build{
				Logger: bard.Logger{},
			}

			// build.Build(ctx)
			buildResult, err := build.Build(ctx)

			fmt.Printf("buildResult: %v\n", buildResult)
			fmt.Printf("err: %v\n", err)

			Expect(buildResult).ShouldNot(BeNil())
			Expect(err).Should(BeNil())

		})

		it("configuration resolver raise error", func() {
			originalConfigurationResolver := configurationResolver

			configurationResolver = func(buildpack libcnb.Buildpack, logger *bard.Logger) (libpak.ConfigurationResolver, error) {
				return libpak.ConfigurationResolver{}, fmt.Errorf("Dummy Error")
			}

			build := Build{
				Logger: bard.Logger{},
			}

			buildResult, err := build.Build(ctx)

			Expect(buildResult).Should(Equal(libcnb.BuildResult{}))
			Expect(err).To(HaveOccurred())

			configurationResolver = originalConfigurationResolver
		})

		it("success", func() {
			tmp := createTempFile()
			ctx = libcnb.BuildContext{
				Buildpack: libcnb.Buildpack{
					Metadata: map[string]interface{}{
						"configurations": []map[string]interface{}{
							{
								"build":   true,
								"default": "test-key",
								"name":    "BP_APP_CONFIG_ENVIRONMENT_VARIABLE_LABEL_NAME",
							},
							{
								"build":   true,
								"default": filepath.Base(tmp.Name()),
								"name":    "BP_APP_CONFIG_ENVIRONMENT_VARIABLE_TARGET_PATTERNS",
							},
						},
					},
					Info: libcnb.BuildpackInfo{
						Name: "test",
					},
				},
			}

			ctx.Application.Path = "/tmp"

			build := Build{
				Logger: bard.Logger{},
			}

			buildResult, err := build.Build(ctx)

			fmt.Printf("buildResult: %v\n", buildResult)
			fmt.Printf("err: %v\n", err)

			Expect(buildResult).ShouldNot(BeNil())
			Expect(err).Should(BeNil())

		})

		it("toLabel failed", func() {
			tmp := createTempFile()
			ctx = libcnb.BuildContext{
				Buildpack: libcnb.Buildpack{
					Metadata: map[string]interface{}{
						"configurations": []map[string]interface{}{
							{
								"build":   true,
								"default": "test-key",
								"name":    "BP_APP_CONFIG_ENVIRONMENT_VARIABLE_LABEL_NAME",
							},
							{
								"build":   true,
								"default": filepath.Base(tmp.Name()),
								"name":    "BP_APP_CONFIG_ENVIRONMENT_VARIABLE_TARGET_PATTERNS",
							},
						},
					},
					Info: libcnb.BuildpackInfo{
						Name: "test",
					},
				},
			}

			ctx.Application.Path = "/tmp"

			// Backup and defer recover
			oldJsonMarshal := jsonMarshal
			defer func() {
				jsonMarshal = oldJsonMarshal
			}()

			// Mock
			jsonMarshal = func(v interface{}) ([]byte, error) {
				return nil, errors.New("forced error to marshal")
			}

			build := Build{
				Logger: bard.Logger{},
			}

			buildResult, err := build.Build(ctx)

			fmt.Printf("buildResult: %v\n", buildResult)
			fmt.Printf("err: %v\n", err)

			Expect(buildResult).ShouldNot(BeNil())
			Expect(err).To(HaveOccurred())

			jsonMarshal = oldJsonMarshal
		})

	})

	when("toLabel", func() {
		it("variable is nil", func() {

			var vars []EnvironmentVariable

			key := "test-key"

			res, err := toLabel(key, vars)

			expected := libcnb.Label{
				Key:   "test-key",
				Value: `[]`,
			}

			Expect(err).Should(BeNil())
			Expect(res).Should(Equal(expected))

		})

		it("variable is empty", func() {

			var vars []EnvironmentVariable

			key := "test-key"

			res, err := toLabel(key, vars)

			expected := libcnb.Label{
				Key:   "test-key",
				Value: `[]`,
			}

			Expect(err).Should(BeNil())
			Expect(res).Should(Equal(expected))

		})

		it("json.Marshal returns error.", func() {

			// Backup and defer recover
			oldJsonMarshal := jsonMarshal
			defer func() {
				jsonMarshal = oldJsonMarshal
			}()

			// Mock
			jsonMarshal = func(v interface{}) ([]byte, error) {
				return nil, errors.New("forced error to marshal")
			}

			vars := []EnvironmentVariable{}

			key := "test-key"

			res, err := toLabel(key, vars)

			expected := libcnb.Label{
				Key:   "test-key",
				Value: `[]`,
			}

			Expect(err).ShouldNot(BeNil())
			Expect(res).Should(Equal(expected))

			jsonMarshal = oldJsonMarshal

		})
	}, spec.Random())

	when("parseTargetPatterns", func() {

		it("split with comma & space", func() {

			res := parseTargetPatterns("base", "a, b, c")

			Expect(res).Should(HaveLen(3))

			Expect(res).Should(ConsistOf("base/a", "base/b", "base/c"))

		})

		it("split with only comma", func() {

			res := parseTargetPatterns("base", "a,b,c")

			Expect(res).Should(HaveLen(3))

			Expect(res).Should(ConsistOf("base/a", "base/b", "base/c"))

		})

		it("split with mixed", func() {

			res := parseTargetPatterns("base", "a,b, c")

			Expect(res).Should(HaveLen(3))

			Expect(res).Should(ConsistOf("base/a", "base/b", "base/c"))

		})

	}, spec.Random())

}
