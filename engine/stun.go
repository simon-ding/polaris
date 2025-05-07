package engine

import (
	"net/url"
	"polaris/ent/downloadclients"
	"polaris/pkg/nat"
	"polaris/pkg/qbittorrent"
	"strconv"
)

func (s *Engine) stunProxyDownloadClient() error {
	downloader, e, err := s.GetDownloadClient()
	if err != nil {
		return err
	}
	if e.Implementation != downloadclients.ImplementationQbittorrent {
		return nil
	}
	d, ok := downloader.(*qbittorrent.Client)
	if !ok {
		return nil
	}
	n, err := nat.NewNatTraversal()
	if err != nil {
		return err
	}
	addr, err := n.StunAddr()
	if err != nil {
		return err
	}
	err = d.SetListenPort(addr.Port)
	if err != nil {
		return err
	}
	u, err := url.Parse(d.URL)
	if err != nil {
		return err
	}

	return n.StartProxy(u.Hostname() + ":" + strconv.Itoa(addr.Port))
}
