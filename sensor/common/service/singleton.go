package service

import (
	"sync"

	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/sensor/common/cache"
	"github.com/stackrox/rox/sensor/common/processsignal"
)

var (
	once sync.Once

	as Service
)

// newService creates a new streaming service with the collector. It should only be called once.
func newService(containerCache *cache.ContainerCache) Service {
	indicators := make(chan *v1.SensorEvent)

	return &serviceImpl{
		queue:           make(chan *v1.Signal, maxBufferSize),
		indicators:      indicators,
		processPipeline: processsignal.NewProcessPipeline(indicators, containerCache),
	}
}

func initialize() {
	// Creates the signal service with the pending cache embedded
	as = newService(cache.Singleton())
}

// Singleton implements a singleton for the client streaming gRPC service between collector and sensor
func Singleton() Service {
	once.Do(initialize)
	return as
}
