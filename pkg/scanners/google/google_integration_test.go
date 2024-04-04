//go:build integration

package google

import (
	"os"
	"testing"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/images/types"
	"github.com/stackrox/rox/pkg/images/utils"
	"github.com/stackrox/rox/pkg/registries/google"
	"github.com/stretchr/testify/require"
)

const project = "acs-san-stackroxci"

func TestGoogle(t *testing.T) {
	serviceAccount := os.Getenv("SERVICE_ACCOUNT")
	if serviceAccount == "" {
		t.Skip("SERVICE_ACCOUNT is required for Google integration test")
		return
	}
	t.Setenv("ROX_REGISTRY_RESPONSE_TIMEOUT", "90s")
	t.Setenv("ROX_REGISTRY_CLIENT_TIMEOUT", "120s")

	integration := &storage.ImageIntegration{
		IntegrationConfig: &storage.ImageIntegration_Google{
			Google: &storage.GoogleConfig{
				Endpoint:       "us.gcr.io",
				ServiceAccount: os.Getenv("SERVICE_ACCOUNT"),
				Project:        project,
			},
		},
	}

	_, creator := google.Creator()

	registry, err := creator(integration)
	require.NoError(t, err)

	scanner, err := newScanner(integration)
	require.NoError(t, err)

	var images = []string{
		"us.gcr.io/acs-san-stackroxci/music-nginx:latest",
		"us.gcr.io/acs-san-stackroxci/nginx:slim",
		"us.gcr.io/acs-san-stackroxci/ubuntu:latest",
	}

	for _, i := range images {
		containerImage, err := utils.GenerateImageFromString(i)
		require.NoError(t, err)

		img := types.ToImage(containerImage)
		metadata, err := registry.Metadata(img)
		require.NoError(t, err)
		img.Metadata = metadata
		img.Id = utils.GetSHA(img)

		scan, err := scanner.GetScan(img)
		require.NoError(t, err)
		require.NotEmpty(t, scan.GetComponents())
		for _, c := range scan.GetComponents() {
			for _, v := range c.Vulns {
				require.NotEmpty(t, v.Cve)
			}
		}
	}
}
