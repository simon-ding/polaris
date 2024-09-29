package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Blocklist holds the schema definition for the Blocklist entity.
type Blocklist struct {
	ent.Schema
}

// Fields of the Blocklist.
func (Blocklist) Fields() []ent.Field {
	return  []ent.Field{
		field.Enum("type").Values("media", "torrent"),
		field.String("value"),
	}
}

// Edges of the Blocklist.
func (Blocklist) Edges() []ent.Edge {
	return nil
}
