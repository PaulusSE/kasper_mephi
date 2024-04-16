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

var Students = newStudentsTable("public", "students", "")

type studentsTable struct {
	postgres.Table

	//Columns
	StudentID       postgres.ColumnString
	UserID          postgres.ColumnString
	FullName        postgres.ColumnString
	SpecID          postgres.ColumnInteger
	ActualSemester  postgres.ColumnInteger
	Years           postgres.ColumnInteger
	StartDate       postgres.ColumnTimestampz
	StudyingStatus  postgres.ColumnString
	GroupID         postgres.ColumnInteger
	Status          postgres.ColumnString
	CanEdit         postgres.ColumnBool
	Progressiveness postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type StudentsTable struct {
	studentsTable

	EXCLUDED studentsTable
}

// AS creates new StudentsTable with assigned alias
func (a StudentsTable) AS(alias string) *StudentsTable {
	return newStudentsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new StudentsTable with assigned schema name
func (a StudentsTable) FromSchema(schemaName string) *StudentsTable {
	return newStudentsTable(schemaName, a.TableName(), a.Alias())
}

func newStudentsTable(schemaName, tableName, alias string) *StudentsTable {
	return &StudentsTable{
		studentsTable: newStudentsTableImpl(schemaName, tableName, alias),
		EXCLUDED:      newStudentsTableImpl("", "excluded", ""),
	}
}

func newStudentsTableImpl(schemaName, tableName, alias string) studentsTable {
	var (
		StudentIDColumn       = postgres.StringColumn("student_id")
		UserIDColumn          = postgres.StringColumn("user_id")
		FullNameColumn        = postgres.StringColumn("full_name")
		SpecIDColumn          = postgres.IntegerColumn("spec_id")
		ActualSemesterColumn  = postgres.IntegerColumn("actual_semester")
		YearsColumn           = postgres.IntegerColumn("years")
		StartDateColumn       = postgres.TimestampzColumn("start_date")
		StudyingStatusColumn  = postgres.StringColumn("studying_status")
		GroupIDColumn         = postgres.IntegerColumn("group_id")
		StatusColumn          = postgres.StringColumn("status")
		CanEditColumn         = postgres.BoolColumn("can_edit")
		ProgressivenessColumn = postgres.IntegerColumn("progressiveness")
		allColumns            = postgres.ColumnList{StudentIDColumn, UserIDColumn, FullNameColumn, SpecIDColumn, ActualSemesterColumn, YearsColumn, StartDateColumn, StudyingStatusColumn, GroupIDColumn, StatusColumn, CanEditColumn, ProgressivenessColumn}
		mutableColumns        = postgres.ColumnList{UserIDColumn, FullNameColumn, SpecIDColumn, ActualSemesterColumn, YearsColumn, StartDateColumn, StudyingStatusColumn, GroupIDColumn, StatusColumn, CanEditColumn, ProgressivenessColumn}
	)

	return studentsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		StudentID:       StudentIDColumn,
		UserID:          UserIDColumn,
		FullName:        FullNameColumn,
		SpecID:          SpecIDColumn,
		ActualSemester:  ActualSemesterColumn,
		Years:           YearsColumn,
		StartDate:       StartDateColumn,
		StudyingStatus:  StudyingStatusColumn,
		GroupID:         GroupIDColumn,
		Status:          StatusColumn,
		CanEdit:         CanEditColumn,
		Progressiveness: ProgressivenessColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
