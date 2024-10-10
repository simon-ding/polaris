package metadata

import (
	"fmt"
	"polaris/log"
	"polaris/pkg/utils"
	"regexp"
	"strconv"
	"strings"
)

type Metadata struct {
	NameEn       string
	NameCn       string
	Season       int
	Episode      int
	Resolution   string
	IsSeasonPack bool
}

func (m *Metadata) IsAcceptable(names... string) bool {
	for _, name := range names {
		re := regexp.MustCompile(`[^\p{L}\w\s]`)
		name = re.ReplaceAllString(strings.ToLower(name), " ")
		nameCN := re.ReplaceAllString(strings.ToLower(m.NameCn), " ")
		nameEN := re.ReplaceAllString(strings.ToLower(m.NameEn), " ")
		name = strings.Join(strings.Fields(name), " ")
		nameCN = strings.Join(strings.Fields(nameCN), " ")
		nameEN = strings.Join(strings.Fields(nameEN), " ")
		if utils.IsASCII(name) { //ascii name should match words
			re := regexp.MustCompile(`\b` + name + `\b`)
			return re.MatchString(nameCN) ||  re.MatchString(nameEN) 
		}

		if  strings.Contains(nameCN, name) || strings.Contains(nameEN, name) {
			return true
		}
	
	}
	return false
}


func ParseTv(name string) *Metadata {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "\u200b", "") //remove unicode hidden character
	if utils.ContainsChineseChar(name) {
		return parseChineseName(name)
	}
	return parseEnglishName(name)
}

func parseEnglishName(name string) *Metadata {
	re := regexp.MustCompile(`[^\p{L}\w\s]`)
	name = re.ReplaceAllString(strings.ToLower(name), " ")
	newSplits := strings.Split(strings.TrimSpace(name), " ")

	seasonRe := regexp.MustCompile(`^s\d{1,2}`)
	resRe := regexp.MustCompile(`^\d{3,4}p`)
	episodeRe := regexp.MustCompile(`e\d{1,3}`)

	var seasonIndex = -1
	var episodeIndex = -1
	var resIndex = -1
	for i, p := range newSplits {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if seasonRe.MatchString(p) {
			//season part
			seasonIndex = i
		} else if resRe.MatchString(p) {
			resIndex = i
		}
		if i >= seasonIndex && episodeRe.MatchString(p) {
			episodeIndex = i
		}
	}

	meta := &Metadata{
		Season:  -1,
		Episode: -1,
	}
	if seasonIndex != -1 {
		//season exists
		ss := seasonRe.FindAllString(newSplits[seasonIndex], -1)
		if len(ss) != 0 {
			//season info

			ssNum := strings.TrimLeft(ss[0], "s")
			n, err := strconv.Atoi(ssNum)
			if err != nil {
				panic(fmt.Sprintf("convert %s error: %v", ssNum, err))
			}
			meta.Season = n
		}
	} else { //maybe like Season 1?
		seasonRe := regexp.MustCompile(`season \d{1,2}`)
		matches := seasonRe.FindAllString(name, -1)
		if len(matches) > 0 {
			for i, s := range newSplits {
				if s == "season" {
					seasonIndex = i
				}
			}
			numRe := regexp.MustCompile(`\d{1,2}`)
			seNum := numRe.FindAllString(matches[0], -1)[0]
			n, err := strconv.Atoi(seNum)
			if err != nil {
				panic(fmt.Sprintf("convert %s error: %v", seNum, err))
			}
			meta.Season = n

		}
	}
	if episodeIndex != -1 {
		ep := episodeRe.FindAllString(newSplits[episodeIndex], -1)
		if len(ep) > 0 {
			//episode info exists
			epNum := strings.TrimLeft(ep[0], "e")
			n, err := strconv.Atoi(epNum)
			if err != nil {
				panic(fmt.Sprintf("convert %s error: %v", epNum, err))
			}
			meta.Episode = n
		}
	} else { //no episode, maybe like  One Punch Man S2 - 08 [1080p].mkv

		// numRe := regexp.MustCompile(`^\d{1,2}$`)
		// for i, p := range newSplits {
		// 	if numRe.MatchString(p) {
		// 		if i > 0 && strings.Contains(newSplits[i-1], "season") { //last word cannot be season
		// 			continue
		// 		}
		// 		if i < seasonIndex {
		// 			//episode number most likely  should comes alfter season number
		// 			continue
		// 		}
		// 		//episodeIndex = i
		// 		n, err := strconv.Atoi(p)
		// 		if err != nil {
		// 			panic(fmt.Sprintf("convert %s error: %v", p, err))
		// 		}
		// 		meta.Episode = n

		// 	}
		// }

	}
	if resIndex != -1 {
		//resolution exists
		meta.Resolution = newSplits[resIndex]
	}
	if meta.Episode == -1 {
		meta.Episode = -1
		meta.IsSeasonPack = true
	}

	if seasonIndex > 0 {
		//name exists
		names := newSplits[0:seasonIndex]
		meta.NameEn = strings.TrimSpace(strings.Join(names, " "))
	} else {
		meta.NameEn = name
	}

	return meta
}

