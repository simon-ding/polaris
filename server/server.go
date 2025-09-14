package server

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"polaris/db"
	"polaris/engine"
	"polaris/log"
	"polaris/pkg/cache"
	"polaris/pkg/tmdb"
	"polaris/ui"
	"strconv"
	"time"

	ginzap "github.com/gin-contrib/zap"

	"github.com/gin-contrib/static"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewServer(db db.Database) *Server {
	s := &Server{
		db:               db,
		srv: &http.Server{},
		language:         db.GetLanguage(),
		monitorNumCache:  cache.NewCache[int, int](10 * time.Minute),
		downloadNumCache: cache.NewCache[int, int](10 * time.Minute),
	}
	s.core = engine.NewEngine(db, s.language)
	s.setupRoutes()
	return s
}

type Server struct {
	srv              *http.Server
	db               db.Database
	core             *engine.Engine
	language         string
	jwtSerect        string
	monitorNumCache  *cache.Cache[int, int]
	downloadNumCache *cache.Cache[int, int]
}

func (s *Server) setupRoutes() {
	s.core.Init()

	r := gin.Default()
	s.jwtSerect = s.db.GetSetting(db.JwtSerectKey)
	//st, _ := fs.Sub(ui.Web, "build/web")
	fs, err := static.EmbedFolder(ui.Web, "build/web")
	if err == nil {
		r.Use(static.Serve("/", fs))
	} else {
		log.Warnf("serve web static files error: %v", err)
	}
	
	//s.r.Use(ginzap.Ginzap(log.Logger().Desugar(), time.RFC3339, false))
	r.Use(ginzap.RecoveryWithZap(log.Logger().Desugar(), true))

	log.SetLogLevel(s.db.GetSetting(db.SettingLogLevel)) //restore log level

	r.POST("/api/login", HttpHandler(s.Login))

	api := r.Group("/api/v1")
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
		setting.POST("/monitoring", HttpHandler(s.ChangeEpisodeMonitoring))
		setting.POST("/cron/trigger", HttpHandler(s.TriggerCronJob))
		setting.GET("/prowlarr", HttpHandler(s.GetProwlarrSetting))
		setting.POST("/prowlarr", HttpHandler(s.SaveProwlarrSetting))
		setting.GET("/limiter", HttpHandler(s.GetSizeLimiter))
		setting.POST("/limiter", HttpHandler(s.SetSizeLimiter))
	}
	activity := api.Group("/activity")
	{
		activity.GET("/", HttpHandler(s.GetAllActivities))
		activity.POST("/delete", HttpHandler(s.RemoveActivity))
		activity.GET("/media/:id", HttpHandler(s.GetMediaDownloadHistory))
		activity.GET("/blacklist", HttpHandler(s.GetAllBlacklistItems))
		activity.DELETE("/blacklist/:id", HttpHandler(s.RemoveBlacklistItem))
		//activity.GET("/torrents", HttpHandler(s.GetAllTorrents))
	}

	tv := api.Group("/media")
	{
		tv.GET("/search", HttpHandler(s.SearchMedia))
		tv.POST("/edit", HttpHandler(s.EditMediaMetadata))
		tv.POST("/tv/watchlist", HttpHandler(s.AddTv2Watchlist))
		tv.GET("/tv/watchlist", HttpHandler(s.GetTvWatchlist))
		tv.POST("/torrents", HttpHandler(s.SearchAvailableTorrents))
		tv.POST("/torrents/download", HttpHandler(s.DownloadTorrent))
		tv.POST("/movie/watchlist", HttpHandler(s.AddMovie2Watchlist))
		tv.GET("/movie/watchlist", HttpHandler(s.GetMovieWatchlist))
		tv.GET("/record/:id", HttpHandler(s.GetMediaDetails))
		tv.DELETE("/record/:id", HttpHandler(s.DeleteFromWatchlist))
		tv.GET("/suggest/tv/:tmdb_id", HttpHandler(s.SuggestedSeriesFolderName))
		tv.GET("/suggest/movie/:tmdb_id", HttpHandler(s.SuggestedMovieFolderName))
		tv.GET("/downloadall/:id", HttpHandler(s.DownloadAll))
		tv.GET("/download/tv", HttpHandler(s.DownloadAllTv))
		tv.GET("/download/movie", HttpHandler(s.DownloadAllMovies))
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
	importlist := api.Group("/importlist")
	{
		importlist.GET("/", HttpHandler(s.getAllImportLists))
		importlist.POST("/add", HttpHandler(s.addImportlist))
		importlist.DELETE("/delete", HttpHandler(s.deleteImportList))
	}
	s.srv.Handler = r

}

func (s *Server) Start(addr string) (int, error) {
	if addr == "" {
		addr = "127.0.0.1:0" // 0 means any available port
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("failed to listen on port: %w", err)
	}

	_, port, _ := net.SplitHostPort(ln.Addr().String())

	p, err := strconv.Atoi(port)
	if err != nil {
		return 0, fmt.Errorf("failed to convert port to int: %w", err)
	}
	go func() {
		defer ln.Close()
		if err := s.srv.Serve(ln); err != nil {
			log.Errorf("failed to serve: %v", err)
		}
	}()

	log.Infof("----------- Polaris Server Successfully Started on Port %d------------", p)

	return p, nil
}

func (s *Server) Stop() error {
	log.Infof("Stopping Polaris Server...")
	return s.srv.Close()
}

func (s *Server) TMDB() (*tmdb.Client, error) {
	api := s.db.GetTmdbApiKey()
	if api == "" {
		return nil, errors.New("TMDB apiKey not set")
	}
	proxy := s.db.GetSetting(db.SettingProxy)
	adult := s.db.GetSetting(db.SettingEnableTmdbAdultContent)
	return tmdb.NewClient(api, proxy, adult == "true")
}

func (s *Server) MustTMDB() *tmdb.Client {
	t, err := s.TMDB()
	if err != nil {
		log.Panicf("get tmdb: %v", err)
	}
	return t
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
