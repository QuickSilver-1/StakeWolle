package config

import "flag"

func init() {
	httpPort := flag.Int("port", 8080, "port on which the application will run")
	pgHost := flag.String("host", "89.46.131.181", "host for PostgreSQL database")
	pgName := flag.String("dbname", "89.46.131.181", "host for PostgreSQL database")
	pgUser := flag.String("dbuser", "roman", "username for PostgreSQL database")
	pgPass := flag.String("dbpass", "030905romaN", "password for PostgreSQL database")
}