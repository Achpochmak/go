package main

import (
	"context"
	"fmt"
	"log"

	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/internal/module"
	"HOMEWORK-1/internal/repository"
	"HOMEWORK-1/internal/repository/transactor"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("ошибка файла конфигурации: %v", err)
	}
}

func connectDB() *pgxpool.Pool {

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
	config.MaxConns = 30

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("нет соединения: %v", err)
	}
	return pool
}

func main() {
	initConfig()
	pool := connectDB()
	defer pool.Close()
	tm := &transactor.TransactionManager{Pool: pool}
	repo := repository.NewRepository(tm)
	pvz := module.NewModule(module.Deps{
		Repository: repo,
		Transactor: tm,
	})
	commands := cli.NewCLI(cli.Deps{Module: pvz})

	for {
		if err := commands.Run(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("done")
	}
}
