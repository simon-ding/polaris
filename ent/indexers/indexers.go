// Code generated by ent, DO NOT EDIT.

package indexers

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the indexers type in the database.
	Label = "indexers"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldImplementation holds the string denoting the implementation field in the database.
	FieldImplementation = "implementation"
	// FieldSettings holds the string denoting the settings field in the database.
	FieldSettings = "settings"
	// FieldEnableRss holds the string denoting the enable_rss field in the database.
	FieldEnableRss = "enable_rss"
	// FieldPriority holds the string denoting the priority field in the database.
	FieldPriority = "priority"
	// FieldSeedRatio holds the string denoting the seed_ratio field in the database.
	FieldSeedRatio = "seed_ratio"
	// FieldDisabled holds the string denoting the disabled field in the database.
	FieldDisabled = "disabled"
	// Table holds the table name of the indexers in the database.
	Table = "indexers"
)

// Columns holds all SQL columns for indexers fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldImplementation,
	FieldSettings,
	FieldEnableRss,
	FieldPriority,
	FieldSeedRatio,
	FieldDisabled,
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
	// DefaultEnableRss holds the default value on creation for the "enable_rss" field.
	DefaultEnableRss bool
	// DefaultPriority holds the default value on creation for the "priority" field.
	DefaultPriority int
	// DefaultSeedRatio holds the default value on creation for the "seed_ratio" field.
	DefaultSeedRatio float32
	// DefaultDisabled holds the default value on creation for the "disabled" field.
	DefaultDisabled bool
)

// OrderOption defines the ordering options for the Indexers queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByImplementation orders the results by the implementation field.
func ByImplementation(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldImplementation, opts...).ToFunc()
}

// BySettings orders the results by the settings field.
func BySettings(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSettings, opts...).ToFunc()
}

// ByEnableRss orders the results by the enable_rss field.
func ByEnableRss(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEnableRss, opts...).ToFunc()
}

// ByPriority orders the results by the priority field.
func ByPriority(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPriority, opts...).ToFunc()
}

// BySeedRatio orders the results by the seed_ratio field.
func BySeedRatio(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSeedRatio, opts...).ToFunc()
}

// ByDisabled orders the results by the disabled field.
func ByDisabled(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDisabled, opts...).ToFunc()
}
