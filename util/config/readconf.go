package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Listen     string              `json:"listen"`
		Path       string              `json:"path"`
		Servicemap map[string][]string `json:"servicemap"`
		TLS        struct {
			Enabled  bool   `json:"enabled"`
			Keyfile  string `json:"keyfile"`
			Certfile string `json:"certfile"`
		} `json:"tls"`
	} `json:"server"`
	Client struct {
		Upstream  string              `json:"upstream"`
		Listenmap map[string][]string `json:"listenmap"`
	} `json:"client"`
}

func (c Config) IsServer() bool {
	return c.Server.Listen != ""
}

func (c Config) IsClient() bool {
	return c.Client.Upstream != ""
}

func LoadFile(path string) *Config {
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
	}
	var c Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Println(err)
	}
	return &c
}
