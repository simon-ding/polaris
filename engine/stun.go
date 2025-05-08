package engine

import (
	"net/url"
	"polaris/ent/downloadclients"
	"polaris/pkg/nat"
	"polaris/pkg/qbittorrent"

	"github.com/pion/stun/v3"
)

func (s *Engine) stunProxyDownloadClient() error {
	downloader, e, err := s.GetDownloadClient()
	if err != nil {
		return err
	}
	if !e.UseNatTraversal {
		return nil
	}
	if e.Implementation != downloadclients.ImplementationQbittorrent {
		return nil
	}
	d, ok := downloader.(*qbittorrent.Client)
	if !ok {
		return nil
	}
	u, err := url.Parse(d.URL)
	if err != nil {
		return err
	}

	n, err := nat.NewNatTraversal(func(xa stun.XORMappedAddress) error {
		return d.SetListenPort(xa.Port)
	}, u.Hostname())
	if err != nil {
		return err
	}
	n.StartProxy()
	return nil
}
