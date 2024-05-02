//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type PatentType string

const (
	PatentType_Software PatentType = "software"
	PatentType_Database PatentType = "database"
)

func (e *PatentType) Scan(value interface{}) error {
	var enumValue string
	switch val := value.(type) {
	case string:
		enumValue = val
	case []byte:
		enumValue = string(val)
	default:
		return errors.New("jet: Invalid scan value for AllTypesEnum enum. Enum value has to be of type string or []byte")
	}

	switch enumValue {
	case "software":
		*e = PatentType_Software
	case "database":
		*e = PatentType_Database
	default:
		return errors.New("jet: Invalid scan value '" + enumValue + "' for PatentType enum")
	}

	return nil
}

func (e PatentType) String() string {
	return string(e)
}
