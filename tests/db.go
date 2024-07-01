
package tests

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetConfigName("test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/OZON/homework-1/tests") 
	
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("ошибка файла конфигурации: %v", err)
	}
}

func ConnectDB() *pgxpool.Pool {
	dbPassword := viper.GetString("database.password")
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetInt("database.port")
	dbUser := viper.GetString("database.user")
	dbName := viper.GetString("database.dbname")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("ошибка url: %v", err)
	}
	config.MaxConns = 10

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("нет соединения: %v", err)
	}
	return pool
}
