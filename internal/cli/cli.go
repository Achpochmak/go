package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"HOMEWORK-1/internal/models"
)

// NewCLI creates a command line interface
func NewCLI(d Deps) *CLI {
	cli := &CLI{
		Deps: d,
		commandList: []command{
			{
				name:        help,
				description: helpDescription,
			},
			{
				name:        addOrder,
				description: addOrderDescription,
			},
			{
				name:        deleteOrder,
				description: deleteOrderDescription,
			},
			{
				name:        deliverOrder,
				description: deliverOrderDescription,
			},
			{
				name:        GetOrderByID,
				description: GetOrderByIDDescription,
			},
			{
				name:        listOrder,
				description: listOrderDescription,
			},
			{
				name:        refund,
				description: refundDescription,
			},
			{
				name:        listRefund,
				description: listRefundDescription,
			},
			{
				name:        setWorkers,
				description: setWorkersDescription,
			},
		},
		taskQueue:     make(chan task, 10),
		numWorkers:    2,
		workerPool:    make(chan struct{}, 2),
		orderLocks:    make(map[models.ID]*sync.Mutex),
		notifications: make(chan string, 10),
	}
	go cli.notificationHandler()
	return cli
}

// Run ..
func (c *CLI) Run() error {
	ctx := context.Background()

	for i := 0; i < c.numWorkers; i++ {
		c.wg.Add(1)
		go c.worker(ctx)
	}

	c.handleSignals()

	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		args := strings.Fields(strings.TrimSpace(input))
		if len(args) == 0 {
			fmt.Println("command isn't set")
			continue
		}

		commandName := args[0]
		if commandName == exit {
			close(c.taskQueue)
			break
		}
		c.taskQueue <- task{commandName: commandName, args: args[1:]}
	}

	c.wg.Wait()
	fmt.Println("All tasks completed. Exiting...")
	os.Exit(0)
	return nil
}
