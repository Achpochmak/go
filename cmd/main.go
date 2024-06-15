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
    viper.SetConfigName("docker-compose")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("ошибка файла конфигурации")
    }
}

func connectDB() *pgxpool.Pool {
	dbURL := viper.GetString("services.app.environment.DATABASE_URL")
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("ошибка url")
	}
	config.MaxConns = 10

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("нет соединения")
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
	})
	commands := cli.NewCLI(cli.Deps{Module: pvz})

	for {
		if err := commands.Run(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("done")
	}

}
