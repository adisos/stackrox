// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/sac/resources"
	"github.com/stackrox/rox/pkg/search"
	"github.com/stackrox/rox/pkg/search/postgres/mapping"
)

var (
	// CreateTableCloudSourcesStmt holds the create statement for table `cloud_sources`.
	CreateTableCloudSourcesStmt = &postgres.CreateStmts{
		GormModel: (*CloudSources)(nil),
		Children:  []*postgres.CreateStmts{},
	}

	// CloudSourcesSchema is the go schema for table `cloud_sources`.
	CloudSourcesSchema = func() *walker.Schema {
		schema := GetSchemaForTable("cloud_sources")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.CloudSource)(nil)), "cloud_sources")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory_CLOUD_SOURCES, "cloudsource", (*storage.CloudSource)(nil)))
		schema.ScopingResource = resources.Integration
		RegisterTable(schema, CreateTableCloudSourcesStmt)
		mapping.RegisterCategoryToTable(v1.SearchCategory_CLOUD_SOURCES, schema)
		return schema
	}()
)

const (
	// CloudSourcesTableName specifies the name of the table in postgres.
	CloudSourcesTableName = "cloud_sources"
)

// CloudSources holds the Gorm model for Postgres table `cloud_sources`.
type CloudSources struct {
	ID         string                   `gorm:"column:id;type:uuid;primaryKey"`
	Name       string                   `gorm:"column:name;type:varchar;unique"`
	Type       storage.CloudSource_Type `gorm:"column:type;type:integer"`
	Serialized []byte                   `gorm:"column:serialized;type:bytea"`
}
