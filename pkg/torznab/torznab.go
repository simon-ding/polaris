package torznab

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"polaris/log"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Response struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Torznab string   `xml:"torznab,attr"`
	Channel struct {
		Text string `xml:",chardata"`
		Link struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Language    string `xml:"language"`
		Category    string `xml:"category"`
		Item        []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Text           string `xml:",chardata"`
	Title          string `xml:"title"`
	Guid           string `xml:"guid"`
	Jackettindexer struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id,attr"`
	} `xml:"jackettindexer"`
	Type        string   `xml:"type"`
	Comments    string   `xml:"comments"`
	PubDate     string   `xml:"pubDate"`
	Size        string   `xml:"size"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	Category    []string `xml:"category"`
	Enclosure   struct {
		Text   string `xml:",chardata"`
		URL    string `xml:"url,attr"`
		Length string `xml:"length,attr"`
		Type   string `xml:"type,attr"`
	} `xml:"enclosure"`
	Attr []struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	} `xml:"attr"`
}

func (i *Item) GetAttr(key string) string {
	for _, a := range i.Attr {
		if a.Name == key {
			return a.Value
		}
	}
	return ""
}
func (r *Response) ToResults() []Result {
	var res []Result
	for _, item := range r.Channel.Item {
		r := Result{
			Name:     item.Title,
			Magnet:   item.Link,
			Size:     mustAtoI(item.Size),
			Seeders:  mustAtoI(item.GetAttr("seeders")),
			Peers:    mustAtoI(item.GetAttr("peers")),
			Category: mustAtoI(item.GetAttr("category")),
			Source:   r.Channel.Title,
		}
		res = append(res, r)
	}
	return res
}

func mustAtoI(key string) int {
	i, err := strconv.Atoi(key)
	if err != nil {
		log.Errorf("must atoi error: %v", err)
		panic(err)
	}
	return i
}
func Search(torznabUrl, api, keyWord string) ([]Result, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, torznabUrl, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}
	var q = url.Values{}
	q.Add("apikey", api)
	q.Add("t", "search")
	q.Add("q", keyWord)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do http")
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read http body")
	}
	var res Response
	err = xml.Unmarshal(data, &res)
	if err != nil {
		return nil, errors.Wrap(err, "json unmarshal")
	}
	return res.ToResults(), nil
}

type Result struct {
	Name     string
	Magnet   string
	Size     int
	Seeders  int
	Peers    int
	Category int
	Source   string
}
