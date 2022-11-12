package config

type (
	Config struct {
		App `yaml:"app"`
		Log `yaml:"logger"`
	}

	App struct {
		Name                     string  `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version                  string  `env-required:"true" yaml:"version" env:"APP_VERSION"`
		IdApp                    string  `env-required:"true" yaml:"id_app" env:"ID_APP"`
		AppIcon                  string  `env-required:"true" yaml:"app_icon" env:"APP_ICON"`
		AppImagePathConnected    string  `env-required:"true" yaml:"image_path_connected" env:"IMAGE_PATH_CONNECTED"`
		AppImagePathBlink        string  `env-required:"true" yaml:"image_path_blink" env:"IMAGE_PATH_BLINK"`
		AppImagePathOpen         string  `env-required:"true" yaml:"image_path_open" env:"IMAGE_PATH_OPEN"`
		AppImagePathDisconnected string  `env-required:"true" yaml:"image_path_disconnected" env:"IMAGE_PATH_DISCONNECTED"`
		TempDir                  string  `env-required:"true" yaml:"temp_dir" env:"TEMP_DIR"`
		ConfigsPath              string  `env-required:"true" yaml:"configs_path" env:"CONFIGS_PATH"`
		StopTimeout              int     `env-required:"true" yaml:"stop_time" env:"STOP_TIME"`
		Height                   float32 `env-required:"true" yaml:"height_app" env:"HEIGHT_APP"`
		Width                    float32 `env-required:"true" yaml:"width_app" env:"WIDTH_APP"`
		UseSystemd               bool    `yaml:"use_systemd" env:"USE_SYSTEMD"`
		TunPersist               bool    `yaml:"tun_persist" env:"TUN_PERSIST"`
		ClockTicks               int     `env-required:"true" yaml:"clock_ticks_time_ms" env:"CLOCK_TICKS_TIME_MS"`
		CountReconn              int     `env-required:"true" yaml:"count_reconnection_attempt" env:"COUNT_RECCONECTION_ATTEMPT"`
		VerbLogs                 bool    `yaml:"verbose_logs" env:"VERB_LOGS"`
		ConnectTimeout           int     `env-required:"true" yaml:"conn_timeout" env:"CONN_TIMEOUT"`
		CompressionMode          string  `env-required:"true" yaml:"mode" env:"MODE"`
		CheckDisableClientCert   bool    `yaml:"disable_cert" env:"DISABLE_CERT"`
	}

	Log struct {
		Level         string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
		Logfile       string `yaml:"log_file" env:"LOG_FILE" `
		LogFileCount  int    `yaml:"log_file_count" env:"LOG_FILE_COUNT"`
		Log_File_Size int64  `yaml:"log_file_size" env:"LOG_FILE_SIZE"`
	}
)
