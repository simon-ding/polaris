package metadata

import (
	"fmt"
	"polaris/log"
	"polaris/pkg/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	NameEn       string
	NameCn       string
	Year         int
	Season       int
	StartEpisode int
	EndEpisode   int
	Resolution   string
	IsSeasonPack bool
}

func (m *Info) ParseExtraDescription(desc string) {
	if m.IsSeasonPack { //try to parse episode number with description
		mm := ParseTv(desc)
		if mm.StartEpisode > 0 { //sometimes they put episode info in desc text
			m.IsSeasonPack = false
			m.StartEpisode = mm.StartEpisode
			m.EndEpisode = mm.EndEpisode
		}
	}
}

func (m *Info) IsAcceptable(names ...string) bool {
	re := regexp.MustCompile(`[^\p{L}\w\s]`)

	nameCN := re.ReplaceAllString(strings.ToLower(m.NameCn), " ")
	nameEN := re.ReplaceAllString(strings.ToLower(m.NameEn), " ")
	nameCN = strings.Join(strings.Fields(nameCN), " ")
	nameEN = strings.Join(strings.Fields(nameEN), " ")

	for _, name := range names {
		name = re.ReplaceAllString(strings.ToLower(name), " ")
		name = strings.Join(strings.Fields(name), " ")
		if utils.IsASCII(name) { //ascii name should match words
			re := regexp.MustCompile(`\b` + name + `\b`)
			if re.MatchString(nameCN) || re.MatchString(nameEN) {
				return true
			} else {
				continue
			}
		}

		if strings.Contains(nameCN, name) || strings.Contains(nameEN, name) {
			return true
		}

	}
	return false
}

func ParseTv(name string) *Info {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "\u200b", "") //remove unicode hidden character

	return parseName(name)
}

func adjacentNumber(s string, start int) (n1 int, l int) {
	runes := []rune(s)
	if start > len(runes)-1 { //out of bound
		return -1, -1
	}
	var n []rune
	for i := start; i < len(runes); i++ {
		k := runes[i]
		if (k < '0' || k > '9') && !chineseNum[k] { //not digit anymore
			break
		}
		n = append(n, k)
	}
	if len(n) == 0 {
		return -1, -1
	}
	m, err := strconv.Atoi(string(n))
	if err != nil {
		return chinese2Num[string(n)], len(n)
	}
	return m, len(n)
}

func findSeason(s string) (n int, p int) {
	//season numner
	seasonRe1 := regexp.MustCompile(`s\d{1,2}`)
	seasonMatches := seasonRe1.FindAllString(s, -1)
	if len(seasonMatches) > 0 {
		seNum := seasonMatches[0][1:]
		n, err := strconv.Atoi(seNum)
		if err != nil {
			panic(fmt.Sprintf("convert %s error: %v", seNum, err))
		}

		return n, strings.Index(s, seNum)
	} else {
		seasonRe1 := regexp.MustCompile(`season \d{1,2}`)
		seasonMatches := seasonRe1.FindAllString(s, -1)
		if len(seasonMatches) > 0 {
			re3 := regexp.MustCompile(`\d{1,2}`)
			seNum := re3.FindAllString(seasonMatches[0], -1)[0]
			n, err := strconv.Atoi(seNum)
			if err != nil {
				panic(fmt.Sprintf("convert %s error: %v", seNum, err))
			}
			return n, strings.Index(s, seasonMatches[0])
		} else {
			seasonRe1 := regexp.MustCompile(`第.{1,2}季`)
			seasonMatches := seasonRe1.FindAllString(s, -1)
			if len(seasonMatches) > 0 {
				m1 := []rune(seasonMatches[0])
				seNum := m1[1 : len(m1)-1]
				n, err := strconv.Atoi(string(seNum))
				if err != nil {
					log.Warnf("parse season number %v error: %v, try to parse using chinese", seNum, err)
					n = chinese2Num[string(seNum)]
				}
				return n, strings.Index(s, seasonMatches[0])
			}
		}
	}
	return -1, -1
}

