//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/postgres"

var ApprovalStatus = &struct {
	Todo       postgres.StringExpression
	Approved   postgres.StringExpression
	OnReview   postgres.StringExpression
	InProgress postgres.StringExpression
	Empty      postgres.StringExpression
}{
	Todo:       postgres.NewEnumValue("todo"),
	Approved:   postgres.NewEnumValue("approved"),
	OnReview:   postgres.NewEnumValue("on review"),
	InProgress: postgres.NewEnumValue("in progress"),
	Empty:      postgres.NewEnumValue("empty"),
}
