package service

import (
	v2ComplianceBenchmark "github.com/stackrox/rox/central/complianceoperator/v2/benchmarks/datastore"
	profileDS "github.com/stackrox/rox/central/complianceoperator/v2/profiles/datastore"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/sync"
)

var (
	serviceInstance     Service
	serviceInstanceInit sync.Once
)

// Singleton returns the singleton instance of the compliance service.
func Singleton() Service {
	if !features.ComplianceEnhancements.Enabled() {
		return nil
	}

	serviceInstanceInit.Do(func() {
		serviceInstance = New(profileDS.Singleton(), v2ComplianceBenchmark.Singleton())
	})
	return serviceInstance
}
