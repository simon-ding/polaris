package core

import (
	"polaris/db"
	"polaris/ent"
	"polaris/ent/downloadclients"
	"polaris/log"
	"polaris/pkg"
	"polaris/pkg/qbittorrent"
	"polaris/pkg/tmdb"
	"polaris/pkg/transmission"
	"polaris/pkg/utils"

	"github.com/pkg/errors"
	"github.com/robfig/cron"
)

func NewClient(db *db.Client, language string) *Client {
	return &Client{
		db:       db,
		cron:     cron.New(),
		tasks:    make(map[int]*Task, 0),
		language: language,
	}
}

type scheduler struct {
	cron string
	f    func() error
}
type Client struct {
	db         *db.Client
	cron       *cron.Cron
	tasks      map[int]*Task
	language   string
	schedulers utils.Map[string, scheduler]
}

func (c *Client) registerCronJob(name string, cron string, f func() error) {
	c.schedulers.Store(name, scheduler{
		cron: cron,
		f:    f,
	})
}

func (c *Client) Init() {
	c.reloadTasks()
	c.addSysCron()
}

func (c *Client) reloadTasks() {
	allTasks := c.db.GetRunningHistories()
	for _, t := range allTasks {
		dl, err := c.db.GetDownloadClient(t.DownloadClientID)
		if err != nil {
			log.Warnf("no download client related: %v", t.SourceTitle)
			continue
		}

		if dl.Implementation == downloadclients.ImplementationTransmission {
			to, err := transmission.NewTorrent(transmission.Config{
				URL:      dl.URL,
				User:     dl.User,
				Password: dl.Password,
			}, t.Link)
			if err != nil {
				log.Warnf("get task error: %v", err)
				continue
			}
			c.tasks[t.ID] = &Task{Torrent: to}
		} else if dl.Implementation == downloadclients.ImplementationQbittorrent {
			to, err := qbittorrent.NewTorrent(qbittorrent.Info{
				URL:      dl.URL,
				User:     dl.User,
				Password: dl.Password,
			}, t.Link)
			if err != nil {
				log.Warnf("get task error: %v", err)
				continue
			}
			c.tasks[t.ID] = &Task{Torrent: to}
		}

	}
}

func (c *Client) GetDownloadClient() (pkg.Downloader, *ent.DownloadClients, error) {
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
		}
	}
	return nil, nil, errors.Errorf("no available download client")
}

func (c *Client) TMDB() (*tmdb.Client, error) {
	api := c.db.GetSetting(db.SettingTmdbApiKey)
	if api == "" {
		return nil, errors.New("TMDB apiKey not set")
	}
	proxy := c.db.GetSetting(db.SettingProxy)
	adult := c.db.GetSetting(db.SettingEnableTmdbAdultContent)
	return tmdb.NewClient(api, proxy, adult == "true")
}

func (c *Client) MustTMDB() *tmdb.Client {
	t, err := c.TMDB()
	if err != nil {
		log.Panicf("get tmdb: %v", err)
	}
	return t
}

func (c *Client) RemoveTaskAndTorrent(id int) error {
	torrent := c.tasks[id]
	if torrent != nil {
		if err := torrent.Remove(); err != nil {
			return errors.Wrap(err, "remove torrent")
		}
		delete(c.tasks, id)
	}
	return nil
}

func (c *Client) GetTasks() map[int]*Task {
	return c.tasks
}
