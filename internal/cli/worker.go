package cli

import (
	"context"
	"fmt"
	"os"
)

func (c *CLI) worker(ctx context.Context) {
	defer c.wg.Done()
	for t := range c.taskQueue {
		var err error
		startMsg := fmt.Sprintf("Началась обработка команды: %s", t.commandName)
		endMsg := fmt.Sprintf("Завершилась обработка команды: %s", t.commandName)
		c.notifications <- startMsg
		switch t.commandName {

		case help:
			c.help()
		case addOrder:
			err = c.addOrder(ctx, t.args)
		case deleteOrder:
			err = c.deleteOrder(ctx, t.args)
		case deliverOrder:
			err = c.deliverOrder(ctx, t.args)
		case listOrder:
			err = c.listOrder(ctx)
		case GetOrderByID:
			err = c.GetOrderByID(ctx, t.args)
		case getOrdersByCustomer:
			err = c.getOrdersByCustomer(ctx, t.args)
		case refund:
			err = c.refund(ctx, t.args)
		case listRefund:
			err = c.listRefund(ctx, t.args)
		case setWorkers:
			err = c.setWorkers(t.args)
		case exit:
			fmt.Println("Exiting...")
			c.mu.Unlock()
			close(c.taskQueue)
			os.Exit(0)
		default:
			fmt.Println("command isn't set")
		}
		if err != nil {
			fmt.Println("Ошибка:", err)
		}
		c.notifications <- endMsg
	}
}
