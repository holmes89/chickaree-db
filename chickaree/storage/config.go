package storage

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/hashicorp/raft"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	StoragePath string `yaml:"storage-path"`
	RaftDir     string `yaml:"raft-dir"`
	Raft        struct {
		raft.Config
		BindAddr    string
		StreamLayer *StreamLayer
		Bootstrap   bool
	}
}

type ServerConfig struct {
	Config          `yaml:"config"`
	ServerTLSConfig *tls.Config
	PeerTLSConfig   *tls.Config
	// DataDir stores the log and raft data.
	DataDir string `yaml:"data-dir"`
	// BindAddr is the address serf runs on.
	BindAddr string `yaml:"bind-addr"`
	// RPCPort is the port for client (and Raft) connections.
	RPCPort int `yaml:"rpc-port"`
	// Raft server id.
	NodeName string `yaml:"node-name"`
	// Bootstrap should be set to true when starting the first node of the cluster.
	StartJoinAddrs []string `yaml:"start-join-addrs"`
	Bootstrap      bool     `yaml:"bootstrap"`
}

func LoadConfiguration() (ServerConfig, error) {
	cfgfilePtr := flag.String("config-file", "", "load configurations from a file")
	flag.Parse() // #2

	hostname, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msg("unable to find hostname")
	}
	cfg := ServerConfig{
		Config: Config{
			StoragePath: "chicakree.db",
			RaftDir:     "/tmp",
		},
		NodeName: hostname,
		BindAddr: "127.0.0.1:8401",
		RPCPort:  8400,
	}

	if cfgfilePtr != nil && *cfgfilePtr != "" { // #3
		if err := cfg.LoadFromFile(*cfgfilePtr); err != nil {
			return cfg, fmt.Errorf("unable to load configuration from json: %s\n", *cfgfilePtr) // #4
		}
	}

	cfg.LoadFromEnv()

	cfg.Raft.LocalID = raft.ServerID(cfg.NodeName)
	return cfg, nil
}

func (config *ServerConfig) LoadFromFile(path string) error {
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

func (config ServerConfig) LoadFromEnv() {
	if val := os.Getenv("DATA_DIR"); val != "" {
		config.DataDir = val
	}
	if val := os.Getenv("STORAGE_PATH"); val != "" {
		config.StoragePath = val
	}
	if val := os.Getenv("NODE_NAME"); val != "" {
		config.NodeName = val
	}
	if val := os.Getenv("BIND_ADDR"); val != "" {
		config.BindAddr = val
	}
	if val := os.Getenv("RPC_PORT"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			config.RPCPort = v
		}
	}
	if val := os.Getenv("BOOTSTRAP"); val != "" {
		config.Bootstrap = (val == "true" || val == "1")
	}
	return
}

func (c ServerConfig) RPCAddr() (string, error) {
	host, _, err := net.SplitHostPort(c.BindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, c.RPCPort), nil
}
