package main

import (
	"flag"
	"log"
	"os"

	"github.com/sileader/llama-run/config"
)

func main() {
	configFile := flag.String("config", "/etc/llama-run/config.yaml", "Path to the config file")
	flag.Parse()

	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}
}

func loadConfig(file string) (*config.Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return config.UnmarshalConfig(data)
}
