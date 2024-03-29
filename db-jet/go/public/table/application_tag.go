//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var ApplicationTag = newApplicationTagTable("public", "application_tag", "")

type applicationTagTable struct {
	postgres.Table

	// Columns
	ID            postgres.ColumnInteger
	ApplicationID postgres.ColumnInteger
	TagID         postgres.ColumnInteger
	CreatedTs     postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ApplicationTagTable struct {
	applicationTagTable

	EXCLUDED applicationTagTable
}

// AS creates new ApplicationTagTable with assigned alias
func (a ApplicationTagTable) AS(alias string) *ApplicationTagTable {
	return newApplicationTagTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ApplicationTagTable with assigned schema name
func (a ApplicationTagTable) FromSchema(schemaName string) *ApplicationTagTable {
	return newApplicationTagTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ApplicationTagTable with assigned table prefix
func (a ApplicationTagTable) WithPrefix(prefix string) *ApplicationTagTable {
	return newApplicationTagTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ApplicationTagTable with assigned table suffix
func (a ApplicationTagTable) WithSuffix(suffix string) *ApplicationTagTable {
	return newApplicationTagTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newApplicationTagTable(schemaName, tableName, alias string) *ApplicationTagTable {
	return &ApplicationTagTable{
		applicationTagTable: newApplicationTagTableImpl(schemaName, tableName, alias),
		EXCLUDED:            newApplicationTagTableImpl("", "excluded", ""),
	}
}

func newApplicationTagTableImpl(schemaName, tableName, alias string) applicationTagTable {
	var (
		IDColumn            = postgres.IntegerColumn("id")
		ApplicationIDColumn = postgres.IntegerColumn("application_id")
		TagIDColumn         = postgres.IntegerColumn("tag_id")
		CreatedTsColumn     = postgres.TimestampzColumn("created_ts")
		allColumns          = postgres.ColumnList{IDColumn, ApplicationIDColumn, TagIDColumn, CreatedTsColumn}
		mutableColumns      = postgres.ColumnList{ApplicationIDColumn, TagIDColumn, CreatedTsColumn}
	)

	return applicationTagTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:            IDColumn,
		ApplicationID: ApplicationIDColumn,
		TagID:         TagIDColumn,
		CreatedTs:     CreatedTsColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