func parseChineseName(name string) *Metadata {
	var meta = parseEnglishName(name)
	if meta.Season != -1 && (meta.Episode != -1 || meta.IsSeasonPack) {
		return meta
	}
	meta = &Metadata{Season: 1}
	//season pack
	packRe := regexp.MustCompile(`(\d{1,2}-\d{1,2})|(全集)`)
	if packRe.MatchString(name) {
		meta.IsSeasonPack = true
	}
	//resolution
	resRe := regexp.MustCompile(`\d{3,4}p`)
	resMatches := resRe.FindAllString(name, -1)
	if len(resMatches) != 0 {
		meta.Resolution = resMatches[0]
	} else {
		if strings.Contains(name, "720") {
			meta.Resolution = "720p"
		} else if strings.Contains(name, "1080") {
			meta.Resolution = "1080p"
		}

	}

	//episode number
	re1 := regexp.MustCompile(`\[\d{1,3}\]`)
	episodeMatches1 := re1.FindAllString(name, -1)
	if len(episodeMatches1) > 0 { //[11] [1080p]
		epNum := strings.TrimRight(strings.TrimLeft(episodeMatches1[0], "["), "]")
		n, err := strconv.Atoi(epNum)
		if err != nil {
			panic(fmt.Sprintf("convert %s error: %v", epNum, err))
		}
		meta.Episode = n
	} else { //【第09話】
		re2 := regexp.MustCompile(`第\d{1,4}(话|話|集)`)
		episodeMatches1 := re2.FindAllString(name, -1)
		if len(episodeMatches1) > 0 {
			re := regexp.MustCompile(`\d{1,4}`)
			epNum := re.FindAllString(episodeMatches1[0], -1)[0]
			n, err := strconv.Atoi(epNum)
			if err != nil {
				panic(fmt.Sprintf("convert %s error: %v", epNum, err))
			}
			meta.Episode = n
		} else { //SHY 靦腆英雄 / Shy -05 ( CR 1920x1080 AVC AAC MKV)
			re3 := regexp.MustCompile(`[^\d\w]\d{1,2}[^\d\w]`)
			epNums := re3.FindAllString(name, -1)
			if len(epNums) > 0 {

				re3 := regexp.MustCompile(`\d{1,2}`)
				epNum := re3.FindAllString(epNums[0], -1)[0]
				n, err := strconv.Atoi(epNum)
				if err != nil {
					panic(fmt.Sprintf("convert %s error: %v", epNum, err))
				}
				meta.Episode = n
			}
		}
	}

	//season numner
	seasonRe1 := regexp.MustCompile(`s\d{1,2}`)
	seasonMatches := seasonRe1.FindAllString(name, -1)
	if len(seasonMatches) > 0 {
		seNum := seasonMatches[0][1:]
		n, err := strconv.Atoi(seNum)
		if err != nil {
			panic(fmt.Sprintf("convert %s error: %v", seNum, err))
		}
		meta.Season = n
	} else {
		seasonRe1 := regexp.MustCompile(`season \d{1,2}`)
		seasonMatches := seasonRe1.FindAllString(name, -1)
		if len(seasonMatches) > 0 {
			re3 := regexp.MustCompile(`\d{1,2}`)
			seNum := re3.FindAllString(seasonMatches[0], -1)[0]
			n, err := strconv.Atoi(seNum)
			if err != nil {
				panic(fmt.Sprintf("convert %s error: %v", seNum, err))
			}
			meta.Season = n
		} else {
			seasonRe1 := regexp.MustCompile(`第.{1,2}季`)
			seasonMatches := seasonRe1.FindAllString(name, -1)
			if len(seasonMatches) > 0 {
				m1 := []rune(seasonMatches[0])
				seNum := m1[1 : len(m1)-1]
				n, err := strconv.Atoi(string(seNum))
				if err != nil {
					log.Warnf("parse season number %v error: %v, try to parse using chinese", seNum, err)
					n = chinese2Num[string(seNum)]
				}
				meta.Season = n

			}
		}
	}

	if meta.IsSeasonPack && meta.Episode != 0 {
		meta.Season = meta.Episode
		meta.Episode = -1
	}

	//tv name

	fields := strings.FieldsFunc(name, func(r rune) bool {
		return r == '[' || r == ']' || r == '【' || r == '】'
	})
	titleCn := ""
	title := ""
	for _, p := range fields { //寻找匹配的最长的字符串，最有可能是名字
		if utils.ContainsChineseChar(p) && len([]rune(p)) > len([]rune(titleCn)) { //最长含中文字符串
			titleCn = p
		}
		if len([]rune(p)) > len([]rune(title)) { //最长字符串
			title = p
		}
	}
	re := regexp.MustCompile(`[^\p{L}\w\s]`)
	title = re.ReplaceAllString(strings.TrimSpace(strings.ToLower(title)), "") //去除标点符号
	titleCn = re.ReplaceAllString(strings.TrimSpace(strings.ToLower(titleCn)), "")

	meta.NameCn = titleCn
	cnRe := regexp.MustCompile(`\p{Han}.*\p{Han}`)
	cnmatches := cnRe.FindAllString(titleCn, -1)

	//titleCn中最长的中文字符
	if len(cnmatches) > 0 {
		for _, t := range cnmatches {
			if len([]rune(t)) > len([]rune(meta.NameCn)) {
				meta.NameCn = strings.ToLower(t)
			}
		}
	}

	//匹配title中最长拉丁字符串
	enRe := regexp.MustCompile(`[[:ascii:]]*`)
	enM := enRe.FindAllString(title, -1)
	if len(enM) > 0 {
		for _, t := range enM {
			if len(t) > len(meta.NameEn) {
				meta.NameEn = strings.TrimSpace(strings.ToLower(t))
			}
		}
	}

	return meta
}

var chinese2Num = map[string]int{
	"一": 1,
	"二": 2,
	"三": 3,
	"四": 4,
	"五": 5,
	"六": 6,
	"七": 7,
	"八": 8,
	"九": 9,
}
