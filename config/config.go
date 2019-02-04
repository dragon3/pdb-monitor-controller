package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	LogLevel       string `envconfig:"LOG_LEVEL" default:"INFO"`
	Env            string `envconfig:"ENV" required:"true" default:"development"`
	KubeConfigPath string `envconfig:"KUBE_CONFIG_PATH"`
}

func NewFromEnv() (*Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return nil, errors.Wrap(err, "failed to process envconfig")
	}
	return &config, nil
}
