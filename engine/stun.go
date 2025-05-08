package engine

import (
	"fmt"
	"net/url"
	"polaris/ent/downloadclients"
	"polaris/log"
	"polaris/pkg/nat"
	"polaris/pkg/qbittorrent"

	"github.com/pion/stun/v3"
)

func (s *Engine) stunProxyDownloadClient() error {
	
	return s.StartStunProxy("")
}


func (s *Engine) StartStunProxy(name string) error {
	downloaders := s.db.GetAllDonloadClients()
	for _, d := range downloaders {
		if !d.Enable {
			continue
		}
		if !d.UseNatTraversal {
			continue
		}
		if name != "" && d.Name != name {
			continue
		}

		if d.Implementation != downloadclients.ImplementationQbittorrent { //TODO only support qbittorrent for now
			continue
		}

		qbt, err := qbittorrent.NewClient(d.URL, d.User, d.Password)
		if err != nil {
			return fmt.Errorf("connect to download client error: %v", d.URL)
		}
		u, err := url.Parse(d.URL)
		if err != nil {
			return err
		}
		log.Infof("start stun proxy for %s", d.Name)
		n, err := nat.NewNatTraversal(func(xa stun.XORMappedAddress) error {
			return qbt.SetListenPort(xa.Port)
		}, u.Hostname())
		if err != nil {
			return err
		}
		n.StartProxy()

	}
	return nil
}