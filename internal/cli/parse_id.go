package cli

import (
	"flag"

	"HOMEWORK-1/internal/models/customErrors"
)

// Парсинг ID заказа
func (c *CLI) parseID(args []string) (int, error) {
	var id int
	fs := flag.NewFlagSet(deleteOrder, flag.ContinueOnError)
	fs.IntVar(&id, "id", 0, "use --id=1")

	if err := fs.Parse(args); err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, customErrors.ErrIDNotFound
	}
	return id, nil
}
