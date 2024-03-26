package go_ipay

type PaymentMethod struct {
	Card *Card
}

type Card struct {
	Name  string
	Token *string
}
