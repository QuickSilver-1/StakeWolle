package config

import (
	"flag"
)

var (
    AppConfig = NewConfig()
    SecretKey = []byte("M8axEo25vLElQ8n85CvmFRmNrFWt0YQq")
)

type Config struct {
    HttpPort    *int
    PgHost      *string
    PgPort      *int
    PgName      *string
    PgUser      *string
    PgPass      *string
}

// NewConfig создает и возвращает новый конфигурационный объект
func NewConfig() *Config {
    config := &Config{
        HttpPort: flag.Int("port", 8080, "port on which the application will run"),           // Порт для запуска приложения
        PgHost: flag.String("host", "89.46.131.181", "host for PostgreSQL database"),          // Хост для базы данных PostgreSQL
        PgPort: flag.Int("dbport", 5432, "port for PostgreSQL database"),                      // Порт для базы данных PostgreSQL
        PgName: flag.String("dbname", "stakewolle", "name of the PostgreSQL database"),       // Имя базы данных PostgreSQL
        PgUser: flag.String("dbuser", "roman", "username for PostgreSQL database"),           // Имя пользователя для базы данных PostgreSQL
        PgPass: flag.String("dbpass", "030905romaN", "password for PostgreSQL database"),     // Пароль для базы данных PostgreSQL
    }

    return config
}
