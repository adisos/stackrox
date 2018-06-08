package main

import (
	"fmt"

	"bitbucket.org/stack-rox/apollo/generated/api/v1"
	"bitbucket.org/stack-rox/apollo/pkg/central"
)

var (
	cfg central.Config
)

func validateConfig(c central.Config, cluster v1.ClusterType) error {
	if err := validateExternal(c.External, cluster); err != nil {
		return err
	}
	return validateHostPath(c.HostPath)
}

func validateHostPath(hostpath *central.HostPathPersistence) error {
	if hostpath == nil {
		return nil
	}
	if hostpath.NodeSelectorKey == "" || hostpath.NodeSelectorValue == "" {
		return fmt.Errorf("Both node selector key and node selector value must be specified when using a hostpath")
	}
	return nil
}

func validateExternal(ext *central.ExternalPersistence, cluster v1.ClusterType) error {
	if ext == nil {
		return nil
	}
	if cluster == v1.ClusterType_SWARM_CLUSTER && ext.Name == "" {
		return fmt.Errorf("name must be specified for external volume in Swarm")
	}
	return nil
}
