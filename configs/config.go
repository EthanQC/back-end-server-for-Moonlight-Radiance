// config/config.go
package config

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
		Port string
	}
}
