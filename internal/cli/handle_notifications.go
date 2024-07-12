package cli

import (
	"HOMEWORK-1/internal/infrastructure/app/sender"
	"fmt"
	"time"
)

const (
	started = "started"
	done    = "done"
)

// Обработка уведомлений
func (c *CLI) notificationHandler() {
	for msg := range c.notifications {
		fmt.Println(msg)
	}
}

// Отправка уведомления старта работы
func (c *CLI) sendStartNotification(t task) {
	c.SendToOutbox(t, started)
	if !c.outputKafka {
		startMsg := fmt.Sprintf("Началась обработка команды: %s", t.commandName)
		c.notifications <- startMsg
	}
}

// Обработка уведомления окончания работы
func (c *CLI) sendEndNotification(t task) {
	c.SendToOutbox(t, done)
	if !c.outputKafka {
		endMsg := fmt.Sprintf("Завершилась обработка команды: %s", t.commandName)
		c.notifications <- endMsg
	}
}

// Запись сообщения в outbox
func (c *CLI) SendToOutbox(t task, status string) {
	c.AnswerID++
	answerID := c.AnswerID

	msg := sender.Message{
		Command:     t.commandName,
		Args:        t.args,
		AnswerID:    int(answerID),
		CreatedAt:   time.Now(),
		Status:      status,
		Success:     false,
		IsAquired:   false,
		IsProcessed: false,
	}

	c.outbox.CreateMessage(&msg)
}
