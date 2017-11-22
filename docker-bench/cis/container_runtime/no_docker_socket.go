package containerruntime

import (
	"strings"

	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type dockerSocketMountBenchmark struct{}

func (c *dockerSocketMountBenchmark) Definition() utils.Definition {
	return utils.Definition{
		BenchmarkDefinition: v1.BenchmarkDefinition{
			Name:        "CIS 5.31",
			Description: "Ensure the Docker socket is not mounted inside any containers",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *dockerSocketMountBenchmark) Run() (result v1.BenchmarkTestResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		for _, containerMount := range container.Mounts {
			if strings.Contains(containerMount.Source, "docker.sock") {
				utils.Warn(&result)
				utils.AddNotef(&result, "Container %v has mounted docker.sock", container.ID)
			}
		}
	}
	return
}

// NewDockerSocketMountBenchmark implements CIS-5.31
func NewDockerSocketMountBenchmark() utils.Benchmark {
	return &dockerSocketMountBenchmark{}
}
