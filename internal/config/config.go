package config

import (
	"flag"
)

var (
	AppConfig = NewConfig()
	SecretKey = []byte("M8axEo25vLElQ8n85CvmFRmNrFWt0YQq")
)

type Config struct {
	HttpPort	*int
	PgHost		*string
	PgPort		*int
	PgName		*string
	PgUser		*string
	PgPass		*string
}

func NewConfig() *Config {
	config := &Config{
		HttpPort: flag.Int("port", 8080, "port on which the application will run"),
		PgHost: flag.String("host", "89.46.131.181", "host for PostgreSQL database"),
		PgPort: flag.Int("dbport", 5432, "port for PostgreSQL database"),
		PgName: flag.String("dbname", "stakewolle", "host for PostgreSQL database"),
		PgUser: flag.String("dbuser", "roman", "username for PostgreSQL database"),
		PgPass: flag.String("dbpass", "030905romaN", "password for PostgreSQL database"),
	}

	return config
}