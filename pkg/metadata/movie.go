package metadata

import (
	"fmt"
	"polaris/pkg/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MovieMetadata struct {
	Name       string
	Year       int
	Resolution string
	IsQingban  bool
}

func (m *MovieMetadata) IsAcceptable(names ...string) bool {
	for _, name := range names {
		re := regexp.MustCompile(`[^\p{L}\w\s]`)
		name = re.ReplaceAllString(strings.ToLower(name), " ")
		name2 := re.ReplaceAllString(strings.ToLower(m.Name), " ")
		name = strings.Join(strings.Fields(name), " ")
		name2 = strings.Join(strings.Fields(name2), " ")
		if utils.IsASCII(name) { //ascii name should match words
			re := regexp.MustCompile(`\b` + name + `\b`)
			return re.MatchString(name2)
		}

		if strings.Contains(name2, name) {
			return true
		}
	}
	return false
}

func findYear(name string) (year int, index int) {
	yearRe := regexp.MustCompile(`\(\d{4}\)`)
	yearMatches := yearRe.FindAllString(name, -1)
	index = -1
	if len(yearMatches) > 0 {
		index = strings.Index(name, yearMatches[0])
		y := yearMatches[0][1 : len(yearMatches[0])-1]
		n, err := strconv.Atoi(y)
		if err != nil {
			panic(fmt.Sprintf("convert %s error: %v", y, err))
		}
		year = n
	} else {
		yearRe := regexp.MustCompile(`\d{4}`)
		yearMatches := yearRe.FindAllString(name, -1)
		if len(yearMatches) > 0 {
			year, index = findYearInMatches(yearMatches, name)

		}
	}
	return
}

func findYearInMatches(matches []string, name string) (year int, index int) {
	if len(matches) == 0 {
		return 0, -1
	}
	for _, y := range matches {
		index = strings.Index(name, y)
		n, err := strconv.Atoi(y)
		if err != nil {
			panic(fmt.Sprintf("convert %s error: %v", y, err))
		}
		if n < 1900 || n > time.Now().Year()+1 { //filter invalid year
			continue
		}
		year = n

	}
	return
}

func ParseMovie(name string) *MovieMetadata {
	name = strings.Join(strings.Fields(name), " ") //remove unnessary spaces
	name = strings.ToLower(strings.TrimSpace(name))
	var meta = &MovieMetadata{}
	year, yearIndex := findYear(name)

	meta.Year = year

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
	qiangbanFilter := []string{"CAMRip", "CAM-Rip", "CAM", "HDCAM", "TS", "TSRip", "HDTS", "TELESYNC", "PDVD", "PreDVDRip", "TC", "HDTC", "TELECINE", "WP", "WORKPRINT"}
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
