package labels

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		ctx    libcnb.DetectContext
	)

	it.Before(func() {

		ctx = libcnb.DetectContext{
			Buildpack: libcnb.Buildpack{
				Metadata: map[string]interface{}{},
			},
		}
	})

	it("Detect Always pass", func() {
		detect := Detect{}

		res, err := detect.Detect(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())

		Expect(res.Plans).To(Equal([]libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "environment-variable-labels"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "spring-boot"},
					{Name: "environment-variable-labels"},
				},
			},
		}))

	})
}
