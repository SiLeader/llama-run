package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/sileader/llama-run/builder"
	"github.com/sileader/llama-run/config"
	"github.com/sileader/llama-run/downloader"
)

func main() {
	configFile := flag.String("config", "/etc/llama-run/config.yaml", "Path to the config file")
	dryRun := flag.Bool("dry-run", false, "Dry run")
	flag.Parse()

	cfg, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	logLevel := slog.LevelInfo
	if err := logLevel.UnmarshalText([]byte(strings.ToLower(cfg.Log.Level))); err != nil {
		log.Printf("Invalid log level, falling back to info: %v\n", err)
	}
	slog.SetLogLoggerLevel(logLevel)

	dlb := downloader.NewBuilder(cfg.Downloader)

	ctx := context.Background()
	llamaServer, err := builder.NewLlamaServerApplicationBuilder(ctx, cfg.LlamaServer, dlb)
	if err != nil {
		log.Fatalf("Failed to create llama server application builder: %v\n", err)
	}
	if err := cfg.Visit(llamaServer); err != nil {
		log.Fatalf("Failed to visit config: %v\n", err)
	}

	if *dryRun {
		slog.Info("Dry run mode enabled. No actions will be performed.")
		return
	}

	slog.Info("Starting llama-server")
	if err := llamaServer.Exec(); err != nil {
		log.Fatalf("Failed to start llama-server: %v\n", err)
	}
}

func loadConfig(file string) (*config.Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return config.UnmarshalConfig(data)
}
