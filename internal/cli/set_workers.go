package cli

import (
	"context"
	"flag"
	"fmt"
	
	"HOMEWORK-1/internal/models/customErrors"

	"github.com/pkg/errors"
)

// Измeнение числа рутин
func (c *CLI) setWorkers(args []string) error {
	num, err := c.parseSetWorkers(args)
	if err != nil {
		return errors.Wrap(err, "некорректный ввод")
	}

	if num > c.numWorkers {
		for i := c.numWorkers; i < num; i++ {
			c.wg.Add(1)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go c.worker(ctx)
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
