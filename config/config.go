package config

import (
	"fmt"
	"os"

	toml "github.com/pelletier/go-toml"
)

type Config struct {
	App      AppConfiguration      `toml:"Application"`
	Database DatabaseConfiguration `toml:"Database"`
	Redis    RedisConfiguration    `toml:"Redis"`
}

type DatabaseConfiguration struct {
	Host     string `toml:"Host"`
	KeySpace string `toml:"KeySpace"`
}

type RedisConfiguration struct {
	Host string `toml:"Host"`
	Port int    `toml:"Port"`
}
type AppConfiguration struct {
	Port int `toml:"Port"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	fmt.Println(path)
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}

	tomlFile, err := toml.LoadReader(file)
	if err != nil {
		return config, err
	}

	err = tomlFile.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil

}
