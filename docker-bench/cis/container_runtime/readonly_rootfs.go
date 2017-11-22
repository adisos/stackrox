package containerruntime

import (
	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type readonlyRootfsBenchmark struct{}

func (c *readonlyRootfsBenchmark) Definition() utils.Definition {
	return utils.Definition{
		BenchmarkDefinition: v1.BenchmarkDefinition{
			Name:        "CIS 5.12",
			Description: "Ensure the container's root filesystem is mounted as read only",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *readonlyRootfsBenchmark) Run() (result v1.BenchmarkTestResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		if !container.HostConfig.ReadonlyRootfs {
			utils.Warn(&result)
			utils.AddNotef(&result, "Container %v does not have a readonly rootfs", container.ID)
		}
	}
	return
}

// NewReadonlyRootfsBenchmark implements CIS-5.12
func NewReadonlyRootfsBenchmark() utils.Benchmark {
	return &readonlyRootfsBenchmark{}
}
