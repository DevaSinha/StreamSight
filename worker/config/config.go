package config

type Config struct {
	ApiEndpoint string
}

func LoadConfig() Config {
	return Config{
		ApiEndpoint: "http://host.docker.internal:8080",
	}
}
