package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

const defaultPath = "./config/config.yaml"

type Config struct {
	Application     Application `yaml:"application"`
	FileStorage     FileStorage `yaml:"fileStorage"`
	MetaStorage     MetaStorage `yaml:"metaStorage"`
	ChunkSize       int         `yaml:"chunkSize"` // bytes
	LoadWorkersCnt  int         `yaml:"loadWorkersCnt"`
	FetchWorkersCnt int         `yaml:"FetchWorkersCnt"`
}

type Application struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type FileStorage struct {
	DirectoryName string `yaml:"directoryName"`
}

type MetaStorage struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"DBName"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		path = defaultPath
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("file is not exists")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config file")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
