package entity

type Transfer struct {
	Payer  int     `json:"payer"`
	Payee  int     `json:"payee"`
	Amount float64 `json:"value"`
}

type Transaction struct {
	Type             string  `json:"type"` // credit or debit
	Amount           float64 `json:"amount"`
	RelatedAccountID int     `json:"related_account_id"`
}
