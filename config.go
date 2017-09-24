package main

import (
	"flag"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	config "gitlab.com/eunleem/gopack/config-v1"
)

type serverConfig struct {
	WebServer config.WebServerConfig `yaml:"server" json:"server"`
	MongoDb   config.MongoDbConfig   `yaml:"mongoDb" json:"mongoDb"`
	Redis     config.RedisConfig     `yaml:"redis" json:"redis"`
}

var conf serverConfig

var isDevMode = true
var configPath = ""

func init() {
	parseFlags()
	loadConfig()
}

func parseFlags() {
	log.Print("Parsing flags...")

	modePtr := flag.Bool("prod", false, "When dev mode is on, it uses config.dev.yaml for configuration")
	configPtr := flag.String("config", "", ".yaml config file to use")

	flag.Parse()

	isDevMode = !(*modePtr)
	configPath = *configPtr
}

func loadConfig() {
	confDirPath := "./configs/"
	confFileName := "server.prod.yaml"
	confPath := ""

	if isDevMode == true {
		confFileName = strings.Replace(confFileName, ".prod.yaml", ".dev.yaml", 1)
	}

	if configPath == "" {
		confPath = filepath.Join(confDirPath, confFileName)
	} else {
		confPath = configPath
	}

	log.Printf("Loading config file: \"%s\"\n", confPath)

	if err := config.LoadFile(&conf, confPath); err != nil {
		panic(err)
	}

	log.Print("Server Name: " + conf.WebServer.Name)
	log.Print("Version: " + conf.WebServer.Version)
	log.Print(conf.WebServer.Domain + ":" + strconv.Itoa(conf.WebServer.Port))
	log.Print(conf.WebServer.WebDir)
	log.Print(conf.WebServer.FullAddress)

	if isDevMode == true {
		log.Print("MongoDB Config")
		log.Print(conf.MongoDb.Host)
		log.Print(conf.MongoDb.Username)
		log.Print(conf.MongoDb.Password)

		log.Print("Redis Config")
		log.Print(conf.Redis.Host)
		log.Print(conf.Redis.Password)
	}
}
