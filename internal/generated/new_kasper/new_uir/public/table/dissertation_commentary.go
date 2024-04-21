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

var DissertationCommentary = newDissertationCommentaryTable("public", "dissertation_commentary", "")

type dissertationCommentaryTable struct {
	postgres.Table

	//Columns
	CommentaryID postgres.ColumnString
	StudentID    postgres.ColumnString
	Semester     postgres.ColumnInteger
	Commentary   postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type DissertationCommentaryTable struct {
	dissertationCommentaryTable

	EXCLUDED dissertationCommentaryTable
}

// AS creates new DissertationCommentaryTable with assigned alias
func (a DissertationCommentaryTable) AS(alias string) *DissertationCommentaryTable {
	return newDissertationCommentaryTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new DissertationCommentaryTable with assigned schema name
func (a DissertationCommentaryTable) FromSchema(schemaName string) *DissertationCommentaryTable {
	return newDissertationCommentaryTable(schemaName, a.TableName(), a.Alias())
}

func newDissertationCommentaryTable(schemaName, tableName, alias string) *DissertationCommentaryTable {
	return &DissertationCommentaryTable{
		dissertationCommentaryTable: newDissertationCommentaryTableImpl(schemaName, tableName, alias),
		EXCLUDED:                    newDissertationCommentaryTableImpl("", "excluded", ""),
	}
}

func newDissertationCommentaryTableImpl(schemaName, tableName, alias string) dissertationCommentaryTable {
	var (
		CommentaryIDColumn = postgres.StringColumn("commentary_id")
		StudentIDColumn    = postgres.StringColumn("student_id")
		SemesterColumn     = postgres.IntegerColumn("semester")
		CommentaryColumn   = postgres.StringColumn("commentary")
		allColumns         = postgres.ColumnList{CommentaryIDColumn, StudentIDColumn, SemesterColumn, CommentaryColumn}
		mutableColumns     = postgres.ColumnList{StudentIDColumn, SemesterColumn, CommentaryColumn}
	)

	return dissertationCommentaryTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		CommentaryID: CommentaryIDColumn,
		StudentID:    StudentIDColumn,
		Semester:     SemesterColumn,
		Commentary:   CommentaryColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
