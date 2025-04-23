package engine

import (
	"polaris/db"
	"polaris/ent"
	"polaris/ent/downloadclients"
	"polaris/log"
	"polaris/pkg"
	"polaris/pkg/buildin"
	"polaris/pkg/qbittorrent"
	"polaris/pkg/tmdb"
	"polaris/pkg/transmission"
	"polaris/pkg/utils"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

func NewEngine(db db.Database, language string) *Engine {
	return &Engine{
		db:         db,
		cron:       cron.New(),
		tasks:      utils.Map[int, *Task]{},
		schedulers: utils.Map[string, scheduler]{},
		language:   language,
	}
}

type scheduler struct {
	cron string
	f    func() error
}
type Engine struct {
	db         db.Database
	cron       *cron.Cron
	tasks      utils.Map[int, *Task]
	language   string
	schedulers utils.Map[string, scheduler]
	buildin    *buildin.Downloader
}

func (c *Engine) registerCronJob(name string, cron string, f func() error) {
	c.schedulers.Store(name, scheduler{
		cron: cron,
		f:    f,
	})
}

func (c *Engine) Init() {
	go c.reloadTasks()
	c.addSysCron()
	go c.checkW500PosterOnStartup()
}

func (c *Engine) GetTask(id int) (*Task, bool) {
	return c.tasks.Load(id)
}

func (c *Engine) reloadUsingBuildinDownloader(h *ent.History) error {
	cl, err := buildin.NewDownloader(c.db.GetDownloadDir())
	if err != nil {
		log.Warnf("buildin downloader error: %v", err)
	}
	t, err := cl.Download(h.Link, h.Hash, c.db.GetDownloadDir())
	if err != nil {
		return errors.Wrap(err, "download torrent")
	}
	c.tasks.Store(h.ID, &Task{Torrent: t})
	return nil
}

func (c *Engine) reloadTasks() {
	allTasks := c.db.GetRunningHistories()
	for _, t := range allTasks {
		if t.DownloadClientID == 0 {
			log.Warnf("assume buildin downloader: %v", t.SourceTitle)
			err := c.reloadUsingBuildinDownloader(t)
			if err != nil {
				log.Warnf("buildin downloader error: %v", err)
			} else {
				log.Infof("success reloading buildin task: %v", t.SourceTitle)
			}
			continue
		}
		dl, err := c.db.GetDownloadClient(t.DownloadClientID)
		if err != nil {
			log.Warnf("no download client related: %v", t.SourceTitle)
			continue
		}

		if dl.Implementation == downloadclients.ImplementationTransmission {
			if t.Hash != "" { //优先使用hash
				to, err := transmission.NewTorrentHash(transmission.Config{
					URL:      dl.URL,
					User:     dl.User,
					Password: dl.Password,
				}, t.Hash)
				if err != nil {
					log.Warnf("get task error: %v", err)
					continue
				}
				c.tasks.Store(t.ID, &Task{Torrent: to})
			} else if t.Link != "" {
				to, err := transmission.NewTorrent(transmission.Config{
					URL:      dl.URL,
					User:     dl.User,
					Password: dl.Password,
				}, t.Link)
				if err != nil {
					log.Warnf("get task error: %v", err)
					continue
				}
				c.tasks.Store(t.ID, &Task{Torrent: to})
			}
		} else if dl.Implementation == downloadclients.ImplementationQbittorrent {
			if t.Hash != "" {
				to, err := qbittorrent.NewTorrentHash(qbittorrent.Info{
					URL:      dl.URL,
					User:     dl.User,
					Password: dl.Password,
				}, t.Hash)
				if err != nil {
					log.Warnf("get task error: %v", err)
					continue
				}
				c.tasks.Store(t.ID, &Task{Torrent: to})

			} else if t.Link != "" {
				to, err := qbittorrent.NewTorrent(qbittorrent.Info{
					URL:      dl.URL,
					User:     dl.User,
					Password: dl.Password,
				}, t.Link)
				if err != nil {
					log.Warnf("get task error: %v", err)
					continue
				}
				c.tasks.Store(t.ID, &Task{Torrent: to})
			}
		}

	}
	log.Infof("------ task reloading done ------")
}

func (c *Engine) buildInDownloader() (pkg.Downloader, error) {
	if c.buildin != nil {
		return c.buildin, nil
	}
	dir := c.db.GetDownloadDir()
	d, err := buildin.NewDownloader(dir)
	if err != nil {
		return nil, errors.Wrap(err, "buildin downloader")
	}
	c.buildin = d
	return d, nil
}

func (c *Engine) GetDownloadClient() (pkg.Downloader, *ent.DownloadClients, error) {
	downloaders := c.db.GetAllDonloadClients()
	for _, d := range downloaders {
		if !d.Enable {
			continue
		}
		if d.Implementation == downloadclients.ImplementationTransmission {
			trc, err := transmission.NewClient(transmission.Config{
				URL:      d.URL,
				User:     d.User,
				Password: d.Password,
			})
			if err != nil {
				log.Warnf("connect to download client error: %v", d.URL)
				continue
			}
			return trc, d, nil

		} else if d.Implementation == downloadclients.ImplementationQbittorrent {
			qbt, err := qbittorrent.NewClient(d.URL, d.User, d.Password)
			if err != nil {
				log.Warnf("connect to download client error: %v", d.URL)
				continue
			}
			return qbt, d, nil
		} else if d.Implementation == downloadclients.ImplementationBuildin {
			bin, err := c.buildInDownloader()
			if err != nil {
				log.Warnf("connect to download client error: %v", err)
				continue
			}
			return bin, d, nil
		}
	}

	return nil, nil, errors.Errorf("no available download client")
}

func (c *Engine) TMDB() (*tmdb.Client, error) {
	api := c.db.GetTmdbApiKey()
	if api == "" {
		return nil, errors.New("TMDB apiKey not set")
	}
	proxy := c.db.GetSetting(db.SettingProxy)
	adult := c.db.GetSetting(db.SettingEnableTmdbAdultContent)
	return tmdb.NewClient(api, proxy, adult == "true")
}

func (c *Engine) MustTMDB() *tmdb.Client {
	t, err := c.TMDB()
	if err != nil {
		log.Panicf("get tmdb: %v", err)
	}
	return t
}

func (c *Engine) RemoveTaskAndTorrent(id int) error {
	torrent, ok := c.tasks.Load(id)
	if ok {
		if err := torrent.Remove(); err != nil {
			return errors.Wrap(err, "remove torrent")
		}
		c.tasks.Delete(id)
	}
	return nil
}

func (c *Engine) GetTasks() utils.Map[int, *Task] {
	return c.tasks
}
