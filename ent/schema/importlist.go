package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// ImportList holds the schema definition for the ImportList entity.
type ImportList struct {
	ent.Schema
}

// Fields of the ImportList.
func (ImportList) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Enum("type").Values("plex", "doulist"),
		field.String("url").Optional(),
		field.String("qulity"),
		field.Int("storage_id"),
		field.JSON("settings", ImportListSettings{}).Optional(),
	}
}

// Edges of the ImportList.
func (ImportList) Edges() []ent.Edge {
	return nil
}

type ImportListSettings struct {
	//Url string `json:"url"`
}