package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string           `mapstructure:"env" yaml:"env"`
	HTTPServer HTTPServerConfig `mapstructure:"http_server" yaml:"http_server"`
	DataBase   DBConfig         `mapstructure:"database" yaml:"database"`
	JWT        JWTConfig        `mapstructure:"jwt" yaml:"jwt"`
	Kafka      KafkaCfg         `mapstructure:"kafka" yaml:"kafka"`
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type HTTPServerConfig struct {
	Host              string        `mapstructure:"host" yaml:"host"`
	Port              string        `mapstructure:"port" yaml:"port"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout" yaml:"read_header_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout"`
}

type DBConfig struct {
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	DBName   string `mapstructure:"db_name" yaml:"db_name"`
	SSLMode  string `mapstructure:"ssl_mode" yaml:"ssl_mode"`
}

type KafkaCfg struct {
	Brokers      []string `mapstructure:"brokers" yaml:"brokers"`
	RequiredAcks int      `mapstructure:"required_acks" yaml:"required_acks"`
	Async        bool     `mapstructure:"async" yaml:"async"`
}

func New() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := setConfigEnv(&cfg); err != nil {
		return nil, fmt.Errorf("set config env: %w", err)
	}

	return &cfg, nil
}

func setConfigEnv(cfg *Config) error {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read .env: %w", err)
	}

	cfg.DataBase.Password = viper.GetString("postgres_password")
	if cfg.DataBase.Password == "" {
		return fmt.Errorf("db_password is required")
	}

	cfg.JWT.Secret = viper.GetString("jwt_secret")
	if cfg.JWT.Secret == "" {
		return fmt.Errorf("jwt_secret is required")
	}

	cfg.JWT.TTL = viper.GetDuration("jwt_ttl")
	if cfg.JWT.TTL == 0 {
		cfg.JWT.TTL = 24 * time.Hour
	}

	return nil
}
