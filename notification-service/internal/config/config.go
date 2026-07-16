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
	Kafka      KafkaConfig      `mapstructure:"kafka" yaml:"kafka"`
	MailGun    MailGunConfig    `mapstructure:"mailgun" yaml:"mailgun"`
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

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers" yaml:"brokers"`
	GroupID string   `mapstructure:"group_id" yaml:"group_id"`
}

type MailGunConfig struct {
	APIKey string `mapstructure:"api_key" yaml:"api_key"`
	Domain string `mapstructure:"domain" yaml:"domain"`
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

	cfg.MailGun.APIKey = viper.GetString("mailgun_api_key")
	if cfg.MailGun.APIKey == "" {
		return fmt.Errorf("mailgun_api_key is required")
	}

	cfg.MailGun.Domain = viper.GetString("mailgun_domain")
	if cfg.MailGun.Domain == "" {
		return fmt.Errorf("mailgun_domain is required")
	}

	return nil
}
