package config

import (
	"log"
	"strings"

	"github.com/mpuzanov/sysmonitor/pkg/logger"
	"github.com/spf13/viper"
)

//Config Структура файла с конфигурацией
type Config struct {
	Host      string         `yaml:"host" mapstructure:"host"`
	Port      string         `yaml:"port" mapstructure:"port"`
	Log       logger.LogConf `yaml:"log" mapstructure:"log"`
	Collector CollectorConf  `yaml:"collector" mapstructure:"collector"`
}

// CollectorConf .
type CollectorConf struct {
	Timeout  int `yaml:"timeout" mapstructure:"timeout"`
	Category struct {
		LoadSystem  bool `yaml:"load_system" mapstructure:"load_system"`
		LoadCPU     bool `yaml:"load_cpu" mapstructure:"load_cpu"`
		LoadDisk    bool `yaml:"load_disk" mapstructure:"load_disk"`
		TopTalkers  bool `yaml:"top_talkers" mapstructure:"top_talkers"`
		StatNetwork bool `yaml:"stat_network" mapstructure:"stat_network"`
	} `yaml:"category" mapstructure:"category"`
}

// LoadConfig Загрузка конфигурации из файла
func LoadConfig(filePath string) (*Config, error) {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("host", "localhost")
	viper.SetDefault("collector.timeout", 5)

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