package db

import "polaris/ent/media"

var Version = "undefined"

const (
	SettingTmdbApiKey             = "tmdb_api_key"
	SettingLanguage               = "language"
	SettingJacketUrl              = "jacket_url"
	SettingJacketApiKey           = "jacket_api_key"
	SettingDownloadDir            = "download_dir"
	SettingLogLevel               = "log_level"
	SettingProxy                  = "proxy"
	SettingPlexMatchEnabled       = "plexmatch_enabled"
	SettingNfoSupportEnabled      = "nfo_support_enabled"
	SettingAllowQiangban          = "filter_qiangban"
	SettingEnableTmdbAdultContent = "tmdb_adult_content"
	SettingTvNamingFormat         = "tv_naming_format"
	SettingMovieNamingFormat      = "movie_naming_format"
	SettingProwlarrInfo           = "prowlarr_info"

	SettingTvSizeLimiter          = "tv_size_limiter"
	SettingMovieSizeLimiter       = "movie_size_limiter"
)

const (
	SettingAuthEnabled = "auth_enbled"
	SettingUsername    = "auth_username"
	SettingPassword    = "auth_password"
)

const (
	IndexerTorznabImpl = "torznab"
	DataPath           = "./data"
	ImgPath            = DataPath + "/img"
	LogPath            = DataPath + "/logs"
)

const (
	LanguageEN = "en-US"
	LanguageCN = "zh-CN"
)

const DefaultNamingFormat = "{{.NameCN}} {{.NameEN}} {{if .Year}} ({{.Year}}) {{end}}"

type NamingInfo struct {
	NameCN string
	NameEN string
	Year   string
	TmdbID int
}

type ResolutionType string

const JwtSerectKey = "jwt_secrect_key"

type MediaSizeLimiter struct {
	P720p SizeLimiter `json:"720p"`
	P1080 SizeLimiter `json:"1080p"`
	P2160 SizeLimiter `json:"2160p"`
}

func (m *MediaSizeLimiter) GetLimiter(r media.Resolution) SizeLimiter {
	if r == media.Resolution1080p {
		return m.P1080
	} else if r == media.Resolution720p {
		return m.P720p
	} else if r == media.Resolution2160p {
		return m.P2160
	}
	return SizeLimiter{}
}

type SizeLimiter struct {
	MaxSIze    int64 `json:"max_size"`
	MinSize    int64 `json:"min_size"`
	PreferSIze int64 `json:"prefer_size"`
}

type ProwlarrSetting struct {
	Disabled bool   `json:"disabled"`
	ApiKey   string `json:"api_key"`
	URL      string `json:"url"`
}
