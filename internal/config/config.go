package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

const envFileName = ".env"

type Config struct {
	env *EnvSetting
}

type EnvSetting struct {
	HTTPBindAddr string `env:"HTTP_BIND_ADDR" env-default:":8080" env-description:"HTTP server bind address"`
}

func findConfigFile() bool {
	_, err := os.Stat(envFileName)
	return err == nil
}

func (e *EnvSetting) GetHelpString() (string, error) {
	baseHeader := "options which can be set with env: "
	helpString, err := cleanenv.GetDescription(e, &baseHeader)
	if err != nil {
		return "", fmt.Errorf("failed to get help string: %w", err)
	}
	return helpString, nil
}

func New(logger *zap.Logger) *Config {
	envSetting := &EnvSetting{}

	helpString, err := envSetting.GetHelpString()
	if err != nil {
		logger.Fatal("failed to get help string", zap.Error(err))
	}
	logger.Info(helpString)

	if findConfigFile() {
		if err := cleanenv.ReadConfig(envFileName, envSetting); err != nil {
			logger.Fatal("failed to read env config", zap.Error(err))
		}
	} else if err := cleanenv.ReadEnv(envSetting); err != nil {
		logger.Fatal("failed to read env config", zap.Error(err))
	}

	return &Config{env: envSetting}
}

func (c *Config) GetHTTPBindAddr() string {
	return c.env.HTTPBindAddr
}
