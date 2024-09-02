package douban

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type DoulistItem struct {
	Name   string
	ImdbID string
}

func ParseDoulist(doulistUrl string) ([]DoulistItem, error) {
	res, err := http.Get(doulistUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)

	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	doc.Find("")
	return nil, nil
}
