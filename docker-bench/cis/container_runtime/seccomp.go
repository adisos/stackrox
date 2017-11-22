package containerruntime

import (
	"strings"

	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type seccompBenchmark struct{}

func (c *seccompBenchmark) Definition() utils.Definition {
	return utils.Definition{
		BenchmarkDefinition: v1.BenchmarkDefinition{
			Name:        "CIS 5.21",
			Description: "Ensure the default seccomp profile is not Disabled",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *seccompBenchmark) Run() (result v1.BenchmarkTestResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		for _, opt := range container.HostConfig.SecurityOpt {
			if strings.Contains(opt, "seccomp:unconfined") {
				utils.Warn(&result)
				utils.AddNotef(&result, "Container %v has seccomp set to unconfined", container.ID)
				break
			}
		}
	}
	return
}

// NewSeccompBenchmark implements CIS-5.21
func NewSeccompBenchmark() utils.Benchmark {
	return &seccompBenchmark{}
}
