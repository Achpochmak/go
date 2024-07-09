package cli

import "fmt"

func (c *CLI) SwitchOutputMode() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.outputKafka = !c.outputKafka
	c.outbox.getMessageFromRepoToOutbox() //выводим команду переключения из аутбокса, чтобы не было лишнего переключения
	if c.outputKafka {
		fmt.Println("Output mode switched to Kafka")
	} else {
		fmt.Println("Output mode switched to console")
	}
}
