package mongodbexporter

import (
	"errors"

	"go.opentelemetry.io/collector/confmap"
)

// Config defines configuration for mongodb exporter.
type Config struct {

	// MongoDB connection uri.
	ConnectionURI    string `mapstructure:"conn_uri"`
	CollectionLogs   string `mapstructure:"logs_collection"`
	CollectionTraces string `mapstructure:"traces_collection"`
	Database         string `mapstructure:"db"`
}

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if cfg.ConnectionURI == "" {
		return errors.New("connection uri must be non-empty")
	}
	return nil
}

// Unmarshal a confmap.Conf into the config struct.
func (cfg *Config) Unmarshal(componentParser *confmap.Conf) error {
	if componentParser == nil {
		return errors.New("empty config for file exporter")
	}
	// first load the config normally
	err := componentParser.Unmarshal(cfg, confmap.WithErrorUnused())
	if err != nil {
		return err
	}
	return nil
}
