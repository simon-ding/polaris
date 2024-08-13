package core

import "encoding/xml"

type Tvshow struct {
	XMLName       xml.Name `xml:"tvshow"`
	Text          string   `xml:",chardata"`
	Title         string   `xml:"title"`
	Originaltitle string   `xml:"originaltitle"`
	Showtitle     string   `xml:"showtitle"`
	Ratings       struct {
		Text   string `xml:",chardata"`
		Rating []struct {
			Text    string `xml:",chardata"`
			Name    string `xml:"name,attr"`
			Max     string `xml:"max,attr"`
			Default string `xml:"default,attr"`
			Value   string `xml:"value"`
			Votes   string `xml:"votes"`
		} `xml:"rating"`
	} `xml:"ratings"`
	Userrating     string `xml:"userrating"`
	Top250         string `xml:"top250"`
	Season         string `xml:"season"`
	Episode        string `xml:"episode"`
	Displayseason  string `xml:"displayseason"`
	Displayepisode string `xml:"displayepisode"`
	Outline        string `xml:"outline"`
	Plot           string `xml:"plot"`
	Tagline        string `xml:"tagline"`
	Runtime        string `xml:"runtime"`
	Thumb          []struct {
		Text    string `xml:",chardata"`
		Spoof   string `xml:"spoof,attr"`
		Cache   string `xml:"cache,attr"`
		Aspect  string `xml:"aspect,attr"`
		Preview string `xml:"preview,attr"`
		Season  string `xml:"season,attr"`
		Type    string `xml:"type,attr"`
	} `xml:"thumb"`
	Fanart struct {
		Text  string `xml:",chardata"`
		Thumb []struct {
			Text    string `xml:",chardata"`
			Colors  string `xml:"colors,attr"`
			Preview string `xml:"preview,attr"`
		} `xml:"thumb"`
	} `xml:"fanart"`
	Mpaa       string     `xml:"mpaa"`
	Playcount  string     `xml:"playcount"`
	Lastplayed string     `xml:"lastplayed"`
	ID         string     `xml:"id"`
	Uniqueid   []UniqueId `xml:"uniqueid"`
	Genre      string     `xml:"genre"`
	Premiered  string     `xml:"premiered"`
	Year       string     `xml:"year"`
	Status     string     `xml:"status"`
	Code       string     `xml:"code"`
	Aired      string     `xml:"aired"`
	Studio     string     `xml:"studio"`
	Trailer    string     `xml:"trailer"`
	Actor      []struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name"`
		Role  string `xml:"role"`
		Order string `xml:"order"`
		Thumb string `xml:"thumb"`
	} `xml:"actor"`
	Namedseason []struct {
		Text   string `xml:",chardata"`
		Number string `xml:"number,attr"`
	} `xml:"namedseason"`
	Resume struct {
		Text     string `xml:",chardata"`
		Position string `xml:"position"`
		Total    string `xml:"total"`
	} `xml:"resume"`
	Dateadded string `xml:"dateadded"`
}

type UniqueId struct {
	Text    string `xml:",chardata"`
	Type    string `xml:"type,attr"`
	Default string `xml:"default,attr"`
}

type Episodedetails struct {
	XMLName   xml.Name `xml:"episodedetails"`
	Text      string   `xml:",chardata"`
	Title     string   `xml:"title"`
	Showtitle string   `xml:"showtitle"`
	Ratings   struct {
		Text   string `xml:",chardata"`
		Rating []struct {
			Text    string `xml:",chardata"`
			Name    string `xml:"name,attr"`
			Max     string `xml:"max,attr"`
			Default string `xml:"default,attr"`
			Value   string `xml:"value"`
			Votes   string `xml:"votes"`
		} `xml:"rating"`
	} `xml:"ratings"`
	Userrating     string `xml:"userrating"`
	Top250         string `xml:"top250"`
	Season         string `xml:"season"`
	Episode        string `xml:"episode"`
	Displayseason  string `xml:"displayseason"`
	Displayepisode string `xml:"displayepisode"`
	Outline        string `xml:"outline"`
	Plot           string `xml:"plot"`
	Tagline        string `xml:"tagline"`
	Runtime        string `xml:"runtime"`
	Thumb          []struct {
		Text    string `xml:",chardata"`
		Spoof   string `xml:"spoof,attr"`
		Cache   string `xml:"cache,attr"`
		Aspect  string `xml:"aspect,attr"`
		Preview string `xml:"preview,attr"`
	} `xml:"thumb"`
	Mpaa       string `xml:"mpaa"`
	Playcount  string `xml:"playcount"`
	Lastplayed string `xml:"lastplayed"`
	ID         string `xml:"id"`
	Uniqueid   []struct {
		Text    string `xml:",chardata"`
		Type    string `xml:"type,attr"`
		Default string `xml:"default,attr"`
	} `xml:"uniqueid"`
	Genre     string   `xml:"genre"`
	Credits   []string `xml:"credits"`
	Director  string   `xml:"director"`
	Premiered string   `xml:"premiered"`
	Year      string   `xml:"year"`
	Status    string   `xml:"status"`
	Code      string   `xml:"code"`
	Aired     string   `xml:"aired"`
	Studio    string   `xml:"studio"`
	Trailer   string   `xml:"trailer"`
	Actor     []struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name"`
		Role  string `xml:"role"`
		Order string `xml:"order"`
		Thumb string `xml:"thumb"`
	} `xml:"actor"`
	Resume struct {
		Text     string `xml:",chardata"`
		Position string `xml:"position"`
		Total    string `xml:"total"`
	} `xml:"resume"`
	Dateadded string `xml:"dateadded"`
}

