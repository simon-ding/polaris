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
		field.String("settings"),
		field.String("priority"),
		field.Bool("remove_completed_downloads"),
		field.Bool("remove_failed_downloads"),
		field.String("tags"),
	}
}

// Edges of the DownloadClients.
func (DownloadClients) Edges() []ent.Edge {
	return nil
}
