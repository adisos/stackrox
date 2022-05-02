// Code generated by pg-bindings generator. DO NOT EDIT.

package schema

import (
	"reflect"

	"github.com/stackrox/rox/central/globaldb"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/postgres"
	"github.com/stackrox/rox/pkg/postgres/walker"
	"github.com/stackrox/rox/pkg/search"
)

var (
	// CreateTableTestg3grandchild1Stmt holds the create statement for table `testg3grandchild1`.
	CreateTableTestg3grandchild1Stmt = &postgres.CreateStmts{
		Table: `
               create table if not exists testg3grandchild1 (
                   Id varchar,
                   Val varchar,
                   serialized bytea,
                   PRIMARY KEY(Id)
               )
               `,
		Indexes:  []string{},
		Children: []*postgres.CreateStmts{},
	}

	// Testg3grandchild1Schema is the go schema for table `testg3grandchild1`.
	Testg3grandchild1Schema = func() *walker.Schema {
		schema := globaldb.GetSchemaForTable("testg3grandchild1")
		if schema != nil {
			return schema
		}
		schema = walker.Walk(reflect.TypeOf((*storage.TestG3GrandChild1)(nil)), "testg3grandchild1")
		schema.SetOptionsMap(search.Walk(v1.SearchCategory(67), "testg3grandchild1", (*storage.TestG3GrandChild1)(nil)))
		globaldb.RegisterTable(schema)
		return schema
	}()
)