func findEpisodes(s string) (start int, end int) {
	var episodeCn = map[rune]bool{
		'话': true,
		'話': true,
		'集': true,
	}

	rr := []rune(s)
	for i := 0; i < len(rr); i++ {
		r := rr[i]
		if r == 'e' {
			n, l := adjacentNumber(s, i+1)

			if n > 0 {
				foundDash := false
				for j := i + l + 1; j < len(rr); j++ {
					r1 := rr[j]
					if r1 == '-' {
						foundDash = true
						continue
					}
					if r1 == ' ' || r1 == 'e' {
						continue
					}

					if foundDash {
						if r1 == 's' {
							s1, l1 := adjacentNumber(s, j+1)
							if s1 > 0 { //S01E01-S01E21
								n1, _ := adjacentNumber(s, j+l1+2)
								if n1 > 0 {
									return n, n1
								}
							}
						}
						n1, _ := adjacentNumber(s, j)
						if n1 > 0 {
							return n, n1
						}
					} else {
						break
					}
				}
				return n, n
			}
		} else if r == '第' {
			n, l := adjacentNumber(s, i+1)
			if len(rr) > i+l+1 && episodeCn[rr[i+l+1]] {
				return n, n
			} else if len(rr) > i+l+1 {
				if rr[i+l+1] == '-' {
					n1, l1 := adjacentNumber(s, i+l+2)
					if episodeCn[rr[i+l+2+l1]] {
						return n, n1
					}
				}
			}

		}
	}
	//episode number
	re1 := regexp.MustCompile(`\[\d{1,4}\]`)
	episodeMatches1 := re1.FindAllString(s, -1)
	if len(episodeMatches1) > 0 { //[11] [1080p], [2022][113][HEVC][GB][4K]
		for _, m := range episodeMatches1 {
			epNum := strings.TrimRight(strings.TrimLeft(m, "["), "]")
			n, err := strconv.Atoi(epNum)
			if err != nil {
				log.Debugf("convert %s error: %v", epNum, err)
				continue
			}
			nowYear := time.Now().Year()
			if n > nowYear-50 { //high possibility is year number
				continue
			}
			return n, n
		}
	} else { //【第09話】
		re2 := regexp.MustCompile(`第\d{1,4}([话話集])`)
		episodeMatches1 := re2.FindAllString(s, -1)
		if len(episodeMatches1) > 0 {
			re := regexp.MustCompile(`\d{1,4}`)
			epNum := re.FindAllString(episodeMatches1[0], -1)[0]
			n, err := strconv.Atoi(epNum)
			if err != nil {
				panic(fmt.Sprintf("convert %s error: %v", epNum, err))
			}
			return n, n
		} else { //The Road Season 2 Episode 12 XviD-AFG
			re3 := regexp.MustCompile(`episode \d{1,4}`)
			epNums := re3.FindAllString(s, -1)
			if len(epNums) > 0 {
				re3 := regexp.MustCompile(`\d{1,4}`)
				epNum := re3.FindAllString(epNums[0], -1)[0]
				n, err := strconv.Atoi(epNum)
				if err != nil {
					panic(fmt.Sprintf("convert %s error: %v", epNum, err))
				}
				return n, n

			} else { //SHY 靦腆英雄 / Shy -05 ( CR 1920x1080 AVC AAC MKV)
				if maybeSeasonPack(s) { //avoid miss match, season pack not use this rule
					return -1, -1
				}
				re3 := regexp.MustCompile(`[^(season)][^\d\w]\d{1,2}[^\d\w]`)
				epNums := re3.FindAllString(s, -1)
				if len(epNums) > 0 {

					re3 := regexp.MustCompile(`\d{1,2}`)
					epNum := re3.FindAllString(epNums[0], -1)[0]
					n, err := strconv.Atoi(epNum)
					if err != nil {
						panic(fmt.Sprintf("convert %s error: %v", epNum, err))
					}
					return n, n
				}
			}
		}
	}

	return -1, -1
}

func matchResolution(s string) string {
	//resolution
	resRe := regexp.MustCompile(`\d{3,4}p`)
	resMatches := resRe.FindAllString(s, -1)
	if len(resMatches) != 0 {
		return resMatches[0]
	} else {
		if strings.Contains(s, "720") {
			return "720p"
		} else if strings.Contains(s, "1080") {
			return "1080p"
		}
	}
	return ""
}

func maybeSeasonPack(s string) bool {
	//season pack
	packRe := regexp.MustCompile(`((\d{1,2}-\d{1,2}))|(complete)|(全集)`)
	if packRe.MatchString(s) {
		return true
	}
	return false
}

