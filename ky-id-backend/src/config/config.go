package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Enviroment  string `yaml:"env"  env-required:"true"  env-default:"local"`
	StoragePath string `yaml:"db_path"  env-required:"false"  env-default:"./database/storage.db"`
	LogPath     string `yaml:"log_path"  env-required:"false"  env-default:"./log/"`
	HTTPServer  `yaml:"http_server"  env-required:"true"`
}

type HTTPServer struct {
	Host        string        `yaml:"host" env-required:"false" env-default:"localhost"`
	Port        string        `yaml:"port" env-required:"false" env-default:"8000"`
	Timeout     time.Duration `yaml:"timeout" env-required:"false" env-default:"2s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-required:"false" env-default:"60s"`
}

type EnviromentVariables struct {
	PermittedUsername string `env:"PERMITTED_USER_NAME" env-required:"true"`
	PermittedPassword string `env:"PERMITTED_USER_PASS" env-required:"true"`
}

// config_path - path to .yaml file (./config/local.yaml)
func MustReadConfig(config_path string) (*Config, error) {

	if config_path == "" {
		return nil, fmt.Errorf("could not resolve config path")
	}

	var cfg Config
	err := cleanenv.ReadConfig(config_path, &cfg)

	if err != nil {
		return nil, fmt.Errorf("could not read config by path %s", config_path)
	}

	fmt.Printf("Config: %v\n", cfg)

	return &cfg, nil
}

// env_path - path to .env file
func MustReadEnv(env_path string) error {

	if env_path == "" {
		return fmt.Errorf("could not resolve env path")
	}

	var env EnviromentVariables

	err := godotenv.Load(env_path)

	if err != nil {
		return fmt.Errorf("could not read enviroment variables file by path %s", env_path)
	}

	err = cleanenv.ReadEnv(&env)

	if err != nil {
		return fmt.Errorf("could not load enviroment variables from file by path %s", env_path)
	}

	fmt.Printf("Env: %v\n", env)

	return nil
}
