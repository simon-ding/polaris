package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Indexers holds the schema definition for the Indexers entity.
type Indexers struct {
	ent.Schema
}

// Fields of the Indexers.
func (Indexers) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("implementation"),
		field.String("settings"),
		field.Bool("enable_rss"),
		field.Int("priority"),
	}
}

// Edges of the Indexers.
func (Indexers) Edges() []ent.Edge {
	return nil
}
