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
		field.Bool("enable_rss").Default(true),
		field.Int("priority").Default(50),
		field.Float32("seed_ratio").Optional().Default(0).Comment("minimal seed ratio requied, before removing torrent"),
		field.Bool("disabled").Optional().Default(false),
	}
}

// Edges of the Indexers.
func (Indexers) Edges() []ent.Edge {
	return nil
}
