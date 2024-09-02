package plexwatchlist

import (
	"encoding/xml"
	"io"
	"net/http"
	"polaris/pkg/importlist"
	"strings"

	"github.com/pkg/errors"
)

type Response struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Atom    string   `xml:"atom,attr"`
	Media   string   `xml:"media,attr"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description string `xml:"description"`
		Category    string `xml:"category"`
		Item        []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			PubDate     string `xml:"pubDate"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			Category    string `xml:"category"`
			Credit      []struct {
				Text string `xml:",chardata"`
				Role string `xml:"role,attr"`
			} `xml:"credit"`
			Thumbnail struct {
				Text string `xml:",chardata"`
				URL  string `xml:"url,attr"`
			} `xml:"thumbnail"`
			Keywords string `xml:"keywords"`
			Rating   struct {
				Text   string `xml:",chardata"`
				Scheme string `xml:"scheme,attr"`
			} `xml:"rating"`
			Guid struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
		} `xml:"item"`
	} `xml:"channel"`
}

func (r *Response) convert() *importlist.Response {
	res := &importlist.Response{}
	for _, im := range r.Channel.Item {
		item := importlist.Item{
			Title: im.Title,
		}
		id := strings.ToLower(im.Guid.Text)
		if strings.HasPrefix(id, "tvdb") {
			tvdbid := strings.TrimPrefix(id, "tvdb://")
			item.TvdbID = tvdbid
		} else if strings.HasPrefix(id, "imdb") {
			imdbid := strings.TrimPrefix(id, "imdb://")
			item.ImdbID = imdbid
		} else if strings.HasPrefix(id, "tmdb") {
			tmdbid := strings.TrimPrefix(id, "tmdb://")
			item.TmdbID = tmdbid
		}
		res.Items = append(res.Items, item)
	}
	return res
}

func ParsePlexWatchlist(url string) (*importlist.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "http get")
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read data")
	}
	var rrr Response
	err = xml.Unmarshal(data, &rrr)
	if err != nil {
		return nil, errors.Wrap(err, "xml")
	}
	return rrr.convert(), nil
}
