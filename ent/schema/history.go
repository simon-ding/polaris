package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// History holds the schema definition for the History entity.
type History struct {
	ent.Schema
}

// Fields of the History.
func (History) Fields() []ent.Field {
	return []ent.Field{
		field.Int("series_id"),
		field.Int("episode_id"),
		field.String("source_title"),
		field.Time("date"),
	}
}

// Edges of the History.
func (History) Edges() []ent.Edge {
	return nil
}
