package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Media holds the schema definition for the Media entity.
type Media struct {
	ent.Schema
}

// Fields of the Media.
func (Media) Fields() []ent.Field {
	return []ent.Field{
		field.Int("tmdb_id"),
		field.String("imdb_id").Optional(),
		field.Enum("media_type").Values("tv", "movie"),
		field.String("name_cn"),
		field.String("name_en"),
		field.String("original_name"),
		field.String("overview"),
		field.Time("created_at").Default(time.Now()),
		field.String("air_date").Default(""),
		field.Enum("resolution").Values("720p", "1080p", "4k").Default("1080p"),
		field.Int("storage_id").Optional(),
		field.String("target_dir").Optional(),
		//field.Bool("download_history_episodes").Optional().Default(false).Comment("tv series only"),
	}
}

// Edges of the Media.
func (Media) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("episodes", Episode.Type),
	}
}
