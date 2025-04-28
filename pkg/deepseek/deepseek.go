package deepseek

import (
	"context"
	"fmt"
	"polaris/log"
	"time"

	"github.com/goccy/go-json"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func NewClient(apiKey string) *Client {
	r := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL("https://api.deepseek.com"))
	return &Client{openai: &r, model: "deepseek-chat"}
}

type Client struct {
	openai *openai.Client
	model  string
}

func (c *Client) Test() error {

	question := `What computer ran the first neural network?
	EXAMPLE JSON OUTPUT:
	{
		"origin": "The origin of the computer",
		"full_name": "The name of the device model",
		"legacy": "Its influence on the field of computing",
		"notable_facts": "A few key facts about the computer
	}
	`

	chat, err_ := c.openai.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		// ...
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				Type: "json_object",
			},
		},
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},

		// only certain models can perform structured outputs
		Model: c.model,
	})
	if err_ != nil {
		return err_
	}

	log.Infof("%+v", chat.Choices[0].Message.Content)
	// extract into a well-typed struct
	return nil
}

type Movies struct {
	Movies []struct {
		Name  string `json:"name"`
		Match string `json:"match"`
	} `json:"movies"`
}

type Tvs struct {
	Tvs []Tv `json:"tvs"`
}

type Tv struct {
	Name            string `json:"name"`
	FileName        string `json:"file_name"`
	Match           string `json:"match"`
	Season          string `json:"season"`
	StartEpisode    string `json:"start_episode"`
	EndEpisode      string `json:"end_episode"`
	Quality         string `json:"quality"`
	IsSeasonPackage string `json:"is_season_package"`
}

func (c *Client) AssessMovieNames(movieName string, releaseYear int, torrentNames []string) (*Movies, error) {
	q := `用户输入的是一些文件名称，你需要判断哪些文件可能属于 %d 年的电影 %s，哪些可能不是。

	EXAMPLE JSON OUTPUT:
	{
		"movies": [
			{
				"name": "The name of the movie",
				"match": "true or false"
			},
		]
	}
	`

	q = fmt.Sprintf(q, releaseYear, movieName)
	chat, err_ := c.openai.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		//
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				Type: "json_object",
			},
		},
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(q),
			openai.UserMessage(fmt.Sprintf("文件名称: %v", torrentNames)),
		},

		// only certain models can perform structured outputs
		Model: c.model,
	})
	if err_ != nil {
		return nil, err_
	}

	log.Infof("%+v", chat.Choices[0].Message.Content)
	var res Movies
	if err := json.Unmarshal([]byte(chat.Choices[0].Message.Content), &res); err != nil {
		return nil, err
	}

	// extract into a well-typed struct
	return &res, nil
}

func (c *Client) AssessTvNames(tvName string, releaseYear int, torrentNames []string) ([]Tv, error) {
	log.Debugf("deepseek tv name: %s, year: %d, torrent name len: %v", tvName, releaseYear, len(torrentNames))
	t := time.Now()
	defer func() {
		log.Infof("deepseek assess tv name cost: %v", time.Since(t))
	}()

	q := `用户输入的是一些文件名称，你需要判断哪些文件可能属于 %d 年的电视剧 %s，哪些可能不是，并返回匹配的文件名。

	EXAMPLE JSON OUTPUT:	
	{
		"tvs": [
		"matched file name 1", "matched file name 2", ...
		]	
	}`
	q = fmt.Sprintf(q, releaseYear, tvName)

	var res []Tv

	chat, err_ := c.openai.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		MaxTokens: openai.Opt(int64(4096)),
		//...
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				Type: "json_object",
			},
		},
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(q),
			openai.UserMessage(fmt.Sprintf("文件名称: %v", torrentNames)),
		},

		// only certain models can perform structured outputs
		Model: c.model,
	})
	if err_ != nil {
		return nil, err_
	}
	log.Infof("%+v", chat.Choices[0].Message.Content)
	var tvs Tvs
	if err := json.Unmarshal([]byte(chat.Choices[0].Message.Content), &tvs); err != nil {
		return nil, err
	}
	res = append(res, tvs.Tvs...)

	return res, nil
}
