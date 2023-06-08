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

var Audit = newAuditTable("public", "audit", "")

type auditTable struct {
	postgres.Table

	// Columns
	ID        postgres.ColumnString
	CreateTs  postgres.ColumnTimestampz
	Username  postgres.ColumnString
	IP        postgres.ColumnString
	RequestID postgres.ColumnString
	Message   postgres.ColumnString
	Meta      postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type AuditTable struct {
	auditTable

	EXCLUDED auditTable
}

// AS creates new AuditTable with assigned alias
func (a AuditTable) AS(alias string) *AuditTable {
	return newAuditTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new AuditTable with assigned schema name
func (a AuditTable) FromSchema(schemaName string) *AuditTable {
	return newAuditTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new AuditTable with assigned table prefix
func (a AuditTable) WithPrefix(prefix string) *AuditTable {
	return newAuditTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new AuditTable with assigned table suffix
func (a AuditTable) WithSuffix(suffix string) *AuditTable {
	return newAuditTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newAuditTable(schemaName, tableName, alias string) *AuditTable {
	return &AuditTable{
		auditTable: newAuditTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newAuditTableImpl("", "excluded", ""),
	}
}

func newAuditTableImpl(schemaName, tableName, alias string) auditTable {
	var (
		IDColumn        = postgres.StringColumn("id")
		CreateTsColumn  = postgres.TimestampzColumn("create_ts")
		UsernameColumn  = postgres.StringColumn("username")
		IPColumn        = postgres.StringColumn("ip")
		RequestIDColumn = postgres.StringColumn("request_id")
		MessageColumn   = postgres.StringColumn("message")
		MetaColumn      = postgres.StringColumn("meta")
		allColumns      = postgres.ColumnList{IDColumn, CreateTsColumn, UsernameColumn, IPColumn, RequestIDColumn, MessageColumn, MetaColumn}
		mutableColumns  = postgres.ColumnList{CreateTsColumn, UsernameColumn, IPColumn, RequestIDColumn, MessageColumn, MetaColumn}
	)

	return auditTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		CreateTs:  CreateTsColumn,
		Username:  UsernameColumn,
		IP:        IPColumn,
		RequestID: RequestIDColumn,
		Message:   MessageColumn,
		Meta:      MetaColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}