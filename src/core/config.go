package core

import (
	"encoding/json"
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var Config config

type config struct {
	BasePath   string
	Logger     Logger
	Database   map[string]Database
	Listen     string
	HttpListen string
	APISecret  string
	AppEnv     string
	AppName    string
	Redis      map[string]RedisConfig
}

type Logger struct {
	Debug        bool
	OutFile      bool
	LogPath      string
	MaxAge       time.Duration //日志最大保存时间单位小时
	RotationTime time.Duration //日志切割时间间隔单位小时
}

type Database struct {
	DriverName     string
	DataSourceName string
	MaxOpenNum     int
	MaxIdleNum     int
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func (c *config) Init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("error get os.Getwd(): ", err.Error())
	}

	var configPath string
	flag.StringVar(&configPath, "c", dir+"/src/config.local.json", "App conifg path")
	flag.Parse()
	log.WithFields(log.Fields{"configPath": configPath}).Info("Config Init")

	file, err := os.Open(configPath)
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&c)
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}

	cwd, _ := os.Getwd()
	if Config.BasePath != "" {
		cwd = Config.BasePath
	} else {
		Config.BasePath = cwd
	}
}
