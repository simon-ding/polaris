package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// NotificationClient holds the schema definition for the NotificationClient entity.
type NotificationClient struct {
	ent.Schema
}

// Fields of the NotificationClient.
func (NotificationClient) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("service"),
		field.String("settings"),
		field.Bool("enabled").Default(true),
	}
}

// Edges of the NotificationClient.
func (NotificationClient) Edges() []ent.Edge {
	return nil
}
