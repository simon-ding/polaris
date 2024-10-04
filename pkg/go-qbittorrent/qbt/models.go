package qbt

// BasicTorrent holds a basic torrent object from qbittorrent
type BasicTorrent struct {
	Category               string `json:"category"`
	CompletionOn           int64  `json:"completion_on"`
	Dlspeed                int    `json:"dlspeed"`
	Eta                    int    `json:"eta"`
	ForceStart             bool   `json:"force_start"`
	Hash                   string `json:"hash"`
	Name                   string `json:"name"`
	NumComplete            int    `json:"num_complete"`
	NumIncomplete          int    `json:"num_incomplete"`
	NumLeechs              int    `json:"num_leechs"`
	NumSeeds               int    `json:"num_seeds"`
	Priority               int    `json:"priority"`
	Progress               int    `json:"progress"`
	Ratio                  int    `json:"ratio"`
	SavePath               string `json:"save_path"`
	SeqDl                  bool   `json:"seq_dl"`
	Size                   int    `json:"size"`
	State                  string `json:"state"`
	SuperSeeding           bool   `json:"super_seeding"`
	Upspeed                int    `json:"upspeed"`
	FirstLastPiecePriority bool   `json:"f_l_piece_prio"`
}

// Torrent holds a torrent object from qbittorrent
// with more information than BasicTorrent
type Torrent struct {
	AdditionDate       int     `json:"addition_date"`
	Comment            string  `json:"comment"`
	CompletionDate     int     `json:"completion_date"`
	CreatedBy          string  `json:"created_by"`
	CreationDate       int     `json:"creation_date"`
	DlLimit            int     `json:"dl_limit"`
	DlSpeed            int     `json:"dl_speed"`
	DlSpeedAvg         int     `json:"dl_speed_avg"`
	Eta                int     `json:"eta"`
	LastSeen           int     `json:"last_seen"`
	NbConnections      int     `json:"nb_connections"`
	NbConnectionsLimit int     `json:"nb_connections_limit"`
	Peers              int     `json:"peers"`
	PeersTotal         int     `json:"peers_total"`
	PieceSize          int     `json:"piece_size"`
	PiecesHave         int     `json:"pieces_have"`
	PiecesNum          int     `json:"pieces_num"`
	Reannounce         int     `json:"reannounce"`
	SavePath           string  `json:"save_path"`
	SeedingTime        int     `json:"seeding_time"`
	Seeds              int     `json:"seeds"`
	SeedsTotal         int     `json:"seeds_total"`
	ShareRatio         float64 `json:"share_ratio"`
	TimeElapsed        int     `json:"time_elapsed"`
	TotalDl            int     `json:"total_downloaded"`
	TotalDlSession     int     `json:"total_downloaded_session"`
	TotalSize          int     `json:"total_size"`
	TotalUl            int     `json:"total_uploaded"`
	TotalUlSession     int     `json:"total_uploaded_session"`
	TotalWasted        int     `json:"total_wasted"`
	UpLimit            int     `json:"up_limit"`
	UpSpeed            int     `json:"up_speed"`
	UpSpeedAvg         int     `json:"up_speed_avg"`
}

type TorrentInfo struct {
	AddedOn           int64   `json:"added_on"`
	AmountLeft        int64   `json:"amount_left"`
	AutoTmm           bool    `json:"auto_tmm"`
	Availability      int64   `json:"availability"`
	Category          string  `json:"category"`
	Completed         int64   `json:"completed"`
	CompletionOn      int64   `json:"completion_on"`
	ContentPath       string  `json:"content_path"`
	DlLimit           int64   `json:"dl_limit"`
	Dlspeed           int64   `json:"dlspeed"`
	Downloaded        int64   `json:"downloaded"`
	DownloadedSession int64   `json:"downloaded_session"`
	Eta               int64   `json:"eta"`
	FLPiecePrio       bool    `json:"f_l_piece_prio"`
	ForceStart        bool    `json:"force_start"`
	Hash              string  `json:"hash"`
	LastActivity      int64   `json:"last_activity"`
	MagnetURI         string  `json:"magnet_uri"`
	MaxRatio          float64 `json:"max_ratio"`
	MaxSeedingTime    int64   `json:"max_seeding_time"`
	Name              string  `json:"name"`
	NumComplete       int64   `json:"num_complete"`
	NumIncomplete     int64   `json:"num_incomplete"`
	NumLeechs         int64   `json:"num_leechs"`
	NumSeeds          int64   `json:"num_seeds"`
	Priority          int64   `json:"priority"`
	Progress          float64 `json:"progress"`
	Ratio             float64 `json:"ratio"`
	RatioLimit        int64   `json:"ratio_limit"`
	SavePath          string  `json:"save_path"`
	SeedingTimeLimit  int64   `json:"seeding_time_limit"`
	SeenComplete      int64   `json:"seen_complete"`
	SeqDl             bool    `json:"seq_dl"`
	Size              int64   `json:"size"`
	State             string  `json:"state"`
	SuperSeeding      bool    `json:"super_seeding"`
	Tags              string  `json:"tags"`
	TimeActive        int64   `json:"time_active"`
	TotalSize         int64   `json:"total_size"`
	Tracker           string  `json:"tracker"`
	TrackersCount     int64   `json:"trackers_count"`
	UpLimit           int64   `json:"up_limit"`
	Uploaded          int64   `json:"uploaded"`
	UploadedSession   int64   `json:"uploaded_session"`
	Upspeed           int64   `json:"upspeed"`
}

