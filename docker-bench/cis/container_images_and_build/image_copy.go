package containerimagesandbuild

import (
	"context"
	"strings"

	"bitbucket.org/stack-rox/apollo/docker-bench/utils"
	"bitbucket.org/stack-rox/apollo/pkg/api/generated/api/v1"
)

type imageCopyBenchmark struct{}

func (c *imageCopyBenchmark) Definition() utils.Definition {
	return utils.Definition{
		BenchmarkDefinition: v1.BenchmarkDefinition{
			Name:        "CIS 4.9",
			Description: "Ensure COPY is used instead of ADD in Dockerfile",
		}, Dependencies: []utils.Dependency{utils.InitImages},
	}
}

func (c *imageCopyBenchmark) Run() (result v1.BenchmarkTestResult) {
	utils.Pass(&result)
	for _, image := range utils.Images {
		historySlice, err := utils.DockerClient.ImageHistory(context.Background(), image.ID)
		if err != nil {
			utils.Warn(&result)
			utils.AddNotef(&result, "Could not get image history for image %v: %+v", err)
			continue
		}
		for _, history := range historySlice {
			cmd := strings.ToLower(history.CreatedBy)
			if strings.Contains(cmd, "add file:") || strings.Contains(cmd, "add dir:") {
				utils.Warn(&result)
				utils.AddNotef(&result, "Image %v has an ADD instead of a COPY command", utils.GetReadableImageName(image))
				break
			}
		}
	}
	return
}

// NewImageCopyBenchmark implements CIS-4.9
func NewImageCopyBenchmark() utils.Benchmark {
	return &imageCopyBenchmark{}
}
