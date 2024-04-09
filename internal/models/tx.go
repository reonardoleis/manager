package models

type Tx struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Charges  int    `json:"charges"`
	Amount   int    `json:"amount"`
	PostDate string `json:"post_date"`
}
