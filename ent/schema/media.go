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
		field.Enum("resolution").Values("720p", "1080p", "2160p", "any").Default("1080p"),
		field.Int("storage_id").Optional(),
		field.String("target_dir").Optional(),
		field.Bool("download_history_episodes").Optional().Default(false).Comment("tv series only"),
		field.JSON("limiter", MediaLimiter{}).Optional(),
		field.JSON("extras", MediaExtras{}).Optional(),
		field.JSON("alternative_titles", []AlternativeTilte{}).Optional(),
		field.Time("create_time").Optional().Default(time.Now).Immutable(),
	}
}

// Edges of the Media.
func (Media) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("episodes", Episode.Type),
	}
}

type AlternativeTilte struct {
	Iso3166_1 string `json:"iso_3166_1"`
	Title     string `json:"title"`
	Type      string `json:"type"`
}

type MediaLimiter struct {
	SizeMin    int64 `json:"size_min"` //in B
	SizeMax    int64 `json:"size_max"` //in B
	PreferSize int64 `json:"prefer_max"`
}

type MediaExtras struct {
	IsAdultMovie bool   `json:"is_adult_movie"`
	JavId        string `json:"javid"`
	//OriginCountry    []string `json:"origin_country"`
	OriginalLanguage string `json:"original_language"`
	Genres           []Genre `json:"genres"`
}

type Genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (m *MediaExtras) IsJav() bool {
	return m.IsAdultMovie && m.JavId != ""
}
