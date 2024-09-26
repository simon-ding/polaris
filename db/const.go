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
	SetttingSizeLimiter           = "size_limiter"
	SettingTvNamingFormat         = "tv_naming_format"
	SettingMovieNamingFormat      = "movie_naming_format"
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
	R720p  Limiter `json:"720p"`
	R1080p Limiter `json:"1080p"`
	R2160p Limiter `json:"2160p"`
}

type Limiter struct {
	Max int `json:"max"`
	Min int `json:"min"`
}
