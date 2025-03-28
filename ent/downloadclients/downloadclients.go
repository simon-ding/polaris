// Code generated by ent, DO NOT EDIT.

package downloadclients

import (
	"fmt"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the downloadclients type in the database.
	Label = "download_clients"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldEnable holds the string denoting the enable field in the database.
	FieldEnable = "enable"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldImplementation holds the string denoting the implementation field in the database.
	FieldImplementation = "implementation"
	// FieldURL holds the string denoting the url field in the database.
	FieldURL = "url"
	// FieldUser holds the string denoting the user field in the database.
	FieldUser = "user"
	// FieldPassword holds the string denoting the password field in the database.
	FieldPassword = "password"
	// FieldSettings holds the string denoting the settings field in the database.
	FieldSettings = "settings"
	// FieldPriority1 holds the string denoting the priority1 field in the database.
	FieldPriority1 = "priority1"
	// FieldRemoveCompletedDownloads holds the string denoting the remove_completed_downloads field in the database.
	FieldRemoveCompletedDownloads = "remove_completed_downloads"
	// FieldRemoveFailedDownloads holds the string denoting the remove_failed_downloads field in the database.
	FieldRemoveFailedDownloads = "remove_failed_downloads"
	// FieldTags holds the string denoting the tags field in the database.
	FieldTags = "tags"
	// Table holds the table name of the downloadclients in the database.
	Table = "download_clients"
)

// Columns holds all SQL columns for downloadclients fields.
var Columns = []string{
	FieldID,
	FieldEnable,
	FieldName,
	FieldImplementation,
	FieldURL,
	FieldUser,
	FieldPassword,
	FieldSettings,
	FieldPriority1,
	FieldRemoveCompletedDownloads,
	FieldRemoveFailedDownloads,
	FieldTags,
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
	// DefaultUser holds the default value on creation for the "user" field.
	DefaultUser string
	// DefaultPassword holds the default value on creation for the "password" field.
	DefaultPassword string
	// DefaultSettings holds the default value on creation for the "settings" field.
	DefaultSettings string
	// DefaultPriority1 holds the default value on creation for the "priority1" field.
	DefaultPriority1 int
	// Priority1Validator is a validator for the "priority1" field. It is called by the builders before save.
	Priority1Validator func(int) error
	// DefaultRemoveCompletedDownloads holds the default value on creation for the "remove_completed_downloads" field.
	DefaultRemoveCompletedDownloads bool
	// DefaultRemoveFailedDownloads holds the default value on creation for the "remove_failed_downloads" field.
	DefaultRemoveFailedDownloads bool
	// DefaultTags holds the default value on creation for the "tags" field.
	DefaultTags string
)

// Implementation defines the type for the "implementation" enum field.
type Implementation string

// Implementation values.
const (
	ImplementationTransmission Implementation = "transmission"
	ImplementationQbittorrent  Implementation = "qbittorrent"
	ImplementationBuildin      Implementation = "buildin"
)

func (i Implementation) String() string {
	return string(i)
}

// ImplementationValidator is a validator for the "implementation" field enum values. It is called by the builders before save.
func ImplementationValidator(i Implementation) error {
	switch i {
	case ImplementationTransmission, ImplementationQbittorrent, ImplementationBuildin:
		return nil
	default:
		return fmt.Errorf("downloadclients: invalid enum value for implementation field: %q", i)
	}
}

// OrderOption defines the ordering options for the DownloadClients queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByEnable orders the results by the enable field.
func ByEnable(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEnable, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByImplementation orders the results by the implementation field.
func ByImplementation(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldImplementation, opts...).ToFunc()
}

// ByURL orders the results by the url field.
func ByURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldURL, opts...).ToFunc()
}

// ByUser orders the results by the user field.
func ByUser(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUser, opts...).ToFunc()
}

// ByPassword orders the results by the password field.
func ByPassword(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPassword, opts...).ToFunc()
}

// BySettings orders the results by the settings field.
func BySettings(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSettings, opts...).ToFunc()
}

// ByPriority1 orders the results by the priority1 field.
func ByPriority1(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPriority1, opts...).ToFunc()
}

// ByRemoveCompletedDownloads orders the results by the remove_completed_downloads field.
func ByRemoveCompletedDownloads(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRemoveCompletedDownloads, opts...).ToFunc()
}

// ByRemoveFailedDownloads orders the results by the remove_failed_downloads field.
func ByRemoveFailedDownloads(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRemoveFailedDownloads, opts...).ToFunc()
}

// ByTags orders the results by the tags field.
func ByTags(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTags, opts...).ToFunc()
}
