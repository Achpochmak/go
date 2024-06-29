package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// NewCLI creates a command line interface
func NewCLI(d Deps) *CLI {
	cli := &CLI{
		Deps:    d,
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
		notifications: make(chan string, 10),
		taskQueueOpen: true,
		wg:            sync.WaitGroup{},
	}
	
	return cli
}

// Run ..
func (c *CLI) Run() error {
	go c.notificationHandler()
	defer close(c.notifications)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < c.numWorkers; i++ {
		c.wg.Add(1)
		go c.worker(ctx)
	}

	c.handleSignals(cancel)

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
			c.mu.Lock()
			if c.taskQueueOpen {
				c.taskQueueOpen = false
				close(c.taskQueue)
			}
			c.mu.Unlock()

			go func() {
				time.Sleep(5 * time.Second)
				cancel()
			}()
			break
		}
		if c.taskQueueOpen {
			c.taskQueue <- task{commandName: commandName, args: args[1:]}
		} else {
			fmt.Println("Доступ закрыт")
		}

	}

	c.wg.Wait()
	fmt.Println("Все задачи завершены.")
	os.Exit(0)
	return nil
}
