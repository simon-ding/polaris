package metadata

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type MovieMetadata struct {
	Name       string
	Year       int
	Resolution string
	IsQingban  bool
}

func (m *MovieMetadata) IsAcceptable(name string) bool {
	re := regexp.MustCompile(`[^\p{L}\w\s]`)
	name = re.ReplaceAllString(strings.ToLower(name), " ")
	name2 := re.ReplaceAllString(strings.ToLower(m.Name), " ")
	name = strings.Join(strings.Fields(name), " ")
	name2 = strings.Join(strings.Fields(name2), " ")
	return strings.Contains(name2, name)
}


func ParseMovie(name string) *MovieMetadata {
	name = strings.Join(strings.Fields(name), " ") //remove unnessary spaces
	name = strings.ToLower(strings.TrimSpace(name))
	var meta = &MovieMetadata{}
	yearRe := regexp.MustCompile(`\(\d{4}\)`)
	yearMatches := yearRe.FindAllString(name, -1)
	var yearIndex = -1
	if len(yearMatches) > 0 {
		yearIndex = strings.Index(name, yearMatches[0])
		y := yearMatches[0][1 : len(yearMatches[0])-1]
		n, err := strconv.Atoi(y)
		if err != nil {
			panic(fmt.Sprintf("convert %s error: %v", y, err))
		}
		meta.Year = n
	} else {
		yearRe := regexp.MustCompile(`\d{4}`)
		yearMatches := yearRe.FindAllString(name, -1)
		if len(yearMatches) > 0 {
			n, err := strconv.Atoi(yearMatches[0])
			if err != nil {
				panic(fmt.Sprintf("convert %s error: %v", yearMatches[0], err))
			}
			meta.Year = n
		}
	}

	if yearIndex != -1 {
		meta.Name = name[:yearIndex]
	} else {
		meta.Name = name
	}
	resRe := regexp.MustCompile(`\d{3,4}p`)
	resMatches := resRe.FindAllString(name, -1)
	if len(resMatches) > 0 {
		meta.Resolution = resMatches[0]
	}
	meta.IsQingban = isQiangban(name)
	return meta
}

// https://en.wikipedia.org/wiki/Pirated_movie_release_types
func isQiangban(name string) bool {
	qiangbanFilter := []string{"CAMRip","CAM-Rip", "CAM", "HDCAM", "TS","TSRip", "HDTS", "TELESYNC", "PDVD", "PreDVDRip", "TC", "HDTC", "TELECINE", "WP", "WORKPRINT"}
	re := regexp.MustCompile(`\W`)
	name = re.ReplaceAllString(strings.ToLower(name), " ")
	fields := strings.Fields(name)
	for _, q := range qiangbanFilter {
		for _, f := range fields {
			if strings.EqualFold(q, f) {
				return true
			}
		}
	}
	return false
}
