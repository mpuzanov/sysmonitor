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

// CollectorConf подсистемы сбора статистики
type CollectorConf struct {
	// timeout для сбора информации в системе
	Timeout int `yaml:"timeout" mapstructure:"timeout"`
	// Category подсистемы сбора информации
	Category CategoryConf `yaml:"category" mapstructure:"category"`
}

// CategoryConf настройки подсистем сбора информации о системе
type CategoryConf struct {
	// LoadCPU подсистема сбора информации по загрузке системы
	LoadSystem bool `yaml:"load_system" mapstructure:"load_system"`
	// LoadCPU подсистема сбора информации по CPU
	LoadCPU bool `yaml:"load_cpu" mapstructure:"load_cpu"`
	// LoadDisk подсистема сбора информации по дискам
	LoadDisk bool `yaml:"load_disk" mapstructure:"load_disk"`
	// TopTalkers подсистема сбора информации по трафику сети
	TopTalkers bool `yaml:"top_talkers" mapstructure:"top_talkers"`
	// StatNetwork подсистема сбора информации по статистеке сетевым соединениям
	StatNetwork bool `yaml:"stat_network" mapstructure:"stat_network"`
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
