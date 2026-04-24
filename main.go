package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/sileader/llama-run/builder"
	"github.com/sileader/llama-run/config"
	"github.com/sileader/llama-run/downloader"
)

func main() {
	configFile := flag.String("config", "/etc/llama-run/config.yaml", "Path to the config file")
	flag.Parse()

	cfg, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	dlb := downloader.NewBuilder(cfg.Downloader)

	ctx := context.Background()
	llamaServer := builder.NewLlamaServerApplicationBuilder(ctx, cfg.LlamaServer, dlb)
	if err := cfg.Visit(llamaServer); err != nil {
		log.Fatalln(err)
	}

	if err := llamaServer.Exec(); err != nil {
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
