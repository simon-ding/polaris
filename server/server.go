package server

import (
	"polaris/db"
	"polaris/log"
	"polaris/pkg/tmdb"

	"github.com/hekmon/transmissionrpc"
	"github.com/robfig/cron"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewServer(db *db.Client) *Server {
	r := gin.Default()
	return &Server{
		r:    r,
		db:   db,
		cron: cron.New(),
		tasks: make(map[string]*transmissionrpc.Torrent),
	}
}

type Server struct {
	r        *gin.Engine
	db       *db.Client
	cron     *cron.Cron
	language string
	tasks map[string]*transmissionrpc.Torrent
}

func (s *Server) scheduler() {
	s.cron.AddFunc("@every 1m", s.checkTasks)
}

func (s *Server) checkTasks() {
	for name, t := range s.tasks {
		log.Infof("task %s percentage done: %f", name, *t.PercentDone)
	}
}

func (s *Server) Serve() error {
	s.scheduler()

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
	indexer := api.Group("/indexer")
	{
		indexer.POST("/add", HttpHandler(s.AddTorznabInfo))
		indexer.POST("/download", HttpHandler(s.SearchAndDownload))
	}

	downloader := api.Group("/downloader")
	{
		downloader.POST("/add", HttpHandler(s.AddDownloadClient))
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
