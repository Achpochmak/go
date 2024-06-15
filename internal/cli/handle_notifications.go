package cli

import "fmt"

// Обработка уведомлений
func (c *CLI) notificationHandler() {
	for msg := range c.notifications {
		fmt.Println(msg)
	}
}
