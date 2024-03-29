//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/sqlite"
)

var TextPair = newTextPairTable("", "text_pair", "")

type textPairTable struct {
	sqlite.Table

	// Columns
	ID        sqlite.ColumnInteger
	UserID    sqlite.ColumnInteger
	Text1     sqlite.ColumnString
	Text2     sqlite.ColumnString
	IsLearned sqlite.ColumnFloat
	CreatedAt sqlite.ColumnTimestamp
	UpdatedAt sqlite.ColumnTimestamp
	DeletedAt sqlite.ColumnTimestamp

	AllColumns     sqlite.ColumnList
	MutableColumns sqlite.ColumnList
}

type TextPairTable struct {
	textPairTable

	EXCLUDED textPairTable
}

// AS creates new TextPairTable with assigned alias
func (a TextPairTable) AS(alias string) *TextPairTable {
	return newTextPairTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new TextPairTable with assigned schema name
func (a TextPairTable) FromSchema(schemaName string) *TextPairTable {
	return newTextPairTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new TextPairTable with assigned table prefix
func (a TextPairTable) WithPrefix(prefix string) *TextPairTable {
	return newTextPairTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new TextPairTable with assigned table suffix
func (a TextPairTable) WithSuffix(suffix string) *TextPairTable {
	return newTextPairTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newTextPairTable(schemaName, tableName, alias string) *TextPairTable {
	return &TextPairTable{
		textPairTable: newTextPairTableImpl(schemaName, tableName, alias),
		EXCLUDED:      newTextPairTableImpl("", "excluded", ""),
	}
}

func newTextPairTableImpl(schemaName, tableName, alias string) textPairTable {
	var (
		IDColumn        = sqlite.IntegerColumn("id")
		UserIDColumn    = sqlite.IntegerColumn("user_id")
		Text1Column     = sqlite.StringColumn("text1")
		Text2Column     = sqlite.StringColumn("text2")
		IsLearnedColumn = sqlite.FloatColumn("is_learned")
		CreatedAtColumn = sqlite.TimestampColumn("created_at")
		UpdatedAtColumn = sqlite.TimestampColumn("updated_at")
		DeletedAtColumn = sqlite.TimestampColumn("deleted_at")
		allColumns      = sqlite.ColumnList{IDColumn, UserIDColumn, Text1Column, Text2Column, IsLearnedColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn}
		mutableColumns  = sqlite.ColumnList{UserIDColumn, Text1Column, Text2Column, IsLearnedColumn, CreatedAtColumn, UpdatedAtColumn, DeletedAtColumn}
	)

	return textPairTable{
		Table: sqlite.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:        IDColumn,
		UserID:    UserIDColumn,
		Text1:     Text1Column,
		Text2:     Text2Column,
		IsLearned: IsLearnedColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,
		DeletedAt: DeletedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
