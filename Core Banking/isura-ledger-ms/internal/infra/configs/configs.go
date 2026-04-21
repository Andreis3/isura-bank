package configs

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Configs struct {
	ApplicationName string   `mapstructure:"application_name"`
	Env             string   `mapstructure:"env"`
	Servers         Servers  `mapstructure:"servers"`
	DataBase        DataBase `mapstructure:"data_base"`
	Version         string   `mapstructure:"version"`
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

type Postgres struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxConnections  int32         `mapstructure:"max_connections"`
	MinConnections  int32         `mapstructure:"max_connections"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

func LoadConfig() *Configs {
	viper.SetConfigName("config") // nome do arquivo sem extensão
	viper.SetConfigType("json")   // tipo do arquivo
	viper.AddConfigPath(".")      // procura na raiz do projeto

	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}

	var configs Configs
	err = viper.Unmarshal(&configs)
	if err != nil {
		return nil
	}

	os.Setenv("ENV", configs.Env)

	return &configs
}
