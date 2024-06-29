package cli

import "context"

type CLIHandler interface {
	AddOrder(ctx context.Context, args []string) error
	DeliverOrder(ctx context.Context, args []string) error
	Refund(ctx context.Context, args []string) error
	GetOrderByID(ctx context.Context, args []string) error
	DeleteOrder(ctx context.Context, args []string) error
	GetOrdersByCustomer(ctx context.Context, args []string) error 
	ListOrder(ctx context.Context) error 
	ListRefund(ctx context.Context, args []string) error 
	SetWorkers(args []string) error 
}

func (c *CLI) SetHandler(handler CLIHandler) {
	c.handler = handler
}

type CLIHandlerImpl struct {
	cli *CLI
}

func (h *CLIHandlerImpl) AddOrder(ctx context.Context, args []string) error {
	return h.cli.addOrder(ctx, args)
}

func (h *CLIHandlerImpl) DeliverOrder(ctx context.Context, args []string) error {
	return h.cli.deliverOrder(ctx, args)
}

func (h *CLIHandlerImpl) Refund(ctx context.Context, args []string) error {
	return h.cli.refund(ctx, args)
}

func (h *CLIHandlerImpl) GetOrderByID(ctx context.Context, args []string) error {
	return h.cli.getOrderByID(ctx, args)
}

func (h *CLIHandlerImpl) GetOrdersByCustomer(ctx context.Context, args []string) error {
	return h.cli.getOrdersByCustomer(ctx, args)
}

func (h *CLIHandlerImpl) DeleteOrder(ctx context.Context, args []string) error {
	return h.cli.deleteOrder(ctx, args)
}

func (h *CLIHandlerImpl) ListOrder(ctx context.Context) error {
	return h.cli.listOrder(ctx)
}

func (h *CLIHandlerImpl) ListRefund(ctx context.Context, args []string) error {
	return h.cli.listRefund(ctx, args)
}

func (h *CLIHandlerImpl) SetWorkers(args []string) error {
	return h.cli.setWorkers(args)
}

func NewCLIHandler(cli *CLI) CLIHandler {
	return &CLIHandlerImpl{cli: cli}
}
