package config

import (
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Configuration  for chain service
type Configuration struct {
	ListenPort      int                    `json:"listen_port" mapstructure:"listen_port"`
	ShutdownTimeout time.Duration          `json:"shutdown_timeout" mapstructure:"shutdown_timeout"`
	ReadTimeout     time.Duration          `json:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout    time.Duration          `json:"write_timeout" mapstructure:"write_timeout"`
	TimeScale       TimeScaleConfiguration `json:"timescale" mapstructure:"timescale"`
	ThorChain       ThorChainConfiguration `json:"thorchain" mapstructure:"thorchain"`
	LogLevel        string                 `json:"log_level" mapstructure:"log_level"`
	FullNodes       []NodeProxy            `json:"full_nodes" mapstructure:"full_nodes"`
}

type TimeScaleConfiguration struct {
	Host                  string        `json:"host" mapstructure:"host"`
	Port                  int           `json:"port" mapstructure:"port"`
	UserName              string        `json:"user_name" mapstructure:"user_name"`
	Password              string        `json:"password" mapstructure:"password"`
	Database              string        `json:"database" mapstructure:"database"`
	Sslmode               string        `json:"sslmode" mapstructure:"sslmode"`
	MigrationsDir         string        `json:"migrationsDir" mapstructure:"migrationsDir"`
	MaxConnections        int           `json:"max_connections" mapstructure:"max_connections"`
	ConnectionMaxLifetime time.Duration `json:"connection_max_lifetime" mapstructure:"connection_max_lifetime"`
}

type ThorChainConfiguration struct {
	Scheme                      string        `json:"scheme" mapstructure:"scheme"`
	Host                        string        `json:"host" mapstructure:"host"`
	RPCHost                     string        `json:"rpc_host" mapstructure:"rpc_host"`
	ReadTimeout                 time.Duration `json:"read_timeout" mapstructure:"read_timeout"`
	EnableScan                  bool          `json:"enable_scan" mapstructure:"enable_scan"` // TODO: Remove this field
	NoEventsBackoff             time.Duration `json:"no_events_backoff" mapstructure:"no_events_backoff"`
	ProxiedWhitelistedEndpoints []string      `json:"proxied_whitelisted_endpoints" mapstructure:"proxied_whitelisted_endpoints"`
	CacheTTL                    time.Duration `json:"cache_ttl" mapstructure:"cache_ttl"`
	CacheCleanup                time.Duration `json:"cache_cleanup" mapstructure:"cache_cleanup"`
	CacheSize                   int           `json:"cache_size" mapstructure:"cache_size"`
}

type NodeProxy struct {
	Chain         string `json:"chain" mapstructure:"chain"`
	Target        string `json:"target" mapstructure:"target"`
	WebsocketPath string `json:"websocket_path" mapstructure:"websocket_path"`
}

func applyDefaultConfig() {
	viper.SetDefault("read_timeout", "30s")
	viper.SetDefault("write_timeout", "30s")
	viper.SetDefault("thorchain.read_timeout", "10s")
	viper.SetDefault("thorchain.no_events_backoff", "5s")
	viper.SetDefault("thorchain.cache_ttl", "5s")
	viper.SetDefault("thorchain.cache_cleanup", "10s")
	viper.SetDefault("thorchain.cache_size", "40000")
	viper.SetDefault("thorchain.scan_start_pos", 1)
	viper.SetDefault("timescale.max_connections", 25)
	viper.SetDefault("timescale.connection_max_lifetime", time.Minute*5)
}

func LoadConfiguration(file string) (*Configuration, error) {
	applyDefaultConfig()
	var cfg Configuration
	viper.SetConfigName(strings.TrimRight(path.Base(file), ".json"))
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Dir(file))
	if err := viper.ReadInConfig(); nil != err {
		return nil, errors.Wrap(err, "failed to read from config file")
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.Unmarshal(&cfg); nil != err {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}
	return &cfg, nil
}
