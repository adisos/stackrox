package renderer

import (
	"path"
	"strings"

	"github.com/pkg/errors"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/image"
	"github.com/stackrox/rox/pkg/charts"
	"github.com/stackrox/rox/pkg/helmutil"
	"github.com/stackrox/rox/pkg/zip"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

var (
	commonScriptMap = FileNameMap{
		"common/ca-setup.sh":     "scripts/ca-setup.sh",
		"common/delete-ca.sh":    "scripts/delete-ca.sh",
		"common/port-forward.sh": "scripts/port-forward.sh",
	}

	kubectlScannerScriptMap = FileNameMap{
		"common/setup-scanner.sh": "scanner/scripts/setup.sh",
	}
)

func renderHelmChart(chartFiles []*loader.BufferedFile, mode mode, valuesFiles []*zip.File) ([]*zip.File, error) {
	ch, err := loader.LoadFiles(chartFiles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load central services chart")
	}

	values, err := loadAndMergeValues(valuesFiles)
	if err != nil {
		return nil, errors.Wrap(err, "loading generated values")
	}

	rendered, err := helmutil.Render(ch, values, helmutil.Options{
		ReleaseOptions: chartutil.ReleaseOptions{
			Name:      "stackrox-central-services",
			Namespace: "stackrox",
			IsInstall: mode == renderAll || mode == scannerOnly,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "rendering Helm chart")
	}

	var renderedFiles []*zip.File
	// Filter out non-empty YAML files, and partition them into central and scanner files.
	for fileName, contents := range rendered {
		if path.Ext(fileName) != ".yaml" {
			continue
		}
		contents = strings.TrimSpace(contents)
		if contents == "" {
			continue
		}
		contents += "\n"

		subDir := "central"
		if strings.HasPrefix(path.Base(fileName), "02-scanner-") {
			subDir = "scanner"
		}
		renderedFiles = append(renderedFiles, zip.NewFile(path.Join(subDir, path.Base(fileName)), []byte(contents), 0))
	}

	return renderedFiles, nil
}

func renderNewBasicFiles(c Config, mode mode) ([]*zip.File, error) {
	helmImage := image.GetDefaultImage()
	valuesFiles, err := renderNewHelmValues(c)
	if err != nil {
		return nil, errors.Wrap(err, "rendering new helm values")
	}

	if c.K8sConfig.DeploymentFormat == v1.DeploymentFormat_HELM_VALUES {
		return valuesFiles, nil
	}

	// Helm (full) or kubectl

	chTpl, err := helmImage.GetCentralServicesChartTemplate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain central services chart template")
	}

	metaVals := charts.DefaultMetaValues()
	metaVals["RenderMode"] = mode.String()

	chartFiles, err := chTpl.InstantiateRaw(metaVals)
	if err != nil {
		return nil, errors.Wrap(err, "failed to instantiate central services chart template")
	}

	var helmChartFiles []*zip.File
	helmChartFiles = append(helmChartFiles, valuesFiles...)
	helmChartFiles = append(helmChartFiles, withPrefix("chart", convertBufferedFiles(chartFiles))...)

	if c.K8sConfig.DeploymentFormat == v1.DeploymentFormat_HELM {
		return helmChartFiles, nil
	}

	// kubectl
	if c.K8sConfig.DeploymentFormat != v1.DeploymentFormat_KUBECTL {
		return nil, errors.Errorf("unsupported deployment format %v", c.K8sConfig.DeploymentFormat)
	}

	renderedFiles, err := renderHelmChart(chartFiles, mode, valuesFiles)
	if err != nil {
		return nil, errors.Wrap(err, "rendering helm chart")
	}

	if mode != renderAll {
		return renderedFiles, nil
	}

	// When rendering everything, embed entire helm chart.
	// We need to create a copy of the config to pretend as if the Helm chart
	// was the primary output format.
	k8sCfgForHelm := *c.K8sConfig
	k8sCfgForHelm.DeploymentFormat = v1.DeploymentFormat_HELM
	cfgForHelm := c
	cfgForHelm.K8sConfig = &k8sCfgForHelm

	chartAuxFiles, err := renderAuxiliaryFiles(cfgForHelm, mode)
	if err != nil {
		return nil, errors.Wrap(err, "rendering auxiliary files for Helm chart")
	}
	helmChartFiles = append(helmChartFiles, chartAuxFiles...)

	// Embed helm deployment files in "helm" subdirectory of result
	renderedFiles = append(renderedFiles, withPrefix("helm", helmChartFiles)...)

	return renderedFiles, nil
}

func renderAuxiliaryFiles(c Config, mode mode) ([]*zip.File, error) {
	if mode != renderAll && mode != scannerOnly {
		return nil, nil
	}

	var auxFiles []*zip.File
	readmeFile, err := generateReadmeFile(&c, mode)
	if err != nil {
		return nil, errors.Wrap(err, "generating readme file")
	}
	auxFiles = append(auxFiles, readmeFile)

	assets, err := LoadAssets(assetFileNameMap)
	if err != nil {
		return nil, errors.Wrap(err, "loading asset files")
	}
	if c.K8sConfig.DeploymentFormat == v1.DeploymentFormat_KUBECTL {
		auxFiles = append(auxFiles, withPrefix("scanner/scripts", assets)...)
		if mode == renderAll {
			auxFiles = append(auxFiles, withPrefix("central/scripts", assets)...)
		}
	} else {
		auxFiles = append(auxFiles, withPrefix("scripts", assets)...)
	}

	if c.K8sConfig.DeploymentFormat == v1.DeploymentFormat_KUBECTL {
		scannerScriptFiles, err := RenderFiles(kubectlScannerScriptMap, &c)
		if err != nil {
			return nil, errors.Wrap(err, "rendering scanner script files")
		}
		auxFiles = append(auxFiles, scannerScriptFiles...)
	}

	if mode == renderAll {
		commonScriptFileMap := make(FileNameMap)
		for k, v := range commonScriptMap {
			commonScriptFileMap[k] = v
		}
		if c.K8sConfig.DeploymentFormat == v1.DeploymentFormat_KUBECTL {
			commonScriptFileMap["common/setup-central.sh"] = "scripts/setup.sh"
		} else {
			commonScriptFileMap["common/setup.sh"] = "scripts/setup.sh"
		}

		commonScriptFiles, err := RenderFiles(commonScriptFileMap, &c)
		if err != nil {
			return nil, errors.Wrap(err, "rendering common script files")
		}

		if c.K8sConfig.DeploymentFormat == v1.DeploymentFormat_KUBECTL {
			auxFiles = append(auxFiles, withPrefix("central", commonScriptFiles)...)
		} else {
			auxFiles = append(auxFiles, commonScriptFiles...)
		}
	}

	return auxFiles, nil
}

func renderNew(c Config, mode mode) ([]*zip.File, error) {
	var allFiles []*zip.File

	basicFiles, err := renderNewBasicFiles(c, mode)
	if err != nil {
		return nil, errors.Wrap(err, "rendering basic output files")
	}
	allFiles = append(allFiles, basicFiles...)

	auxFiles, err := renderAuxiliaryFiles(c, mode)
	if err != nil {
		return nil, errors.Wrap(err, "rendering auxiliary output files")
	}
	allFiles = append(allFiles, auxFiles...)

	return allFiles, nil
}
