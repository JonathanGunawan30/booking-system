package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port                  int      `mapstructure:"PORT"`
	AppName               string   `mapstructure:"APP_NAME"`
	AppEnv                string   `mapstructure:"APP_ENV"`
	SignatureKey          string   `mapstructure:"SIGNATURE_KEY"`
	Database              Database `mapstructure:",squash"`
	RateLimiterMaxRequest float64  `mapstructure:"RATE_LIMITER_MAX_REQUEST"`
	RateLimiterTimeSecond int      `mapstructure:"RATE_LIMITER_TIME_SECOND"`
	JwtSecretKey          string   `mapstructure:"JWT_SECRET_KEY"`
	JwtExpirationTime     int      `mapstructure:"JWT_EXPIRATION_TIME"`
	JwtIssuer             string   `mapstructure:"JWT_ISSUER"`

	Timezone                  string `mapstructure:"TIMEZONE"`
	ConsulHTTPURL             string `mapstructure:"CONSUL_HTTP_URL"`
	ConsulHTTPKey             string `mapstructure:"CONSUL_HTTP_KEY"`
	ConsulHTTPToken           string `mapstructure:"CONSUL_HTTP_TOKEN"`
	ConsulWatchIntervalSecond int    `mapstructure:"CONSUL_WATCH_INTERVAL_SECONDS"`
}

type Database struct {
	Host                  string `mapstructure:"DB_HOST"`
	Port                  int    `mapstructure:"DB_PORT"`
	Name                  string `mapstructure:"DB_NAME"`
	Username              string `mapstructure:"DB_USERNAME"`
	Password              string `mapstructure:"DB_PASSWORD"`
	SSLMode               string `mapstructure:"DB_SSL_MODE"`
	MaxOpenConnections    int    `mapstructure:"DB_MAX_OPEN_CONNECTION"`
	MaxIdleConnections    int    `mapstructure:"DB_MAX_IDLE_CONNECTION"`
	MaxLifetimeConnection int    `mapstructure:"DB_MAX_LIFETIME"`
	MaxIdleTime           int    `mapstructure:"DB_MAX_IDLE_TIME"`
}

func bindAllEnvironmentVariables() {
	envs := []string{
		"PORT", "APP_NAME", "APP_ENV", "SIGNATURE_KEY",
		"RATE_LIMITER_MAX_REQUEST", "RATE_LIMITER_TIME_SECOND",
		"JWT_SECRET_KEY", "JWT_EXPIRATION_TIME", "JWT_ISSUER",
		"TIMEZONE", "CONSUL_HTTP_URL", "CONSUL_HTTP_KEY", "CONSUL_HTTP_TOKEN", "CONSUL_WATCH_INTERVAL_SECONDS",
		"DB_HOST", "DB_PORT", "DB_NAME", "DB_USERNAME", "DB_PASSWORD", "DB_SSL_MODE",
		"DB_MAX_OPEN_CONNECTION", "DB_MAX_IDLE_CONNECTION", "DB_MAX_LIFETIME", "DB_MAX_IDLE_TIME",
	}

	for _, env := range envs {
		viper.BindEnv(env)
	}
}

func LoadConfig() Config {
	viper.SetConfigFile(".env")

	_ = viper.ReadInConfig()

	viper.AutomaticEnv()

	bindAllEnvironmentVariables()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.Fatalf("failed to load config: %v", err)
	}

	validateConfig(cfg)

	return cfg
}

func validateConfig(cfg Config) {
	useConsul := cfg.ConsulHTTPURL != ""

	if useConsul {
		logrus.Info("using Consul as config source")
		return
	}

	logrus.Info("using ENV config — ensure .env is loaded and required variables are not empty")

}
