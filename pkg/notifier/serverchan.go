package notifier

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type ServerChanConfig struct {
	Key string `json:"key"`
}

func NewServerChanClient(s string) (NotificationClient, error) {
	var cfg ServerChanConfig
	if err := json.Unmarshal([]byte(s), &cfg); err != nil {
		return nil, errors.Wrap(err, "json")
	}

	return &ServerChan{Key: cfg.Key}, nil
}

type ServerChan struct {
	Key string
}

func (s *ServerChan) SendMsg(msg string) error {
	return scSend("Polaris", msg, s.Key)
}

func scSend(text string, desp string, key string) error {
	data := url.Values{}
	data.Set("text", text)
	data.Set("desp", desp)

	// 根据 sendkey 是否以 "sctp" 开头决定 API 的 URL
	var apiUrl string
	if strings.HasPrefix(key, "sctp") {
		// 使用正则表达式提取数字部分
		re := regexp.MustCompile(`sctp(\d+)t`)
		matches := re.FindStringSubmatch(key)
		if len(matches) > 1 {
			num := matches[1]
			apiUrl = fmt.Sprintf("https://%s.push.ft07.com/send/%s.send", num, key)
		} else {
			return errors.New("invalid sendkey format for sctp")
		}
	} else {
		apiUrl = fmt.Sprintf("https://sctapi.ftqq.com/%s.send", key)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	d, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var r response
	if err := json.Unmarshal(d, &r); err != nil {
		return err
	}

	if r.Code != 0 {
		return errors.New(r.Message)
	}

	return nil
}

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
