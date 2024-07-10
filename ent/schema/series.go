package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
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
		field.String("imdb_id").Optional(),
		field.String("name"),
		field.String("name_en").Optional(),
		field.String("original_name"),
		field.String("overview"),
		field.String("poster_path").Optional(),
		field.Time("created_at").Default(time.Now()),
		field.String("air_date").Default(""),
		field.String("resolution").Default(""),
		field.Int("storage_id").Optional(),
	}
}

// Edges of the Series.
func (Series) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("episodes", Episode.Type),
	}

}
