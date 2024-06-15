package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Обработка сигналов
func (c *CLI) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		fmt.Printf("Получена команда %s. Exiting...\n", sig)
		close(c.taskQueue)
		c.wg.Wait()
		os.Exit(0)
	}()
}
