//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"github.com/google/uuid"
)

type SemesterCount struct {
	CountID  uuid.UUID `sql:"primary_key"`
	Amount   int32
	Archived bool
}
