package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Series holds the schema definition for the Series entity.
type Series struct {
	ent.Schema
}

// Fields of the Series.
func (Series) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tmdb_id"),
		field.String("imdb_id"),
		field.String("title"),
		field.String("original_name"),
		field.String("overview"),
		field.String("path"),
	}
}

// Edges of the Series.
func (Series) Edges() []ent.Edge {
	return nil
}
