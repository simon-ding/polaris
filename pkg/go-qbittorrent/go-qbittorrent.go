package qbittorrent

import (
	"fmt"
	"polaris/pkg/go-qbittorrent/qbt"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	// connect to qbittorrent client
	qb := qbt.NewClient("http://localhost:8181")

	// login to the client
	loginOpts := qbt.LoginOptions{
		Username: "username",
		Password: "password",
	}
	err := qb.Login(loginOpts)
	if err != nil {
		fmt.Println(err)
	}

	// ********************
	// DOWNLOAD A TORRENT *
	// ********************

	// were not using any filters so the options map is empty
	downloadOpts := qbt.DownloadOptions{}
	// set the path to the file
	//path := "/Users/me/Downloads/Source.Code.2011.1080p.BluRay.H264.AAC-RARBG-[rarbg.to].torrent"
	links := []string{"http://rarbg.to/download.php?id=9buc5hp&h=d73&f=Courage.the.Cowardly.Dog.1999.S01.1080p.AMZN.WEBRip.DD2.0.x264-NOGRP%5Brartv%5D-[rarbg.to].torrent"}
	// download the torrent using the file
	// the wrapper will handle opening and closing the file for you
	err = qb.DownloadLinks(links, downloadOpts)

	if err != nil {
		fmt.Println("[-] Download torrent from link")
		fmt.Println(err)
	} else {
		fmt.Println("[+] Download torrent from link")
	}

	// ******************
	// GET ALL TORRENTS *
	// ******************
	torrentsOpts := qbt.TorrentsOptions{}
	filter := "inactive"
	sort := "name"
	hash := "d739f78a12b241ba62719b1064701ffbb45498a8"
	torrentsOpts.Filter = &filter
	torrentsOpts.Sort = &sort
	torrentsOpts.Hashes = []string{hash}
	torrents, err := qb.Torrents(torrentsOpts)
	if err != nil {
		fmt.Println("[-] Get torrent list")
		fmt.Println(err)
	} else {
		fmt.Println("[+] Get torrent list")
		if len(torrents) > 0 {
			spew.Dump(torrents[0])
		} else {
			fmt.Println("No torrents found")
		}
	}
}
