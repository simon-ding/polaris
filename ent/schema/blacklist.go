package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Blacklist holds the schema definition for the Blacklist entity.
type Blacklist struct {
	ent.Schema
}

// Fields of the Blacklist.
func (Blacklist) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").Values("media", "torrent").Default("torrent"),
		field.String("torrent_hash").Optional(),
		field.String("torrent_name").Optional(),
		field.Int("media_id").Optional(),
		field.Time("create_time").Optional().Default(time.Now).Immutable(),
		field.String("notes").Optional(),
	}
}

// Edges of the Blacklist.
func (Blacklist) Edges() []ent.Edge {
	return nil
}
