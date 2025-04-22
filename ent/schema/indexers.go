package schema

import (
	"time"

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
		field.String("settings").Optional().Default("").Comment("deprecated, use api_key and url"),
		field.Bool("enable_rss").Default(true),
		field.Int("priority").Default(50),
		field.Float32("seed_ratio").Optional().Default(0).Comment("minimal seed ratio requied, before removing torrent"),
		field.Bool("disabled").Optional().Default(false),
		field.Bool("tv_search").Optional().Default(true),
		field.Bool("movie_search").Optional().Default(true),
		field.String("api_key").Optional(),
		field.String("url").Optional(),
		field.Bool("synced").Optional().Default(false).Comment("synced from prowlarr"),
		field.Time("create_time").Optional().Default(time.Now).Immutable(),
	}
}

// Edges of the Indexers.
func (Indexers) Edges() []ent.Edge {
	return nil
}
