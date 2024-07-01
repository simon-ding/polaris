// Code generated by ent, DO NOT EDIT.

package history

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the history type in the database.
	Label = "history"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldSeriesID holds the string denoting the series_id field in the database.
	FieldSeriesID = "series_id"
	// FieldEpisodeID holds the string denoting the episode_id field in the database.
	FieldEpisodeID = "episode_id"
	// FieldSourceTitle holds the string denoting the source_title field in the database.
	FieldSourceTitle = "source_title"
	// FieldDate holds the string denoting the date field in the database.
	FieldDate = "date"
	// Table holds the table name of the history in the database.
	Table = "histories"
)

// Columns holds all SQL columns for history fields.
var Columns = []string{
	FieldID,
	FieldSeriesID,
	FieldEpisodeID,
	FieldSourceTitle,
	FieldDate,
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

// OrderOption defines the ordering options for the History queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// BySeriesID orders the results by the series_id field.
func BySeriesID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSeriesID, opts...).ToFunc()
}

// ByEpisodeID orders the results by the episode_id field.
func ByEpisodeID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEpisodeID, opts...).ToFunc()
}

// BySourceTitle orders the results by the source_title field.
func BySourceTitle(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSourceTitle, opts...).ToFunc()
}

// ByDate orders the results by the date field.
func ByDate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDate, opts...).ToFunc()
}
