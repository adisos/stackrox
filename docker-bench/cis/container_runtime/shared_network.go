package containerruntime

import (
	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type sharedNetworkBenchmark struct{}

func (c *sharedNetworkBenchmark) Definition() utils.Definition {
	return utils.Definition{
		BenchmarkDefinition: v1.BenchmarkDefinition{
			Name:        "CIS 5.9",
			Description: "Ensure the host's network namespace is not shared",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *sharedNetworkBenchmark) Run() (result v1.BenchmarkTestResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		if container.HostConfig.NetworkMode.IsHost() {
			utils.Warn(&result)
			utils.AddNotef(&result, "Container %v has network set to --net=host", container.ID)
		}
	}
	return
}

// NewSharedNetworkBenchmark implements CIS-5.9
func NewSharedNetworkBenchmark() utils.Benchmark {
	return &sharedNetworkBenchmark{}
}
