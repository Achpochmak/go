package cli

import (
	"context"
	"fmt"
)

func (c *CLI) worker(ctx context.Context) {
	defer c.wg.Done()
	for {
		select {
		case t, ok := <-c.taskQueue:
			if !ok {
				return
			}
			var err error
			c.sendStartNotification(t)
			switch t.commandName {
			case help:
				c.help()
			case addOrder:
				err = c.AddOrder(ctx, t.args)
			case deleteOrder:
				err = c.DeleteOrder(ctx, t.args)
			case deliverOrder:
				err = c.DeliverOrder(ctx, t.args)
			case listOrder:
				err = c.ListOrder(ctx)
			case GetOrderByID:
				err = c.GetOrderByID(ctx, t.args)
			case getOrdersByCustomer:
				err = c.GetOrdersByCustomer(ctx, t.args)
			case refund:
				err = c.Refund(ctx, t.args)
			case listRefund:
				err = c.ListRefund(ctx, t.args)
			case setWorkers:
				err = c.setWorkers(t.args)
			case exit:
				fmt.Println("Exiting...")
				c.mu.Lock()
				close(c.taskQueue)
				c.mu.Unlock()
				return
			default:
				fmt.Println("command isn't set")
			}
			if err != nil {
				fmt.Println("Ошибка:", err)
			}
			c.sendEndNotification(t)
		case <-ctx.Done():
			return
		}
	}
}
