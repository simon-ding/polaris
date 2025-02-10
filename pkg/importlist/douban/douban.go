package douban

import (
	"fmt"
	"io"
	"net/http"
	"polaris/log"
	"polaris/pkg/importlist"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"

func ParseDoulist(doulistUrl string) (*importlist.Response, error) {
	if !strings.Contains(doulistUrl, "doulist") {
		return nil, fmt.Errorf("not doulist")
	}
	res, err := doHttpReq("GET", doulistUrl, nil)
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
	var items []importlist.Item
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
					n1, err := strconv.Atoi(strings.TrimSpace(n))
					if err != nil {
						log.Errorf("convert year number %s to int error: %v", n, err)
						continue
					}
					year = n1
				}
			}
		}
		_, err := parseDetailPage(strings.TrimSpace(href))
		if err != nil {
			log.Errorf("get detail page: %v", err)
			return
		}

		item := importlist.Item{
			Title: strings.TrimSpace(link.Text()),
			Year:  year,
		}
		items = append(items, item)
		_ = item
		//println(link.Text(), href)
	})

	return &importlist.Response{Items: items}, nil
}

func parseDetailPage(url string) (string, error) {
	println(url)

	res, err := doHttpReq("GET", url, nil)
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

	doc.Find("div[class='subject clearfix']").Each(func(i int, se *goquery.Selection) {
		println(se.Text())
		se.Children().Get(1)
		imdb := se.Find("div[class='info']").First().Children().Last()
		println(imdb.Text())
	})

	_ = doc
	return "", nil
}
func NewDoubanWishlist(personId string) *DoubanWishlist {
	return &DoubanWishlist{PersonId: personId}
}

type DoubanWishlist struct {
	PersonId string
}

const wishlistUrl = "https://movie.douban.com/people/%s/wish?sort=time&start=%d&mode=grid&tags_sort=count"

func (d *DoubanWishlist) GetWishlist(page int) (*importlist.Response, error) {
	c := colly.NewCollector(colly.UserAgent(ua))
	c.Limit(&colly.LimitRule{
		DomainRegexp: "*",
		Delay:       10 * time.Second,
		RandomDelay: 2 * time.Second,
	})
	url := fmt.Sprintf(wishlistUrl, d.PersonId, (page-1)*15)
	c.OnHTML("div[class='item comment-item']", func(e *colly.HTMLElement) {
		if !strings.HasPrefix(e.Request.URL.String(), "https://movie.douban.com/people") {
			return
		}
		e.DOM.Find("div[class='pic'] a[title]").Each(func(i int, selection *goquery.Selection) {
			println(selection.Attr("href"))
			url, ok := selection.Attr("href")
			if ok {
				c.Visit(url)
			}
		})
	})

	c.OnHTML("#content", func(h *colly.HTMLElement) {
		var item importlist.Item
		h.DOM.Find("h1").Each(func(i int, selection *goquery.Selection) {
			selection.Find("span[property]").Each(func(i int, selection *goquery.Selection) {
				println(selection.Text())
				item.Title = selection.Text()
			})
			selection.Find("span[class='year']").Each(func(i int, selection *goquery.Selection) {
				n, _ := strconv.Atoi(selection.Text())
				item.Year = n
			})

		})
		h.DOM.Find("#info").Each(func(i int, s *goquery.Selection) {
			info := strings.TrimSpace(s.Text()) 
			lines := strings.Split(info, "\n")
			if len(lines) == 0 {
				return
			}
			last := lines[len(lines)-1]
			if !strings.HasPrefix(strings.ToLower(last), "imdb") {
				return
			}
			ss := strings.Split(last, ":")
			for _, p := range ss {
				p := strings.TrimSpace(strings.ToLower(p))
				if strings.HasPrefix(p, "tt") {
					item.ImdbID = p
				}
			}
		})
		log.Info(item)
	})

	return nil, c.Visit(url)
}

func doHttpReq(method, url string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", ua)
	return http.DefaultClient.Do(req)
}
