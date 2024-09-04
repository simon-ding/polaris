package importlist

type Item struct {
	Title  string
	Year   int
	ImdbID string
	TvdbID string
	TmdbID string
}

type Response struct {
	Items []Item
}