// Tracker holds a tracker object from qbittorrent
type Tracker struct {
	Msg           string `json:"msg"`
	NumPeers      int    `json:"num_peers"`
	NumSeeds      int    `json:"num_seeds"`
	NumLeeches    int    `json:"num_leeches"`
	NumDownloaded int    `json:"num_downloaded"`
	Tier          int    `json:"tier"`
	Status        int    `json:"status"`
	URL           string `json:"url"`
}

// WebSeed holds a webseed object from qbittorrent
type WebSeed struct {
	URL string `json:"url"`
}

// TorrentFile holds a torrent file object from qbittorrent
type TorrentFile struct {
	Index        int     `json:"index"`
	IsSeed       bool    `json:"is_seed"`
	Name         string  `json:"name"`
	Availability float32 `json:"availability"`
	Priority     int     `json:"priority"`
	Progress     int     `json:"progress"`
	Size         int     `json:"size"`
	PieceRange   []int   `json:"piece_range"`
}

// Sync holds the sync response struct which contains
// the server state and a map of infohashes to Torrents
type Sync struct {
	Categories  []string `json:"categories"`
	FullUpdate  bool     `json:"full_update"`
	Rid         int      `json:"rid"`
	ServerState struct {
		ConnectionStatus  string `json:"connection_status"`
		DhtNodes          int    `json:"dht_nodes"`
		DlInfoData        int    `json:"dl_info_data"`
		DlInfoSpeed       int    `json:"dl_info_speed"`
		DlRateLimit       int    `json:"dl_rate_limit"`
		Queueing          bool   `json:"queueing"`
		RefreshInterval   int    `json:"refresh_interval"`
		UpInfoData        int    `json:"up_info_data"`
		UpInfoSpeed       int    `json:"up_info_speed"`
		UpRateLimit       int    `json:"up_rate_limit"`
		UseAltSpeedLimits bool   `json:"use_alt_speed_limits"`
	} `json:"server_state"`
	Torrents map[string]Torrent `json:"torrents"`
}

type BuildInfo struct {
	QTVersion         string `json:"qt"`
	LibtorrentVersion string `json:"libtorrent"`
	BoostVersion      string `json:"boost"`
	OpenSSLVersion    string `json:"openssl"`
	AppBitness        string `json:"bitness"`
}

