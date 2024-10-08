// Code generated by ent, DO NOT EDIT.

package blacklist

import (
	"fmt"
	"polaris/ent/schema"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the blacklist type in the database.
	Label = "blacklist"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldValue holds the string denoting the value field in the database.
	FieldValue = "value"
	// FieldNotes holds the string denoting the notes field in the database.
	FieldNotes = "notes"
	// Table holds the table name of the blacklist in the database.
	Table = "blacklists"
)

// Columns holds all SQL columns for blacklist fields.
var Columns = []string{
	FieldID,
	FieldType,
	FieldValue,
	FieldNotes,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultValue holds the default value on creation for the "value" field.
	DefaultValue schema.BlacklistValue
)

// Type defines the type for the "type" enum field.
type Type string

// Type values.
const (
	TypeMedia   Type = "media"
	TypeTorrent Type = "torrent"
)

func (_type Type) String() string {
	return string(_type)
}

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type Type) error {
	switch _type {
	case TypeMedia, TypeTorrent:
		return nil
	default:
		return fmt.Errorf("blacklist: invalid enum value for type field: %q", _type)
	}
}

// OrderOption defines the ordering options for the Blacklist queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByNotes orders the results by the notes field.
func ByNotes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldNotes, opts...).ToFunc()
}
