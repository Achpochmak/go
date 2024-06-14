package main

import (
	"context"
	"fmt"

	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/internal/module"
	"HOMEWORK-1/internal/repository"
	"HOMEWORK-1/internal/repository/transactor"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {

	dbURL := "postgres://postgres:password@localhost:5432/oms"
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		fmt.Println("ошибка url")
	}

	config.MaxConns = 10

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		fmt.Println("нет соединения")
	}
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