//func parseEnglishName(name string) *Info {
//	meta := &Info{
//		//Season:  -1,
//		Episode: -1,
//	}
//
//	start, end := findEpisodes(name)
//	if start > 0 && end > 0 {
//		meta.Episode = start
//	}
//
//	re := regexp.MustCompile(`[^\p{L}\w\s]`)
//	name = re.ReplaceAllString(strings.ToLower(name), " ")
//	newSplits := strings.Split(strings.TrimSpace(name), " ")
//
//	seasonRe := regexp.MustCompile(`^s\d{1,2}`)
//	resRe := regexp.MustCompile(`^\d{3,4}p`)
//	episodeRe := regexp.MustCompile(`e\d{1,3}`)
//
//	var seasonIndex = -1
//	var episodeIndex = -1
//	var resIndex = -1
//	for i, p := range newSplits {
//		p = strings.TrimSpace(p)
//		if p == "" {
//			continue
//		}
//		if seasonRe.MatchString(p) {
//			//season part
//			seasonIndex = i
//		} else if resRe.MatchString(p) {
//			resIndex = i
//		}
//		if i >= seasonIndex && episodeRe.MatchString(p) {
//			episodeIndex = i
//		}
//	}
//
//	if seasonIndex != -1 {
//		//season exists
//		ss := seasonRe.FindAllString(newSplits[seasonIndex], -1)
//		if len(ss) != 0 {
//			//season info
//
//			ssNum := strings.TrimLeft(ss[0], "s")
//			n, err := strconv.Atoi(ssNum)
//			if err != nil {
//				panic(fmt.Sprintf("convert %s error: %v", ssNum, err))
//			}
//			meta.Season = n
//		}
//	} else { //maybe like Season 1?
//		seasonRe := regexp.MustCompile(`season \d{1,2}`)
//		matches := seasonRe.FindAllString(name, -1)
//		if len(matches) > 0 {
//			for i, s := range newSplits {
//				if s == "season" {
//					seasonIndex = i
//				}
//			}
//			numRe := regexp.MustCompile(`\d{1,2}`)
//			seNum := numRe.FindAllString(matches[0], -1)[0]
//			n, err := strconv.Atoi(seNum)
//			if err != nil {
//				panic(fmt.Sprintf("convert %s error: %v", seNum, err))
//			}
//			meta.Season = n
//
//		}
//	}
//
//	if episodeIndex != -1 {
//		//	ep := episodeRe.FindAllString(newSplits[episodeIndex], -1)
//		//if len(ep) > 0 {
//		//	//episode info exists
//		//	epNum := strings.TrimLeft(ep[0], "e")
//		//	n, err := strconv.Atoi(epNum)
//		//	if err != nil {
//		//		panic(fmt.Sprintf("convert %s error: %v", epNum, err))
//		//	}
//		//	meta.Episode = n
//		//}
//	} else { //no episode, maybe like  One Punch Man S2 - 08 [1080p].mkv
//
//		// numRe := regexp.MustCompile(`^\d{1,2}$`)
//		// for i, p := range newSplits {
//		// 	if numRe.MatchString(p) {
//		// 		if i > 0 && strings.Contains(newSplits[i-1], "season") { //last word cannot be season
//		// 			continue
//		// 		}
//		// 		if i < seasonIndex {
//		// 			//episode number most likely  should comes alfter season number
//		// 			continue
//		// 		}
//		// 		//episodeIndex = i
//		// 		n, err := strconv.Atoi(p)
//		// 		if err != nil {
//		// 			panic(fmt.Sprintf("convert %s error: %v", p, err))
//		// 		}
//		// 		meta.Episode = n
//
//		// 	}
//		// }
//
//	}
//	if resIndex != -1 {
//		//resolution exists
//		meta.Resolution = newSplits[resIndex]
//	}
//	if meta.Episode == -1 {
//		meta.Episode = -1
//		meta.IsSeasonPack = true
//	}
//
//	if seasonIndex > 0 {
//		//name exists
//		names := newSplits[0:seasonIndex]
//		meta.NameEn = strings.TrimSpace(strings.Join(names, " "))
//	} else {
//		meta.NameEn = name
//	}
//
//	return meta
//}

func parseName(name string) *Info {
	meta := &Info{Season: 1}
	if strings.TrimSpace(name) == "" {
		return meta
	}

	season, p := findSeason(name)
	if season == -1 {
		log.Debugf("not find season info: %s", name)
		if !utils.IsASCII(name) {
			season = 1
		}
		p = len(name) - 1
	}
	meta.Season = season

	start, end := findEpisodes(name)
	if start > 0 && end > 0 {
		meta.StartEpisode = start
		meta.EndEpisode = end
	} else {
		meta.IsSeasonPack = true
	}

	meta.Resolution = matchResolution(name)

	//if meta.IsSeasonPack && meta.Episode != 0 {
	//	meta.Season = meta.Episode
	//	meta.Episode = -1
	//}

	//tv name
	if utils.IsASCII(name) && p < len(name) && p-1 > 0 {
		meta.NameEn = strings.TrimSpace(name[:p-1])
		meta.NameCn = meta.NameEn
	} else {
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
		meta.NameEn = title

		////匹配title中最长拉丁字符串
		//enRe := regexp.MustCompile(`[[:ascii:]]*`)
		//enM := enRe.FindAllString(title, -1)
		//if len(enM) > 0 {
		//	for _, t := range enM {
		//		if len(t) > len(meta.NameEn) {
		//			meta.NameEn = strings.TrimSpace(strings.ToLower(t))
		//		}
		//	}
		//}

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

var chineseNum = map[rune]bool{
	'一': true,
	'二': true,
	'三': true,
	'四': true,
	'五': true,
	'六': true,
	'七': true,
	'八': true,
	'九': true,
}
