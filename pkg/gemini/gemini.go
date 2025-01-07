package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"polaris/log"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func NewClient(apiKey, modelName string) (*Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &Client{apiKey: apiKey, modelName: modelName, c: client}, nil

}

type Client struct {
	apiKey    string
	modelName string
	c         *genai.Client
}

type TvInfo struct {
	TitleEnglish     string `json:"title_english"`
	TitleChinses     string `json:"title_chinese"`
	Season           int    `json:"season"`
	StartEpisode     int    `json:"start_episode"`
	EndEpisode       int    `json:"end_episode"`
	Resolution       string `json:"resolution"`
	Subtitle         string `json:"subtitle"`
	ReleaseGroup     string `json:"release_group"`
	Year             int    `json:"year"`
	AudioLanguage    string `json:"audio_language"`
	IsCompleteSeason bool   `json:"is_complete_season"`
}

func (c *Client) ParseTvInfo(q string) (*TvInfo, error) {
	log.Info(q)
	ctx := context.Background()

	model := c.c.GenerativeModel(c.modelName)

	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title_english": {Type: genai.TypeString},
			"title_chinese": {Type: genai.TypeString},
			"season":        {Type: genai.TypeInteger, Description: "season number"},
			"start_episode": {Type: genai.TypeInteger},
			"end_episode":   {Type: genai.TypeInteger},
			//"episodes":       {Type: genai.TypeString},
			"resolution":         {Type: genai.TypeString},
			"subtitle":           {Type: genai.TypeString},
			"release_group":      {Type: genai.TypeString},
			"year":               {Type: genai.TypeInteger},
			"audio_language":     {Type: genai.TypeString},
			"is_complete_season": {Type: genai.TypeBoolean},
		},
		Required: []string{"title_english", "title_chinese", "season", "start_episode", "resolution"},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(q))
	if err != nil {
		return nil, err
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			var info TvInfo
			if err := json.Unmarshal([]byte(txt), &info); err != nil {
				return nil, err
			}
			return &info, nil
		}
	}
	return nil, fmt.Errorf("not found")

}

type MovieInfo struct {
	TitleEnglish     string `json:"title_english"`
	TitleChinses     string `json:"title_chinese"`
	Resolution       string `json:"resolution"`
	Subtitle         string `json:"subtitle"`
	ReleaseGroup     string `json:"release_group"`
	Year             int    `json:"year"`
	AudioLanguage    string `json:"audio_language"`
}

func (c *Client) ParseMovieInfo(q string) (*MovieInfo, error) {
	log.Info(q)
	ctx := context.Background()

	model := c.c.GenerativeModel(c.modelName)

	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title_english":  {Type: genai.TypeString},
			"title_chinese":  {Type: genai.TypeString},
			"resolution":     {Type: genai.TypeString},
			"subtitle":       {Type: genai.TypeString},
			"release_group":  {Type: genai.TypeString},
			"year":           {Type: genai.TypeInteger},
			"audio_language": {Type: genai.TypeString},
		},
		Required: []string{"title_english", "title_chinese", "resolution"},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(q))
	if err != nil {
		return nil, err
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			var info MovieInfo
			if err := json.Unmarshal([]byte(txt), &info); err != nil {
				return nil, err
			}
			return &info, nil
		}
	}
	return nil, fmt.Errorf("not found")

}

func (c *Client) isTvSeries(q string) (bool, error) {
	ctx := context.Background()

	model := c.c.GenerativeModel(c.modelName)

	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeBoolean, Nullable: true, Description: "whether the input text implies a tv series",
	}

	resp, err := model.GenerateContent(ctx, genai.Text(q))
	if err != nil {
		return false, err
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			return strings.ToLower(string(txt)) == "true", nil
		}
	}
	return false, fmt.Errorf("error")
}

func (c *Client) ImpliesSameTvOrMovie(torrentName, mediaName string) bool {
	ctx := context.Background()

	model := c.c.GenerativeModel(c.modelName)

	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type: genai.TypeBoolean, Nullable: true,
	}
	q := fmt.Sprintf("whether this file name \"%s\" implies the same TV series or movie with name \"%s\"?", torrentName, mediaName)
	resp, err := model.GenerateContent(ctx, genai.Text(q))
	if err != nil {
		return false
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			return strings.ToLower(string(txt)) == "true"
		}
	}
	return false

}

func (c *Client) FilterTvOrMovies(resourcesNames []string, titles ...string) ([]string, error) {
	ctx := context.Background()

	model := c.c.GenerativeModel(c.modelName)

	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type:  genai.TypeArray,
		Items: &genai.Schema{Type: genai.TypeString},
	}
	for i, s := range titles {
		titles[i] = "\"" + s + "\""
	}
	p := &bytes.Buffer{}
	p.WriteString(`the following list of file names, list all of which implies the same TV series or movie of name`)
	p.WriteString(strings.Join(titles, " or "))
	p.WriteString(":\n")

	for _, r := range resourcesNames {
		p.WriteString(" * ")
		p.WriteString(r)
		p.WriteString("\n")
	}
	log.Debugf("FilterTvOrMovies prompt is %s", p.String())

	resp, err := model.GenerateContent(ctx, genai.Text(p.String()))
	if err != nil {
		return nil, err
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {

			var s []string
			if err := json.Unmarshal([]byte(txt), &s); err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	return nil, fmt.Errorf("nothing found")

}


func (c *Client) FilterMovies(resourcesNames []string, year int, titles ...string) ([]string, error) {
	ctx := context.Background()

	model := c.c.GenerativeModel(c.modelName)

	model.ResponseMIMEType = "application/json"
	model.ResponseSchema = &genai.Schema{
		Type:  genai.TypeArray,
		Items: &genai.Schema{Type: genai.TypeString},
	}
	for i, s := range titles {
		titles[i] = "\"" + s + "\""
	}
	p := &bytes.Buffer{}
	p.WriteString( fmt.Sprint("the following list of file names, list all of which match following criteria: 1. Is movie 2. Released in year %d 3. Have name of ", year))
	p.WriteString(strings.Join(titles, " or "))
	p.WriteString(":\n")

	for _, r := range resourcesNames {
		p.WriteString(" * ")
		p.WriteString(r)
		p.WriteString("\n")
	}
	log.Debugf("FilterTvOrMovies prompt is %s", p.String())

	resp, err := model.GenerateContent(ctx, genai.Text(p.String()))
	if err != nil {
		return nil, err
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {

			var s []string
			if err := json.Unmarshal([]byte(txt), &s); err != nil {
				return nil, err
			}
			return s, nil
		}
	}
	return nil, fmt.Errorf("nothing found")

}
