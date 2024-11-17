package alist

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"polaris/log"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Resposne[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type Config struct {
	Username string
	Password string
	URL      string
}

func New(cfg *Config) *Client {
	cfg.URL = strings.Trim(cfg.URL, "/")
	return &Client{
		cfg:  cfg,
		http: http.DefaultClient,
	}
}

type Client struct {
	cfg   *Config
	http  *http.Client
	token string
}

func (c *Client) Login() (string, error) {
	p := map[string]string{
		"username": c.cfg.Username,
		"password": c.cfg.Password,
	}
	data, _ := json.Marshal(p)
	resp, err := c.http.Post(c.cfg.URL+loginUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", errors.Wrap(err, "login")
	}
	defer resp.Body.Close()
	d1, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read body")
	}
	var rp Resposne[map[string]string]

	err = json.Unmarshal(d1, &rp)
	if err != nil {
		return "", errors.Wrap(err, "json")
	}
	if rp.Code != 200 {
		return "", errors.Errorf("alist error: code %d, %s", rp.Code, rp.Message)
	}
	c.token = rp.Data["token"]
	return c.token, nil
}

type LsInfo struct {
	Content []struct {
		Name     string    `json:"name"`
		Size     int       `json:"size"`
		IsDir    bool      `json:"is_dir"`
		Modified time.Time `json:"modified"`
		Created  time.Time `json:"created"`
		Sign     string    `json:"sign"`
		Thumb    string    `json:"thumb"`
		Type     int       `json:"type"`
		Hashinfo string    `json:"hashinfo"`
		HashInfo any       `json:"hash_info"`
	} `json:"content"`
	Total    int    `json:"total"`
	Readme   string `json:"readme"`
	Header   string `json:"header"`
	Write    bool   `json:"write"`
	Provider string `json:"provider"`
}

func (c *Client) Ls(dir string) (*LsInfo, error) {
	in := map[string]string{
		"path": dir,
	}

	resp, err := c.post(c.cfg.URL+lsUrl, in)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}

	var out Resposne[LsInfo]
	err = json.Unmarshal(resp, &out)
	if err != nil {
		return nil, err
	}
	if out.Code != 200 {
		return nil, errors.Errorf("alist error: code %d, %s", out.Code, out.Message)
	}
	return &out.Data, nil
}

func (c *Client) Mkdir(dir string) error {
	in := map[string]string{
		"path": dir,
	}
	resp, err := c.post(c.cfg.URL+mkdirUrl, in)
	if err != nil {
		return errors.Wrap(err, "http")
	}
	var out Resposne[any]
	err = json.Unmarshal(resp, &out)
	if err != nil {
		return err
	}
	if out.Code != 200 {
		return errors.Errorf("alist error: code %d, %s", out.Code, out.Message)
	}
	return nil
}

func (c *Client) post(url string, body interface{}) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	req.Header.Add("Authorization", c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http")
	}
	defer resp.Body.Close()
	d1, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body")
	}
	return d1, nil
}

type UploadStreamResponse struct {
	Task struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		State    int    `json:"state"`
		Status   string `json:"status"`
		Progress int    `json:"progress"`
		Error    string `json:"error"`
	} `json:"task"`
}

func (c *Client) UploadStream(reader io.Reader, size int64, toDir string) (*UploadStreamResponse, error) {
	req, err := http.NewRequest(http.MethodPut, c.cfg.URL+streamUploadUrl, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.token)
	req.Header.Add("File-Path", url.PathEscape(toDir))
	//req.Header.Add("As-Task", "true")
	req.Header.Add("Content-Type", "application/octet-stream")
	req.ContentLength = size

	log.Infof("headers: %+v, %v", req.Header, req.URL.String())
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	d1, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var out Resposne[UploadStreamResponse]
	err = json.Unmarshal(d1, &out)
	if err != nil {
		return nil, err
	}
	if out.Code != 200 {
		return nil, errors.Errorf("alist error: code %d, %s", out.Code, out.Message)
	}

	return &out.Data, nil
}
