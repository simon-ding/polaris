package db

var Version = "undefined"

const (
	SettingTmdbApiKey       = "tmdb_api_key"
	SettingLanguage         = "language"
	SettingJacketUrl        = "jacket_url"
	SettingJacketApiKey     = "jacket_api_key"
	SettingDownloadDir      = "download_dir"
	SettingLogLevel         = "log_level"
	SettingProxy            = "proxy"
	SettingPlexMatchEnabled = "plexmatch_enabled"
	SettingAllowQiangban    = "filter_qiangban"
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

type ResolutionType string

const JwtSerectKey = "jwt_secrect_key"
