package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

type Config struct {
	Database               DatabaseConfig         `yaml:"database"`
	Kafka                  KafkaConfig            `yaml:"kafka"`
	Server                 ServerConfig           `yaml:"server"`
	ProfileServiceSettings ProfileServiceSettings `yaml:"profileServiceSettings"`
}

type DatabaseConfig struct {
	Shards      []DatabaseShardConfig `yaml:"shards"`
	BucketCount int                   `yaml:"bucket_count"`
}

type DatabaseShardConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	SSLMode  string `yaml:"ssl_mode"`
}

type KafkaConfig struct {
	Host                    string `yaml:"host"`
	Port                    int    `yaml:"port"`
	MenuGenerationTopicName string `yaml:"menu_generation_topic_name"`
}

type ServerConfig struct {
	GRPCPort int `yaml:"grpc_port"`
	HTTPPort int `yaml:"http_port"`
}

type ProfileServiceSettings struct {
	MinUsernameLen int `yaml:"minUsernameLen"`
	MaxUsernameLen int `yaml:"maxUsernameLen"`
	MinPasswordLen int `yaml:"minPasswordLen"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return &config, nil
}
