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
	Binance         BinanceConfiguration   `json:"binance" mapstructure:"binance"`
	IsTestNet       bool                   `json:"is_testnet" mapstructure:"is_testnet"`
	LogLevel        string                 `json:"log_level" mapstructure:"log_level"`
}

type TimeScaleConfiguration struct {
	Host          string `json:"host" mapstructure:"host"`
	Port          int    `json:"port" mapstructure:"port"`
	UserName      string `json:"user_name" mapstructure:"user_name"`
	Password      string `json:"password" mapstructure:"password"`
	Database      string `json:"database" mapstructure:"database"`
	Sslmode       string `json:"sslmode" mapstructure:"sslmode"`
	MigrationsDir string `json:"migrationsDir" mapstructure:"migrationsDir"`
}

type ThorChainConfiguration struct {
	Scheme               string        `json:"scheme" mapstructure:"scheme"`
	Host                 string        `json:"host" mapstructure:"host"`
	RPCHost              string        `json:"rpc_host" mapstructure:"rpc_host"`
	ReadTimeout          time.Duration `json:"read_timeout" mapstructure:"read_timeout"`
	EnableScan           bool          `json:"enable_scan" mapstructure:"enable_scan"`
	NoEventsBackoff      time.Duration `json:"no_events_backoff" mapstructure:"no_events_backoff"`
	WhitelistedEndpoints []string      `json:"whitelisted_endpoints" mapstructure:"whitelisted_endpoints"`
}

// BinanceConfiguration settings for binance client
type BinanceConfiguration struct {
	DEXHost              string        `json:"dex_host" mapstructure:"dex_host"`
	Scheme               string        `json:"scheme" mapstructure:"scheme"`
	RequestTimeout       time.Duration `json:"request_timeout" mapstructure:"request_timeout"`
	MarketsCacheDuration time.Duration `json:"markets_cache_duration" mapstructure:"markets_cache_duration"`
	TokensCacheDuration  time.Duration `json:"tokens_cache_duration" mapstructure:"tokens_cache_duration"`
	// FullNodeHost         string        `json:"full_node_host" mapstructure:"full_node_host"`
	FullNodeScheme string `json:"full_node_scheme" mapstructure:"full_node_scheme"`
	IsTestNet      bool   `json:"is_testnet" mapstructure:"is_testnet"`
}

func applyDefaultConfig() {
	viper.SetDefault("read_timeout", "30s")
	viper.SetDefault("write_timeout", "30s")
	viper.SetDefault("thorchain.read_timeout", "10s")
	viper.SetDefault("thorchain.no_events_backoff", "30s")
	viper.SetDefault("thorchain.scan_start_pos", 1)
	viper.SetDefault("binance.request_timeout", "30s")
	viper.SetDefault("binance.markets_cache_duration", "24h")
	viper.SetDefault("binance.tokens_cache_duration", "24h")
}

func LoadConfiguration(file string) (*Configuration, error) {
	applyDefaultConfig()
	var cfg Configuration
	viper.SetConfigName(strings.TrimRight(path.Base(file), ".json"))
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Dir(file))
	if err := viper.ReadInConfig(); nil != err {
		return nil, errors.Wrap(err, "fail to read from config file")
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.Unmarshal(&cfg); nil != err {
		return nil, errors.Wrap(err, "fail to unmarshal")
	}
	return &cfg, nil
}
