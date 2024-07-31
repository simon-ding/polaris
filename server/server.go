package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"polaris/db"
	"polaris/log"
	"polaris/pkg/tmdb"
	"polaris/pkg/transmission"
	"polaris/ui"

	ginzap "github.com/gin-contrib/zap"

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
	r         *gin.Engine
	db        *db.Client
	cron      *cron.Cron
	language  string
	tasks     map[int]*Task
	jwtSerect string
}

func (s *Server) Serve() error {
	s.scheduler()
	s.reloadTasks()
	s.restoreProxy()

	s.jwtSerect = s.db.GetSetting(db.JwtSerectKey)
	//st, _ := fs.Sub(ui.Web, "build/web")
	s.r.Use(static.Serve("/", static.EmbedFolder(ui.Web, "build/web")))
	//s.r.Use(ginzap.Ginzap(log.Logger().Desugar(), time.RFC3339, false))
	s.r.Use(ginzap.RecoveryWithZap(log.Logger().Desugar(), true))

	log.SetLogLevel(s.db.GetSetting(db.SettingLogLevel)) //restore log level

	s.r.POST("/api/login", HttpHandler(s.Login))

	api := s.r.Group("/api/v1")
	api.Use(s.authModdleware)
	api.StaticFS("/img", http.Dir(db.ImgPath))
	api.StaticFS("/logs", http.Dir(db.LogPath))
	api.Any("/posters/*proxyPath", s.proxyPosters)

	setting := api.Group("/setting")
	{
		setting.GET("/logout", HttpHandler(s.Logout))
		setting.POST("/general", HttpHandler(s.SetSetting))
		setting.GET("/general", HttpHandler(s.GetSetting))
		setting.POST("/auth", HttpHandler(s.EnableAuth))
		setting.GET("/auth", HttpHandler(s.GetAuthSetting))
		setting.GET("/logfiles", HttpHandler(s.GetAllLogs))
		setting.GET("/about", HttpHandler(s.About))
		setting.POST("/parse/tv", HttpHandler(s.ParseTv))
		setting.POST("/parse/movie", HttpHandler(s.ParseMovie))
	}
	activity := api.Group("/activity")
	{
		activity.GET("/", HttpHandler(s.GetAllActivities))
		activity.DELETE("/:id", HttpHandler(s.RemoveActivity))
		activity.GET("/media/:id", HttpHandler(s.GetMediaDownloadHistory))
		activity.GET("/torrents", HttpHandler(s.GetAllTorrents))
	}

	tv := api.Group("/media")
	{
		tv.GET("/search", HttpHandler(s.SearchMedia))
		tv.POST("/tv/watchlist", HttpHandler(s.AddTv2Watchlist))
		tv.GET("/tv/watchlist", HttpHandler(s.GetTvWatchlist))
		tv.POST("/torrents", HttpHandler(s.SearchAvailableTorrents))
		tv.POST("/torrents/download/", HttpHandler(s.DownloadTorrent))
		tv.POST("/movie/watchlist", HttpHandler(s.AddMovie2Watchlist))
		tv.GET("/movie/watchlist", HttpHandler(s.GetMovieWatchlist))
		tv.GET("/record/:id", HttpHandler(s.GetMediaDetails))
		tv.DELETE("/record/:id", HttpHandler(s.DeleteFromWatchlist))
		tv.GET("/suggest/tv/:tmdb_id", HttpHandler(s.SuggestedSeriesFolderName))
		tv.GET("/suggest/movie/:tmdb_id", HttpHandler(s.SuggestedMovieFolderName))
	}
	indexer := api.Group("/indexer")
	{
		indexer.GET("/", HttpHandler(s.GetAllIndexers))
		indexer.POST("/add", HttpHandler(s.AddTorznabInfo))
		indexer.POST("/download", HttpHandler(s.SearchTvAndDownload))
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
	notifier := api.Group("/notifier")
	{
		notifier.GET("/all", HttpHandler(s.GetAllNotificationClients))
		notifier.GET("/id/:id", HttpHandler(s.GetNotificationClient))
		notifier.DELETE("/id/:id", HttpHandler(s.DeleteNotificationClient))
		notifier.POST("/add", HttpHandler(s.AddNotificationClient))
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

func (s *Server) proxyPosters(c *gin.Context) {
	remote, _ := url.Parse("https://image.tmdb.org")
	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = fmt.Sprintf("/t/p/w500/%v", c.Param("proxyPath"))
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
