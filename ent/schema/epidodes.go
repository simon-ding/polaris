package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Epidodes holds the schema definition for the Epidodes entity.
type Epidodes struct {
	ent.Schema
}

// Fields of the Epidodes.
func (Epidodes) Fields() []ent.Field {
	return []ent.Field{
		field.Int("series_id"),
		field.Int("season_number"),
		field.Int("episode_number"),
		field.String("title"),
		field.String("overview"),
		field.String("air_date"),
	}
}

// Edges of the Epidodes.
func (Epidodes) Edges() []ent.Edge {
	return nil
}
