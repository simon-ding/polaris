package db

const (
	SettingTmdbApiKey = "tmdb_api_key"
	SettingLanguage = "language"
	SettingJacketUrl = "jacket_url"
	SettingJacketApiKey = "jacket_api_key"
	SettingDownloadDir = "download_dir"
)

const (
	SettingAuthEnabled = "auth_enbled"
	SettingUsername = "auth_username"
	SettingPassword = "auth_password"
)

const (
	IndexerTorznabImpl = "torznab"
)

type ResolutionType string

const (
	Any ResolutionType = "any" 
	R720p ResolutionType = "720p"
	R1080p ResolutionType = "1080p"
	R4k ResolutionType = "4k"
)

func (r ResolutionType) String() string {
	return string(r)
}

const (
	ImplLocal = "local"
	ImplWebdav = "webdav"
)

func StorageImplementations() []string {
	return []string{
		ImplLocal,
		ImplWebdav,
	}
}