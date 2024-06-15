package cli

import (
	"HOMEWORK-1/internal/models/customErrors"
	"context"
	"flag"
	"fmt"
)

// Измeнение числа рутин
func (c *CLI) setWorkers(args []string) error {
	num, err := c.parseSetWorkers(args)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if num > c.numWorkers {
		for i := c.numWorkers; i < num; i++ {
			c.wg.Add(1)
			go c.worker(context.Background())
		}
	} else if num < c.numWorkers {
		for i := num; i < c.numWorkers; i++ {
			c.taskQueue <- task{commandName: "exit"}
		}
	}

	c.numWorkers = num
	fmt.Printf("Число рутин %d\n", c.numWorkers)
	return nil
}

// Парсинг параметров изменения количества рутин
func (c *CLI) parseSetWorkers(args []string) (int, error) {
	var num int
	fs := flag.NewFlagSet("setWorkers", flag.ContinueOnError)
	fs.IntVar(&num, "num", c.numWorkers, "use --num=1")

	if err := fs.Parse(args); err != nil {
		return 0, err
	}

	if num < 1 {
		return 0, customErrors.ErrWorkersLessThanOne
	}
	return num, nil
}
