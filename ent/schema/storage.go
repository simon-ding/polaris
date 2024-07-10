package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Storage holds the schema definition for the Storage entity.
type Storage struct {
	ent.Schema
}

// Fields of the Storage.
func (Storage) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("implementation"),
		field.String("path"),
		field.String("user").Optional(),
		field.String("password").Optional(),
		field.Bool("deleted").Default(false),
		field.Bool("default").Default(false),
	}
}

// Edges of the Storage.
func (Storage) Edges() []ent.Edge {
	return nil
}
