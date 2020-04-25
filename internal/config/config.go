package config

import (
	"log"
	"strings"

	"github.com/mpuzanov/sysmonitor/pkg/logger"
	"github.com/spf13/viper"
)

//Config Структура файла с конфигурацией
type Config struct {
	Log      logger.LogConf `yaml:"log" mapstructure:"log"`
	GRPCAddr string         `yaml:"grpc_listen" mapstructure:"grpc_listen"`
}

// LoadConfig Загрузка конфигурации из файла
func LoadConfig(filePath string) (*Config, error) {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("log.level", "info")
	viper.SetDefault("grpc_listen", "localhost:50051")

	if filePath != "" {
		log.Printf("Parsing config: %s\n", filePath)
		viper.SetConfigFile(filePath)
		viper.SetConfigType("yaml")
		//log.Println(viper.ConfigFileUsed())
		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("Config file is not specified.")
	}
	//log.Println(viper.AllSettings())

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	if config.Log.Level == "debug" {
		log.Printf("config: %+v", config)
	}
	return &config, nil
}