type Preferences struct {
	Locale                             string                 `json:"locale"`
	CreateSubfolderEnabled             bool                   `json:"create_subfolder_enabled"`
	StartPausedEnabled                 bool                   `json:"start_paused_enabled"`
	AutoDeleteMode                     int                    `json:"auto_delete_mode"`
	PreallocateAll                     bool                   `json:"preallocate_all"`
	IncompleteFilesExt                 bool                   `json:"incomplete_files_ext"`
	AutoTMMEnabled                     bool                   `json:"auto_tmm_enabled"`
	TorrentChangedTMMEnabled           bool                   `json:"torrent_changed_tmm_enabled"`
	SavePathChangedTMMEnabled          bool                   `json:"save_path_changed_tmm_enabled"`
	CategoryChangedTMMEnabled          bool                   `json:"category_changed_tmm_enabled"`
	SavePath                           string                 `json:"save_path"`
	TempPathEnabled                    bool                   `json:"temp_path_enabled"`
	TempPath                           string                 `json:"temp_path"`
	ScanDirs                           map[string]interface{} `json:"scan_dirs"`
	ExportDir                          string                 `json:"export_dir"`
	ExportDirFin                       string                 `json:"export_dir_fin"`
	MailNotificationEnabled            string                 `json:"mail_notification_enabled"`
	MailNotificationSender             string                 `json:"mail_notification_sender"`
	MailNotificationEmail              string                 `json:"mail_notification_email"`
	MailNotificationSMPTP              string                 `json:"mail_notification_smtp"`
	MailNotificationSSLEnabled         bool                   `json:"mail_notification_ssl_enabled"`
	MailNotificationAuthEnabled        bool                   `json:"mail_notification_auth_enabled"`
	MailNotificationUsername           string                 `json:"mail_notification_username"`
	MailNotificationPassword           string                 `json:"mail_notification_password"`
	AutorunEnabled                     bool                   `json:"autorun_enabled"`
	AutorunProgram                     string                 `json:"autorun_program"`
	QueueingEnabled                    bool                   `json:"queueing_enabled"`
	MaxActiveDls                       int                    `json:"max_active_downloads"`
	MaxActiveTorrents                  int                    `json:"max_active_torrents"`
	MaxActiveUls                       int                    `json:"max_active_uploads"`
	DontCountSlowTorrents              bool                   `json:"dont_count_slow_torrents"`
	SlowTorrentDlRateThreshold         int                    `json:"slow_torrent_dl_rate_threshold"`
	SlowTorrentUlRateThreshold         int                    `json:"slow_torrent_ul_rate_threshold"`
	SlowTorrentInactiveTimer           int                    `json:"slow_torrent_inactive_timer"`
	MaxRatioEnabled                    bool                   `json:"max_ratio_enabled"`
	MaxRatio                           float64                `json:"max_ratio"`
	MaxRatioAct                        bool                   `json:"max_ratio_act"`
	ListenPort                         int                    `json:"listen_port"`
	UPNP                               bool                   `json:"upnp"`
	RandomPort                         bool                   `json:"random_port"`
	DlLimit                            int                    `json:"dl_limit"`
	UlLimit                            int                    `json:"up_limit"`
	MaxConnections                     int                    `json:"max_connec"`
	MaxConnectionsPerTorrent           int                    `json:"max_connec_per_torrent"`
	MaxUls                             int                    `json:"max_uploads"`
	MaxUlsPerTorrent                   int                    `json:"max_uploads_per_torrent"`
	UTPEnabled                         bool                   `json:"enable_utp"`
	LimitUTPRate                       bool                   `json:"limit_utp_rate"`
	LimitTCPOverhead                   bool                   `json:"limit_tcp_overhead"`
	LimitLANPeers                      bool                   `json:"limit_lan_peers"`
	AltDlLimit                         int                    `json:"alt_dl_limit"`
	AltUlLimit                         int                    `json:"alt_up_limit"`
	SchedulerEnabled                   bool                   `json:"scheduler_enabled"`
	ScheduleFromHour                   int                    `json:"schedule_from_hour"`
	ScheduleFromMin                    int                    `json:"schedule_from_min"`
	ScheduleToHour                     int                    `json:"schedule_to_hour"`
	ScheduleToMin                      int                    `json:"schedule_to_min"`
	SchedulerDays                      int                    `json:"scheduler_days"`
	DHTEnabled                         bool                   `json:"dht"`
	DHTSameAsBT                        bool                   `json:"dhtSameAsBT"`
	DHTPort                            int                    `json:"dht_port"`
	PexEnabled                         bool                   `json:"pex"`
	LSDEnabled                         bool                   `json:"lsd"`
	Encryption                         int                    `json:"encryption"`
	AnonymousMode                      bool                   `json:"anonymous_mode"`
	ProxyType                          int                    `json:"proxy_type"`
	ProxyIP                            string                 `json:"proxy_ip"`
	ProxyPort                          int                    `json:"proxy_port"`
	ProxyPeerConnections               bool                   `json:"proxy_peer_connections"`
	ForceProxy                         bool                   `json:"force_proxy"`
	ProxyAuthEnabled                   bool                   `json:"proxy_auth_enabled"`
	ProxyUsername                      string                 `json:"proxy_username"`
	ProxyPassword                      string                 `json:"proxy_password"`
	IPFilterEnabled                    bool                   `json:"ip_filter_enabled"`
	IPFilterPath                       string                 `json:"ip_filter_path"`
	IPFilterTrackers                   string                 `json:"ip_filter_trackers"`
	WebUIDomainList                    string                 `json:"web_ui_domain_list"`
	WebUIAddress                       string                 `json:"web_ui_address"`
	WebUIPort                          int                    `json:"web_ui_port"`
	WebUIUPNPEnabled                   bool                   `json:"web_ui_upnp"`
	WebUIUsername                      string                 `json:"web_ui_username"`
	WebUIPassword                      string                 `json:"web_ui_password"`
	WebUICSRFProtectionEnabled         bool                   `json:"web_ui_csrf_protection_enabled"`
	WebUIClickjackingProtectionEnabled bool                   `json:"web_ui_clickjacking_protection_enabled"`
	BypassLocalAuth                    bool                   `json:"bypass_local_auth"`
	BypassAuthSubnetWhitelistEnabled   bool                   `json:"bypass_auth_subnet_whitelist_enabled"`
	BypassAuthSubnetWhitelist          string                 `json:"bypass_auth_subnet_whitelist"`
	AltWebUIEnabled                    bool                   `json:"alternative_webui_enabled"`
	AltWebUIPath                       string                 `json:"alternative_webui_path"`
	UseHTTPS                           bool                   `json:"use_https"`
	SSLKey                             string                 `json:"ssl_key"`
	SSLCert                            string                 `json:"ssl_cert"`
	DynDNSEnabled                      bool                   `json:"dyndns_enabled"`
	DynDNSService                      int                    `json:"dyndns_service"`
	DynDNSUsername                     string                 `json:"dyndns_username"`
	DynDNSPassword                     string                 `json:"dyndns_password"`
	DynDNSDomain                       string                 `json:"dyndns_domain"`
	RSSRefreshInterval                 int                    `json:"rss_refresh_interval"`
	RSSMaxArtPerFeed                   int                    `json:"rss_max_articles_per_feed"`
	RSSProcessingEnabled               bool                   `json:"rss_processing_enabled"`
	RSSAutoDlEnabled                   bool                   `json:"rss_auto_downloading_enabled"`
}

