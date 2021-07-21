package storage

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

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
	SefPort int `yaml:"serf-port"`
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

	hostNameSplit := strings.Split(hostname, "-")
	defaultBootstrap := hostNameSplit[len(hostNameSplit)-1] == "0"

	cfg := ServerConfig{
		Config: Config{
			StoragePath: "chicakree.db",
			RaftDir:     "/tmp",
		},
		NodeName:  hostname,
		RPCPort:   8400,
		SefPort:   8401,
		Bootstrap: defaultBootstrap,
	}

	if cfgfilePtr != nil && *cfgfilePtr != "" { // #3
		if err := cfg.LoadFromFile(*cfgfilePtr); err != nil {
			return cfg, fmt.Errorf("unable to load configuration from json: %s\n", *cfgfilePtr) // #4
		}
	}

	cfg.LoadFromEnv()

	if strings.Contains(cfg.BindAddr, "$HOSTNAME") {
		cfg.BindAddr = strings.Replace(cfg.BindAddr, "$HOSTNAME", hostname, 1)
	}

	if cfg.Bootstrap {
		cfg.StartJoinAddrs = nil
	}
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

func (config *ServerConfig) LoadFromEnv() {
	if val := os.Getenv("DATA_DIR"); val != "" {
		log.Info().Str("data-dir", val).Msg("update config from env")
		config.DataDir = val
	}
	if val := os.Getenv("RAFT_DIR"); val != "" {
		log.Info().Str("raft-dir", val).Msg("update config from env")
		config.Config.RaftDir = val
	}
	if val := os.Getenv("STORAGE_PATH"); val != "" {
		log.Info().Str("storage-path", val).Msg("update config from env")
		config.Config.StoragePath = val
	}
	if val := os.Getenv("NODE_NAME"); val != "" {
		log.Info().Str("node-name", val).Msg("update config from env")
		config.NodeName = val
	}
	if val := os.Getenv("BIND_ADDR"); val != "" {
		log.Info().Str("bind-addr", val).Msg("update config from env")
		config.BindAddr = val
	}
	if val := os.Getenv("RPC_PORT"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			log.Info().Str("rpc-port", val).Msg("update config from env")
			config.RPCPort = v
		}
	}
	if val := os.Getenv("BOOTSTRAP"); val != "" {
		log.Info().Str("bootstrap", val).Msg("update config from env")
		config.Bootstrap = (val == "true" || val == "1")
	}
	if val := os.Getenv("START_JOIN_ADDRS"); val != "" {
		log.Info().Str("start-join-addrs", val).Msg("update config from env")
		config.StartJoinAddrs = []string{val}
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
