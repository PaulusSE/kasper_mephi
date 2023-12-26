//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
	"time"
)

type Students struct {
	ClientID           uuid.UUID
	StudentID          uuid.UUID `sql:"primary_key"`
	FullName           string
	Department         string
	EnrollmentOrder    string
	TitlePagePath      string
	ExplanatoryNoteURL string
	Specialization     *StudentSpecialization
	ActualSemester     int32
	SupervisorID       uuid.UUID
	StartDate          *time.Time
	AcademicLeave      bool
	DissertationTitle  string
	GroupNumber        *string
}
