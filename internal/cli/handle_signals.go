package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Обработка сигналов
func (c *CLI) handleSignals(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		fmt.Printf("Получена команда %s. Exiting...\n", sig)
		if c.taskQueueOpen {
			close(c.taskQueue)
			c.taskQueueOpen = false
		}
		go func() {
			time.Sleep(5 * time.Second)
			cancel()
		}()
		c.wg.Wait()
		os.Exit(0)
	}()
}
