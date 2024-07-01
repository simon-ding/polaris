package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// DownloadClients holds the schema definition for the DownloadClients entity.
type DownloadClients struct {
	ent.Schema
}

// Fields of the DownloadClients.
func (DownloadClients) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("enable"),
		field.String("name"),
		field.String("implementation"),
		field.String("url"),
		field.String("user").Default(""),
		field.String("password").Default(""),
		field.String("settings").Default(""),
		field.String("priority").Default(""),
		field.Bool("remove_completed_downloads").Default(true),
		field.Bool("remove_failed_downloads").Default(true),
		field.String("tags").Default(""),
	}
}

// Edges of the DownloadClients.
func (DownloadClients) Edges() []ent.Edge {
	return nil
}
