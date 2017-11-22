package containerruntime

import (
	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type ipcNamespaceBenchmark struct{}

func (c *ipcNamespaceBenchmark) Definition() utils.Definition {
	return utils.Definition{
		BenchmarkDefinition: v1.BenchmarkDefinition{
			Name:        "CIS 5.16",
			Description: "Ensure the host's IPC namespace is not shared",
		}, Dependencies: []utils.Dependency{utils.InitContainers},
	}
}

func (c *ipcNamespaceBenchmark) Run() (result v1.BenchmarkTestResult) {
	utils.Pass(&result)
	for _, container := range utils.ContainersRunning {
		if container.HostConfig.IpcMode.IsHost() {
			utils.Warn(&result)
			utils.AddNotef(&result, "Container %v has ipc mode set to host", container.ID)
		}
	}
	return
}

// NewIpcNamespaceBenchmark implements CIS-5.16
func NewIpcNamespaceBenchmark() utils.Benchmark {
	return &ipcNamespaceBenchmark{}
}
