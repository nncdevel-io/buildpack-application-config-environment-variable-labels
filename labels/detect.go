package labels

import (
	"github.com/buildpacks/libcnb/v2"
)

type Detect struct{}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "environment-variable-labels"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "spring-boot"},
					{Name: "environment-variable-labels"},
				},
			},
		},
	}, nil
}