type Movie struct {
	XMLName       xml.Name `xml:"movie"`
	Text          string   `xml:",chardata"`
	Title         string   `xml:"title"`
	Originaltitle string   `xml:"originaltitle"`
	Sorttitle     string   `xml:"sorttitle"`
	Ratings       struct {
		Text   string `xml:",chardata"`
		Rating []struct {
			Text    string `xml:",chardata"`
			Name    string `xml:"name,attr"`
			Max     string `xml:"max,attr"`
			Default string `xml:"default,attr"`
			Value   string `xml:"value"`
			Votes   string `xml:"votes"`
		} `xml:"rating"`
	} `xml:"ratings"`
	Userrating string `xml:"userrating"`
	Top250     string `xml:"top250"`
	Outline    string `xml:"outline"`
	Plot       string `xml:"plot"`
	Tagline    string `xml:"tagline"`
	Runtime    string `xml:"runtime"`
	Thumb      []struct {
		Text    string `xml:",chardata"`
		Spoof   string `xml:"spoof,attr"`
		Cache   string `xml:"cache,attr"`
		Aspect  string `xml:"aspect,attr"`
		Preview string `xml:"preview,attr"`
	} `xml:"thumb"`
	Fanart struct {
		Text  string `xml:",chardata"`
		Thumb struct {
			Text    string `xml:",chardata"`
			Colors  string `xml:"colors,attr"`
			Preview string `xml:"preview,attr"`
		} `xml:"thumb"`
	} `xml:"fanart"`
	Mpaa       string     `xml:"mpaa"`
	Playcount  string     `xml:"playcount"`
	Lastplayed string     `xml:"lastplayed"`
	ID         string     `xml:"id"`
	Uniqueid   []UniqueId `xml:"uniqueid"`
	Genre      string     `xml:"genre"`
	Country    []string   `xml:"country"`
	Set        struct {
		Text     string `xml:",chardata"`
		Name     string `xml:"name"`
		Overview string `xml:"overview"`
	} `xml:"set"`
	Tag                   []string `xml:"tag"`
	Videoassettitle       string   `xml:"videoassettitle"`
	Videoassetid          string   `xml:"videoassetid"`
	Videoassettype        string   `xml:"videoassettype"`
	Hasvideoversions      string   `xml:"hasvideoversions"`
	Hasvideoextras        string   `xml:"hasvideoextras"`
	Isdefaultvideoversion string   `xml:"isdefaultvideoversion"`
	Credits               []string `xml:"credits"`
	Director              string   `xml:"director"`
	Premiered             string   `xml:"premiered"`
	Year                  string   `xml:"year"`
	Status                string   `xml:"status"`
	Code                  string   `xml:"code"`
	Aired                 string   `xml:"aired"`
	Studio                string   `xml:"studio"`
	Trailer               string   `xml:"trailer"`
	Fileinfo              struct {
		Text          string `xml:",chardata"`
		Streamdetails struct {
			Text  string `xml:",chardata"`
			Video struct {
				Text              string `xml:",chardata"`
				Codec             string `xml:"codec"`
				Aspect            string `xml:"aspect"`
				Width             string `xml:"width"`
				Height            string `xml:"height"`
				Durationinseconds string `xml:"durationinseconds"`
				Stereomode        string `xml:"stereomode"`
				Hdrtype           string `xml:"hdrtype"`
			} `xml:"video"`
			Audio struct {
				Text     string `xml:",chardata"`
				Codec    string `xml:"codec"`
				Language string `xml:"language"`
				Channels string `xml:"channels"`
			} `xml:"audio"`
			Subtitle struct {
				Text     string `xml:",chardata"`
				Language string `xml:"language"`
			} `xml:"subtitle"`
		} `xml:"streamdetails"`
	} `xml:"fileinfo"`
	Actor []struct {
		Text  string `xml:",chardata"`
		Name  string `xml:"name"`
		Role  string `xml:"role"`
		Order string `xml:"order"`
		Thumb string `xml:"thumb"`
	} `xml:"actor"`
}
