package metadata

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type MovieMetadata struct {
	NameEn     string
	NameCN     string
	Year       int
	Resolution string
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
	}
	if yearIndex != -1 {
		meta.NameEn = name[:yearIndex]
	} else {
		meta.NameEn = name
	}
	resRe := regexp.MustCompile(`\d{3,4}p`)
	resMatches := resRe.FindAllString(name, -1)
	if len(resMatches) > 0 {
		meta.Resolution = resMatches[0]
	}
	return meta
}
