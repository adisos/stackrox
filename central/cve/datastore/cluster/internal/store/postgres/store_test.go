// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
)

type ClusterCvesStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
	store       Store
	pool        *pgxpool.Pool
}

func TestClusterCvesStore(t *testing.T) {
	suite.Run(t, new(ClusterCvesStoreSuite))
}

func (s *ClusterCvesStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}

	ctx := sac.WithAllAccess(context.Background())

	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	s.Require().NoError(err)
	pool, err := pgxpool.ConnectConfig(ctx, config)
	s.Require().NoError(err)

	Destroy(ctx, pool)

	s.pool = pool
	s.store = New(ctx, pool)
}

func (s *ClusterCvesStoreSuite) TearDownTest() {
	if s.pool != nil {
		s.pool.Close()
	}
	s.envIsolator.RestoreAll()
}

func (s *ClusterCvesStoreSuite) TestStore() {
	ctx := sac.WithAllAccess(context.Background())

	store := s.store

	cVE := &storage.CVE{}
	s.NoError(testutils.FullInit(cVE, testutils.SimpleInitializer(), testutils.JSONFieldsFilter))

	foundCVE, exists, err := store.Get(ctx, cVE.GetId(), cVE.GetCve(), cVE.GetOperatingSystem())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundCVE)

	s.NoError(store.Upsert(ctx, cVE))
	foundCVE, exists, err = store.Get(ctx, cVE.GetId(), cVE.GetCve(), cVE.GetOperatingSystem())
	s.NoError(err)
	s.True(exists)
	s.Equal(cVE, foundCVE)

	cVECount, err := store.Count(ctx)
	s.NoError(err)
	s.Equal(1, cVECount)

	cVEExists, err := store.Exists(ctx, cVE.GetId(), cVE.GetCve(), cVE.GetOperatingSystem())
	s.NoError(err)
	s.True(cVEExists)
	s.NoError(store.Upsert(ctx, cVE))

	foundCVE, exists, err = store.Get(ctx, cVE.GetId(), cVE.GetCve(), cVE.GetOperatingSystem())
	s.NoError(err)
	s.True(exists)
	s.Equal(cVE, foundCVE)

	s.NoError(store.Delete(ctx, cVE.GetId(), cVE.GetCve(), cVE.GetOperatingSystem()))
	foundCVE, exists, err = store.Get(ctx, cVE.GetId(), cVE.GetCve(), cVE.GetOperatingSystem())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundCVE)

	var cVEs []*storage.CVE
	for i := 0; i < 200; i++ {
		cVE := &storage.CVE{}
		s.NoError(testutils.FullInit(cVE, testutils.UniqueInitializer(), testutils.JSONFieldsFilter))
		cVEs = append(cVEs, cVE)
	}

	s.NoError(store.UpsertMany(ctx, cVEs))

	cVECount, err = store.Count(ctx)
	s.NoError(err)
	s.Equal(200, cVECount)
}
