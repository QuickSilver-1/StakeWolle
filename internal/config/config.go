package config

import "flag"

func init() {
	HttpPort := flag.Int("port", 8080, "port on which the application will run")
	PgHost := flag.String("host", "89.46.131.181", "host for PostgreSQL database")
	PgPort := flag.Int("dbport", 5432, "port for PostgreSQL database")
	PgName := flag.String("dbname", "89.46.131.181", "host for PostgreSQL database")
	PgUser := flag.String("dbuser", "roman", "username for PostgreSQL database")
	PgPass := flag.String("dbpass", "030905romaN", "password for PostgreSQL database")
}