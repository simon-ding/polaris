package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Episode holds the schema definition for the Epidodes entity.
type Episode struct {
	ent.Schema
}

// Fields of the Episode.
func (Episode) Fields() []ent.Field {
	return []ent.Field{
		field.Int("series_id").Optional(),
		field.Int("season_number").StructTag("json:\"season_number\""),
		field.Int("episode_number").StructTag("json:\"episode_number\""),
		field.String("title"),
		field.String("overview"),
		field.String("air_date"),
		field.String("file_in_storage").Optional(),
	}
}

// Edges of the Episode.
func (Episode) Edges() []ent.Edge {
	return []ent.Edge{
        edge.From("series", Series.Type).
            Ref("episodes").
            Unique().
			Field("series_id"),
    }

}
