package cli

import "fmt"

func (c *CLI) switchOutputMode() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.outputKafka = !c.outputKafka
	if c.outputKafka {
		fmt.Println("Output mode switched to Kafka")
	} else {
		fmt.Println("Output mode switched to console")
	}
}
