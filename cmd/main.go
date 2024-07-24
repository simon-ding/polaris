package main

import (
	"polaris/log"
	"polaris/pkg/metadata"
	"polaris/pkg/utils"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	b := utils.IsNameAcceptable("legal high 2_勝利即是正", "胜利即是正义")
	log.Info(b)
	m := metadata.ParseMovie("	Inside Out (2013) 1080p WEBRip x264 -YTS")
	log.Infof("%+v", m)
	// dbClient, err := db.Open()
	// if err != nil {
	// 	log.Panicf("init db error: %v", err)
	// }

	// s := server.NewServer(dbClient)
	// if err := s.Serve(); err != nil {
	// 	log.Errorf("server start error: %v", err)
	// }
}

func preProcess(name string) string {
	re := regexp.MustCompile(`[^\p{L}\w\s]`)
	name1 := re.ReplaceAllString(strings.ToLower(name), "")
	return name1
}

func asciiString(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}
