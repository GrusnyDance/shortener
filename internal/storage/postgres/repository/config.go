package repository

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"gopkg.in/yaml.v3"
	"net/url"
	"os"
)

// Connection settings
type Config struct {
	Storage struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbname"`
		Timeout  string `yaml:"timeout"`
	} `yaml:"storage"`
}

// Connection str constructor
func NewPoolConfig() (*pgxpool.Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open("./internal/storage/postgres/config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decoder
	decoder := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err = decoder.Decode(&config); err != nil {
		return nil, err
	}

	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%s",
		"postgres",
		url.QueryEscape(config.Storage.Username),
		url.QueryEscape(config.Storage.Password),
		config.Storage.Host,
		config.Storage.Port,
		config.Storage.DbName,
		config.Storage.Timeout)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	return poolConfig, nil
}
