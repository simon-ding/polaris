package schema

import (
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
		field.Enum("type").Values("media", "torrent"),
		field.JSON("value", BlacklistValue{}).Default(BlacklistValue{}),
		field.String("notes").Optional(),
	}
}

// Edges of the Blacklist.
func (Blacklist) Edges() []ent.Edge {
	return nil
}

type BlacklistValue struct {
	TmdbID      int    `json:"tmdb_id"`
	TorrentHash string `json:"torrent_hash"`
}
