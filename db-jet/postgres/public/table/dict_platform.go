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

var DictPlatform = newDictPlatformTable("public", "dict_platform", "")

type dictPlatformTable struct {
	postgres.Table

	// Columns
	ID          postgres.ColumnInteger
	Name        postgres.ColumnString
	DisplayName postgres.ColumnString
	CreatedTs   postgres.ColumnTimestampz
	ModifiedTs  postgres.ColumnTimestampz
	DeletedTs   postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type DictPlatformTable struct {
	dictPlatformTable

	EXCLUDED dictPlatformTable
}

// AS creates new DictPlatformTable with assigned alias
func (a DictPlatformTable) AS(alias string) *DictPlatformTable {
	return newDictPlatformTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new DictPlatformTable with assigned schema name
func (a DictPlatformTable) FromSchema(schemaName string) *DictPlatformTable {
	return newDictPlatformTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new DictPlatformTable with assigned table prefix
func (a DictPlatformTable) WithPrefix(prefix string) *DictPlatformTable {
	return newDictPlatformTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new DictPlatformTable with assigned table suffix
func (a DictPlatformTable) WithSuffix(suffix string) *DictPlatformTable {
	return newDictPlatformTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newDictPlatformTable(schemaName, tableName, alias string) *DictPlatformTable {
	return &DictPlatformTable{
		dictPlatformTable: newDictPlatformTableImpl(schemaName, tableName, alias),
		EXCLUDED:          newDictPlatformTableImpl("", "excluded", ""),
	}
}

func newDictPlatformTableImpl(schemaName, tableName, alias string) dictPlatformTable {
	var (
		IDColumn          = postgres.IntegerColumn("id")
		NameColumn        = postgres.StringColumn("name")
		DisplayNameColumn = postgres.StringColumn("display_name")
		CreatedTsColumn   = postgres.TimestampzColumn("created_ts")
		ModifiedTsColumn  = postgres.TimestampzColumn("modified_ts")
		DeletedTsColumn   = postgres.TimestampzColumn("deleted_ts")
		allColumns        = postgres.ColumnList{IDColumn, NameColumn, DisplayNameColumn, CreatedTsColumn, ModifiedTsColumn, DeletedTsColumn}
		mutableColumns    = postgres.ColumnList{NameColumn, DisplayNameColumn, CreatedTsColumn, ModifiedTsColumn, DeletedTsColumn}
	)

	return dictPlatformTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		Name:        NameColumn,
		DisplayName: DisplayNameColumn,
		CreatedTs:   CreatedTsColumn,
		ModifiedTs:  ModifiedTsColumn,
		DeletedTs:   DeletedTsColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
