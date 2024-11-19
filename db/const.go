package db

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
	Setting720pSizeLimiter        = "720p_size_limiter"
	Setting1080ppSizeLimiter      = "1080p_size_limiter"
	Setting2160ppSizeLimiter      = "2160p_size_limiter"
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
