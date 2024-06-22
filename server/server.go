package server

import (
	"polaris/db"
	"polaris/log"
	"polaris/pkg/tmdb"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewServer(db *db.Client) *Server {
	r := gin.Default()
	return &Server{
		r: r,
		db: db,
	}
}

type Server struct {
	r    *gin.Engine
	db *db.Client
	language string
}

func (s *Server) Serve() error {
	api := s.r.Group("/api/v1")

	setting := api.Group("/setting")
	{
		setting.POST("/do", HttpHandler(s.SetSetting))
		setting.GET("/do", HttpHandler(s.GetSetting))
	}

	tv := api.Group("/tv")
	{
		tv.GET("/search", HttpHandler(s.SearchTvSeries))
		tv.POST("/watchlist", HttpHandler(s.AddWatchlist))
		tv.GET("/watchlist", HttpHandler(s.GetWatchlist))
	}

	s.language = s.db.GetLanguage()
	return s.r.Run(":8080")
}


func (s *Server) TMDB() (*tmdb.Client, error) {
	api := s.db.GetSetting(db.SettingTmdbApiKey)
	if api == "" {
		return nil, errors.New("tmdb api not set")
	}
	return tmdb.NewClient(api)
}

func (s *Server) MustTMDB() *tmdb.Client {
	t, err := s.TMDB()
	if err != nil {
		log.Panicf("get tmdb: %v", err)
	}
	return t
}

