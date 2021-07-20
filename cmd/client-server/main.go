package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/holmes89/chickaree-db/chickaree"
	"github.com/holmes89/chickaree-db/chickaree/redis"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

func main() {

	cfg, err := LoadConfiguration()
	if err != nil {
		log.Fatal().Msg("unable to load config")
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}
	conn, err := grpc.Dial(cfg.StorageServer, opts...)
	if err != nil {
		log.Fatal().Err(err).Str("url", cfg.StorageServer).Msg("failed to dial GRPC")
	}
	defer conn.Close()
	client := chickaree.NewChickareeDBClient(conn)

	tcpServer := redis.NewTCPServer(fmt.Sprintf(":%d", cfg.Port), client)
	defer tcpServer.Close()

	log.Error().Err(<-tcpServer.Run()).Msg("terminated")
}

type Config struct {
	StorageServer string `yaml:"storage-server"`
	Port          int    `yaml:"port"`
}

func LoadConfiguration() (Config, error) {
	cfgfilePtr := flag.String("config-file", "", "load configurations from a file")
	flag.Parse()

	cfg := Config{
		Port:          6379,
		StorageServer: ":8080",
	}

	if cfgfilePtr != nil && *cfgfilePtr != "" {
		if err := cfg.LoadFromFile(*cfgfilePtr); err != nil {
			return cfg, fmt.Errorf("unable to load configuration from json: %s\n", *cfgfilePtr) // #4
		}
	}

	cfg.LoadFromEnv()
	return cfg, nil
}

func (config *Config) LoadFromFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("unable to load configuration file")
		return errors.New("unable to load configuration")
	}
	if err := yaml.Unmarshal(b, config); err != nil {
		log.Error().Err(err).Str("path", path).Msg("unable to parse configuration file")
		return errors.New("unable to load configuration")
	}
	return nil
}

func (config Config) LoadFromEnv() {
	if val := os.Getenv("STORAGE_SERVER"); val != "" {
		config.StorageServer = val
	}
	if val := os.Getenv("PORT"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			config.Port = v
		}
	}

	return
}
