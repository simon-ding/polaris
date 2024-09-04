package douban

import (
	"fmt"
	"net/http"
	"polaris/log"
	"polaris/pkg/importlist"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"

func ParseDoulist(doulistUrl string) (*importlist.Response, error) {
	if !strings.Contains(doulistUrl, "doulist") {
		return nil, fmt.Errorf("not doulist")
	}

	req, err := http.NewRequest("GET", doulistUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ua)

	res, err := http.DefaultClient.Do(req)
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
	doc.Find("div[class=doulist-item]").Each(func(i int, selection *goquery.Selection) {
		titleDiv := selection.Find("div[class=title]")
		link := titleDiv.Find("div>a")
		href, ok := link.Attr("href")
		if !ok {
			return
		}
		abstract := selection.Find("div[class=abstract]")

		lines := strings.Split(abstract.Text(), "\n")
		year := 0
		for _, l := range lines {
			if strings.Contains(l, "年份") {
				ppp := strings.Split(l, ":")
				if len(ppp) < 2 {
					continue
				} else {
					n := ppp[1]
					n1, err := strconv.Atoi(n)
					if err != nil {
						log.Errorf("convert year number %s to int error: %v", n, err)
						continue
					}
					year = n1
				}
			}
		}

		item := importlist.Item{
			Title: strings.TrimSpace(link.Text()),
			Year:  year,
		}
		_ = item
		println(link.Text(), href)
	})
	return nil, nil
}

func parseDetailPage(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", ua)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)

	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}
	_ = doc
	return "", nil
}
