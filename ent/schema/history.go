package schema

import (
	"time"

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
		field.Int("media_id"),
		//field.Int("episode_id").Optional().Comment("deprecated"),
		field.Ints("episode_nums").Optional(),
		field.Int("season_num").Optional(),
		field.String("source_title"),
		field.Time("date"),
		field.String("target_dir"),
		field.Int("size").Default(0),
		field.Int("download_client_id").Optional(),
		field.Int("indexer_id").Optional(),
		field.String("link").Optional().Comment("deprecated, use hash instead"), //should be magnet link
		field.String("hash").Optional().Comment("torrent hash"),
		field.Enum("status").Values("running", "success", "fail", "uploading", "seeding", "removed"),
		field.Time("create_time").Optional().Default(time.Now).Immutable(),
		//field.String("saved").Optional().Comment("deprecated"), //deprecated
	}
}

// Edges of the History.
func (History) Edges() []ent.Edge {
	return nil
}
