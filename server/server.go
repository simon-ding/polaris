package server

import (
	"polaris/db"
	"polaris/log"
	"polaris/pkg/tmdb"
	"polaris/pkg/transmission"
	"polaris/ui"

	"github.com/gin-contrib/static"
	"github.com/robfig/cron"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewServer(db *db.Client) *Server {
	r := gin.Default()
	return &Server{
		r:     r,
		db:    db,
		cron:  cron.New(),
		tasks: make(map[int]*Task),
	}
}

type Server struct {
	r        *gin.Engine
	db       *db.Client
	cron     *cron.Cron
	language string
	tasks    map[int]*Task
}

func (s *Server) Serve() error {
	s.scheduler()
	s.reloadTasks()
	//st, _ := fs.Sub(ui.Web, "build/web")
	s.r.Use(static.Serve("/", static.EmbedFolder(ui.Web, "build/web")))

	s.r.POST("/api/login", HttpHandler(s.Login))

	api := s.r.Group("/api/v1")
	api.Use(s.authModdleware)

	setting := api.Group("/setting")
	{
		setting.POST("/do", HttpHandler(s.SetSetting))
		setting.GET("/do", HttpHandler(s.GetSetting))
		setting.POST("/auth", HttpHandler(s.EnableAuth))
		setting.GET("/auth", HttpHandler(s.GetAuthSetting))
	}
	activity := api.Group("/activity")
	{
		activity.GET("/", HttpHandler(s.GetAllActivities))
		activity.DELETE("/:id", HttpHandler(s.RemoveActivity))
	}

	tv := api.Group("/tv")
	{
		tv.GET("/search", HttpHandler(s.SearchTvSeries))
		tv.POST("/watchlist", HttpHandler(s.AddWatchlist))
		tv.GET("/watchlist", HttpHandler(s.GetWatchlist))
		tv.GET("/series/:id", HttpHandler(s.GetTvDetails))
		tv.GET("/resolutions", HttpHandler(s.GetAvailableResolutions))
	}
	indexer := api.Group("/indexer")
	{
		indexer.GET("/", HttpHandler(s.GetAllIndexers))
		indexer.POST("/add", HttpHandler(s.AddTorznabInfo))
		indexer.POST("/download", HttpHandler(s.SearchAndDownload))
		indexer.DELETE("/del/:id", HttpHandler(s.DeleteTorznabInfo))
	}

	downloader := api.Group("/downloader")
	{
		downloader.GET("/", HttpHandler(s.GetAllDonloadClients))
		downloader.POST("/add", HttpHandler(s.AddDownloadClient))
		downloader.DELETE("/del/:id", HttpHandler(s.DeleteDownloadCLient))
	}
	storage := api.Group("/storage")
	{
		storage.GET("/", HttpHandler(s.GetAllStorage))
		storage.POST("/", HttpHandler(s.AddStorage))
		storage.DELETE("/:id", HttpHandler(s.DeleteStorage))
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

func (s *Server) reloadTasks() {
	runningTasks := s.db.GetRunningHistories()
	if len(runningTasks) == 0 {
		return
	}
	for _, t := range runningTasks {
		log.Infof("reloading task: %d %s", t.ID, t.SourceTitle)
		torrent, err := transmission.ReloadTorrent(t.Saved)
		if err != nil {
			log.Errorf("relaod task %s failed: %v", t.SourceTitle, err)
			continue
		}
		s.tasks[t.ID] = &Task{Torrent: torrent}
	}
}
