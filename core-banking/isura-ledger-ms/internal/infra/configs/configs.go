package configs

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Configs struct {
	ApplicationName string        `mapstructure:"application_name"`
	Env             string        `mapstructure:"env"`
	Servers         Servers       `mapstructure:"servers"`
	DataBase        DataBase      `mapstructure:"data_base"`
	OpenTelemetry   OpemTelemetry `mapstructure:"open_telemetry"`
	Version         string        `mapstructure:"version"`
}

type Servers struct {
	GRPC GRPC `mapstructure:"grpc"`
	HTTP HTTP `mapstructure:"http"`
}
type GRPC struct {
	Port string `mapstructure:"port"`
}

type HTTP struct {
	Port string `mapstructure:"port"`
}

type DataBase struct {
	Postgres Postgres `mapstructure:"postgres"`
}

type OpemTelemetry struct {
	Host string `mapstructure:"host"`
}

type Postgres struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxConnections  int32         `mapstructure:"max_connections"`
	MinConnections  int32         `mapstructure:"min_connections"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

func LoadConfig() *Configs {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}

	bindEnvs()

	var configs Configs
	err = viper.Unmarshal(&configs)
	if err != nil {
		return nil
	}

	os.Setenv("ENV", configs.Env)

	return &configs
}

// bindEnvs mapeia variáveis de ambiente para as chaves do config.json.
// Env vars têm precedência sobre o arquivo de configuração.
func bindEnvs() {
	viper.BindEnv("env", "APP_ENV")
	viper.BindEnv("servers.grpc.port", "GRPC_PORT")
	viper.BindEnv("servers.http.port", "HTTP_PORT")
	viper.BindEnv("data_base.postgres.host", "POSTGRES_HOST")
	viper.BindEnv("data_base.postgres.port", "POSTGRES_PORT")
	viper.BindEnv("data_base.postgres.user", "POSTGRES_USER")
	viper.BindEnv("data_base.postgres.password", "POSTGRES_PASSWORD")
	viper.BindEnv("data_base.postgres.database", "POSTGRES_DB")
	viper.BindEnv("data_base.postgres.ssl_mode", "POSTGRES_SSL_MODE")
	viper.BindEnv("data_base.postgres.max_connections", "POSTGRES_MAX_CONNECTIONS")
	viper.BindEnv("data_base.postgres.min_connections", "POSTGRES_MIN_CONNECTIONS")
	viper.BindEnv("data_base.postgres.max_conn_lifetime", "POSTGRES_MAX_CONN_LIFETIME")
	viper.BindEnv("data_base.postgres.max_conn_idle_time", "POSTGRES_MAX_CONN_IDLE_TIME")
	viper.BindEnv("open_telemetry.host", "OTEL_HOST")
}
