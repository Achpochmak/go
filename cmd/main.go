package main

import (
	"fmt"

	"HOMEWORK-1/internal/cli"
	"HOMEWORK-1/internal/module"
	"HOMEWORK-1/internal/storage"
)

const (
	fileName = "pvz.json"
)

func main() {
	storageJSON := storage.NewStorage(fileName)
	pvz := module.NewModule(module.Deps{
		Storage: storageJSON,
	})
	commands := cli.NewCLI(cli.Deps{Module: pvz})
	
	for {
		if err := commands.Run(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("done")
	}
}


