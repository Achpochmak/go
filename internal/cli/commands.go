package cli

const (
	help          = "help"
	addOrder    = "add"
	deleteOrder = "delete"
	deliverOrder = "deliver"
	listOrder   = "list"
	OrdersByCustomer   = "customer"
	findOrder   = "find"
	Refund = "refund"
	listRefund = "listrefund"
	exit ="exit"
)

type command struct {
	name        string
	description string
}
