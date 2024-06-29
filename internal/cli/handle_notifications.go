package cli

import "fmt"

// Обработка уведомлений
func (c *CLI) notificationHandler() {
	for msg := range c.notifications {
		fmt.Println(msg)
	}
}

// Отправка уведомления старта работы 
func (c *CLI) sendStartNotification(t task) {
	startMsg := fmt.Sprintf("Началась обработка команды: %s", t.commandName)
	c.notifications <- startMsg
}

// Обработка уведомления окончания работы
func (c *CLI) sendEndNotification(t task) {
	endMsg := fmt.Sprintf("Завершилась обработка команды: %s", t.commandName)
	c.notifications <- endMsg
}