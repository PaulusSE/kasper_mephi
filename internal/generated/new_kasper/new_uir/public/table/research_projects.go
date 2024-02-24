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

var ResearchProjects = newResearchProjectsTable("public", "research_projects", "")

type researchProjectsTable struct {
	postgres.Table

	//Columns
	ProjectID   postgres.ColumnString
	WorksID     postgres.ColumnString
	ProjectName postgres.ColumnString
	StartAt     postgres.ColumnTimestampz
	EndAt       postgres.ColumnTimestampz
	AddInfo     postgres.ColumnString
	Grantee     postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ResearchProjectsTable struct {
	researchProjectsTable

	EXCLUDED researchProjectsTable
}

// AS creates new ResearchProjectsTable with assigned alias
func (a ResearchProjectsTable) AS(alias string) *ResearchProjectsTable {
	return newResearchProjectsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ResearchProjectsTable with assigned schema name
func (a ResearchProjectsTable) FromSchema(schemaName string) *ResearchProjectsTable {
	return newResearchProjectsTable(schemaName, a.TableName(), a.Alias())
}

func newResearchProjectsTable(schemaName, tableName, alias string) *ResearchProjectsTable {
	return &ResearchProjectsTable{
		researchProjectsTable: newResearchProjectsTableImpl(schemaName, tableName, alias),
		EXCLUDED:              newResearchProjectsTableImpl("", "excluded", ""),
	}
}

func newResearchProjectsTableImpl(schemaName, tableName, alias string) researchProjectsTable {
	var (
		ProjectIDColumn   = postgres.StringColumn("project_id")
		WorksIDColumn     = postgres.StringColumn("works_id")
		ProjectNameColumn = postgres.StringColumn("project_name")
		StartAtColumn     = postgres.TimestampzColumn("start_at")
		EndAtColumn       = postgres.TimestampzColumn("end_at")
		AddInfoColumn     = postgres.StringColumn("add_info")
		GranteeColumn     = postgres.StringColumn("grantee")
		allColumns        = postgres.ColumnList{ProjectIDColumn, WorksIDColumn, ProjectNameColumn, StartAtColumn, EndAtColumn, AddInfoColumn, GranteeColumn}
		mutableColumns    = postgres.ColumnList{WorksIDColumn, ProjectNameColumn, StartAtColumn, EndAtColumn, AddInfoColumn, GranteeColumn}
	)

	return researchProjectsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ProjectID:   ProjectIDColumn,
		WorksID:     WorksIDColumn,
		ProjectName: ProjectNameColumn,
		StartAt:     StartAtColumn,
		EndAt:       EndAtColumn,
		AddInfo:     AddInfoColumn,
		Grantee:     GranteeColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