// Log
type Log struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
	Type      int    `json:"type"`
}

// PeerLog
type PeerLog struct {
	ID        int    `json:"id"`
	IP        string `json:"ip"`
	Blocked   bool   `json:"blocked"`
	Timestamp int    `json:"timestamp"`
	Reason    string `json:"reason"`
}

// Info
type Info struct {
	ConnectionStatus  string `json:"connection_status"`
	DHTNodes          int    `json:"dht_nodes"`
	DlInfoData        int    `json:"dl_info_data"`
	DlInfoSpeed       int    `json:"dl_info_speed"`
	DlRateLimit       int    `json:"dl_rate_limit"`
	UlInfoData        int    `json:"up_info_data"`
	UlInfoSpeed       int    `json:"up_info_speed"`
	UlRateLimit       int    `json:"up_rate_limit"`
	Queueing          bool   `json:"queueing"`
	UseAltSpeedLimits bool   `json:"use_alt_speed_limits"`
	RefreshInterval   int    `json:"refresh_interval"`
}

type TorrentsOptions struct {
	Filter   *string  // all, downloading, completed, paused, active, inactive => optional
	Category *string  // => optional
	Sort     *string  // => optional
	Reverse  *bool    // => optional
	Limit    *int     // => optional (no negatives)
	Offset   *int     // => optional (negatives allowed)
	Hashes   []string // separated by | => optional
}

// Category of torrent
type Category struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

// Categories mapping
type Categories struct {
	Category map[string]Category
}

// LoginOptions contains all options for /login endpoint
type LoginOptions struct {
	Username string
	Password string
}

// AddTrackersOptions contains all options for /addTrackers endpoint
type AddTrackersOptions struct {
	Hash     string
	Trackers []string
}

// EditTrackerOptions contains all options for /editTracker endpoint
type EditTrackerOptions struct {
	Hash    string
	OrigURL string
	NewURL  string
}

// RemoveTrackersOptions contains all options for /removeTrackers endpoint
type RemoveTrackersOptions struct {
	Hash     string
	Trackers []string
}

type DownloadOptions struct {
	Savepath                   *string
	Cookie                     *string
	Category                   *string
	SkipHashChecking           *bool
	Paused                     *bool
	RootFolder                 *bool
	Rename                     *string
	UploadSpeedLimit           *int
	DownloadSpeedLimit         *int
	SequentialDownload         *bool
	AutomaticTorrentManagement *bool
	FirstLastPiecePriority     *bool
}

type InfoOptions struct {
	Filter   *string
	Category *string
	Sort     *string
	Reverse  *bool
	Limit    *int
	Offset   *int
	Hashes   []string
}

type PriorityValues int

const (
	Do_not_download  PriorityValues = 0
	Normal_priority  PriorityValues = 1
	High_priority    PriorityValues = 6
	Maximal_priority PriorityValues = 7
)
