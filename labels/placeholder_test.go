package labels

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak/v2/log"
	"github.com/sclevine/spec"
)

func createTempFile() *os.File {

	tmp, err := os.CreateTemp("/tmp", "NewTextPlaceHolderExtractor")

	if err != nil {
		panic(err)
	}

	_, err = tmp.WriteString("a=${placeholder_a}")

	if err != nil {
		panic(err)
	}

	return tmp
}

func testPlaceholder(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
	)

	when("NewTextPlaceHolderExtractorChain", func() {
		tmp := createTempFile()
		it("returns new instance", func() {
			targets := []string{"/tmp/dummy", tmp.Name()}
			extractors := NewTextPlaceHolderExtractorChain(log.NewDiscardLogger(), targets)

			envs := extractors.Extract()

			Expect(envs).Should(HaveLen(1))
		})
	})

	when("NewPlaceholderExtractorChain", func() {
		it("returns new instance", func() {
			extractor := NewPlaceholderExtractorChain(log.NewDiscardLogger(), []PlaceholderExtractor{})

			Expect(extractor.Extractors).Should(HaveLen(0))
		})

		it("Extract with no extractor", func() {
			extractor := NewPlaceholderExtractorChain(log.NewDiscardLogger(), []PlaceholderExtractor{})

			envs := extractor.Extract()

			Expect(envs).Should(BeEmpty())

		})

		it("Extract with single extractor", func() {

			tmp := createTempFile()

			extractors := []PlaceholderExtractor{
				NewTextPlaceHolderExtractor(log.NewDiscardLogger(), tmp.Name()),
			}

			extractor := NewPlaceholderExtractorChain(log.NewDiscardLogger(), extractors)

			envs := extractor.Extract()

			Expect(envs).Should(HaveLen(1))

		})

		it("Extract with multiple extractor", func() {

			tmp := createTempFile()

			extractors := []PlaceholderExtractor{
				NewTextPlaceHolderExtractor(log.NewDiscardLogger(), "/tmp/dummy"),
				NewTextPlaceHolderExtractor(log.NewDiscardLogger(), tmp.Name()),
			}

			extractor := NewPlaceholderExtractorChain(log.NewDiscardLogger(), extractors)

			envs := extractor.Extract()

			Expect(envs).Should(HaveLen(1))

		})
	})

	when("TextPlaceHolderExtractor", func() {
		it("NewTextPlaceHolderExtractor", func() {
			tmp, err := os.CreateTemp("/tmp", "NewTextPlaceHolderExtractor")

			if err != nil {
				panic(err)
			}

			extractor := NewTextPlaceHolderExtractor(log.NewDiscardLogger(), tmp.Name())

			Expect(extractor).ShouldNot(BeNil())
		})

		it("can extract placeholder in file", func() {
			tmp := createTempFile()

			extractor := NewTextPlaceHolderExtractor(log.NewDiscardLogger(), tmp.Name())

			envs := extractor.Extract()

			Expect(envs).Should(HaveLen(1))
		})

		it("could not extract placeholder when target file not exists", func() {
			extractor := NewTextPlaceHolderExtractor(log.NewDiscardLogger(), "/tmp/not-exists")

			envs := extractor.Extract()

			Expect(envs).Should(HaveLen(0))
		})

		it("read file failed", func() {

			tmp, err := os.CreateTemp("/tmp", "NewTextPlaceHolderExtractor")

			if err != nil {
				panic(err)
			}

			t.Logf("temporary name: %s\n", tmp.Name())

			err = os.Chmod(tmp.Name(), 0000)
			if err != nil {
				panic(err)
			}

			extractor := NewTextPlaceHolderExtractor(log.NewDiscardLogger(), tmp.Name())

			envs := extractor.Extract()

			Expect(envs).Should(HaveLen(0))

			err = os.Chmod(tmp.Name(), 0600)
			if err != nil {
				panic(err)
			}
			err = os.Remove(tmp.Name())
			if err != nil {
				return
			}
		})

	})

	when("extractEnvironmentVariablePlaceholders", func() {
		it("no placeholder", func() {
			res := extractEnvironmentVariablePlaceholders("test", log.NewDiscardLogger())

			Expect(res).Should(BeEmpty())
		})

		it("simple", func() {
			input := `
			a=${placeholder_a}
			`

			res := extractEnvironmentVariablePlaceholders(input, log.NewDiscardLogger())

			expect := EnvironmentVariable{
				Name:         "placeholder_a",
				Required:     true,
				DefaultValue: "",
			}

			Expect(res).Should(HaveLen(1))
			Expect(res[0]).Should(Equal(expect))
		})

		it("duplicated", func() {
			input := `
			a=${placeholder_1}
			b=${placeholder_1:default}
			`

			res := extractEnvironmentVariablePlaceholders(input, log.NewDiscardLogger())

			expect := EnvironmentVariable{
				Name:         "placeholder_1",
				Required:     true,
				DefaultValue: "",
			}

			Expect(res).Should(HaveLen(1))
			Expect(res[0]).Should(Equal(expect))
		})

		it("multi placeholder in single line", func() {
			input := `
			a=${placeholder_a}_${placeholder_b}
			`

			res := extractEnvironmentVariablePlaceholders(input, log.NewDiscardLogger())

			Expect(res).Should(HaveLen(2))
			Expect(res[0].Name).Should(Equal("placeholder_a"))
			Expect(res[1].Name).Should(Equal("placeholder_b"))

		})

		it("multiline", func() {
			input := `
			a=${placeholder_a}
			b=${placeholder_b}
			`

			res := extractEnvironmentVariablePlaceholders(input, log.NewDiscardLogger())

			Expect(res).Should(HaveLen(2))
			Expect(res[0].Name).Should(Equal("placeholder_a"))
			Expect(res[1].Name).Should(Equal("placeholder_b"))

		})

		it("environment variable order", func() {
			input := `
			a=${placeholder_b}
			b=${placeholder_a}
			`

			res := extractEnvironmentVariablePlaceholders(input, log.NewDiscardLogger())

			Expect(res).Should(HaveLen(2))
			Expect(res[0].Name).Should(Equal("placeholder_a"))
			Expect(res[1].Name).Should(Equal("placeholder_b"))

		})
	})

	when("targetFileIsNotExists", func() {
		it("target file not exists.", func() {

			res := targetFileIsNotExists("/tmp/not-exists")

			Expect(res).To(Equal(true))

		})

		it("target file not exists.", func() {
			target, err := os.CreateTemp("/tmp", "targetFileIsNotExists")

			if err != nil {
				panic(err)
			}

			res := targetFileIsNotExists(target.Name())

			Expect(res).To(Equal(false))
		})
	})

	when("ParsePlaceholder", func() {

		it("no default value", func() {

			env := ParsePlaceholder("abc")

			expected := EnvironmentVariable{
				Name:         "abc",
				Required:     true,
				DefaultValue: "",
			}

			Expect(env).Should(Equal(expected))

		})

		it("with empty defalt value", func() {

			env := ParsePlaceholder("abc:")

			expected := EnvironmentVariable{
				Name:         "abc",
				Required:     false,
				DefaultValue: "",
			}

			Expect(env).Should(Equal(expected))

		})

		it("with defalt value", func() {

			env := ParsePlaceholder("abc:default")

			expected := EnvironmentVariable{
				Name:         "abc",
				Required:     false,
				DefaultValue: "default",
			}

			Expect(env).Should(Equal(expected))

		})

	})
}
