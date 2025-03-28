package configs

import (
	"os"
	"strconv"
)

type Config struct {
	Database struct {
		DSN string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	JWT struct {
		Secret string
	}
	Server struct {
		Port     string
		Mode     string
		LogLevel string
	}
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	var cfg Config

	// 数据库配置
	cfg.Database.DSN = os.Getenv("DB_DSN")

	// Redis配置
	cfg.Redis.Addr = os.Getenv("REDIS_ADDR")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Redis.DB, _ = strconv.Atoi(os.Getenv("REDIS_DB"))

	// JWT配置
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	// 服务器配置
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	// cfg.Server.Mode = os.Getenv("GIN_MODE")
	// if cfg.Server.Mode == "" {
	// 	cfg.Server.Mode = "release"
	// }

	return &cfg
}
