package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Sev    Server  `mapstructure:"server"`
	Drv    Driver  `mapstructure:"driver"`
	LogCon LogConf `mapstructure:"log"`
	TagCon []string `mapstructure:"tag"`
}

type Server struct {
	Port		string
}

type Driver struct {
	Username			string
	Password 			string
	Network 			string
	Server 				string
	Port 				string
	Database 			string
	InterpolateParams 	bool
}

type LogConf struct {
	FilePath			string
	FileName 			string
}

var Conf = new(Config)

func InitConfig() {
	viper.AddConfigPath("./conf")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	}
	// 监控配置文件变化
	viper.WatchConfig()
	//配置文件发生变化后同步到全局变量Conf
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Config file changed, reloading...")
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
		}
	})


	logrus.Info("Config init successfully")
}











