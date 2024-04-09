package models

import (
	"encoding/json"
)

type Bill struct {
	Txs []Tx `json:"line_items"`
}

func (b Bill) TxsWithDate(dates []string) []Tx {
	txs := []Tx{}

	ondate := func(date string) bool {
		for _, d := range dates {
			if d == date {
				return true
			}
		}
		return false
	}

	for _, tx := range b.Txs {
		if ondate(tx.PostDate) {
			txs = append(txs, tx)
		}
	}

	return txs
}

func (b *Bill) FromJSON(v string) error {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(v), &m)
	if err != nil {
		return err
	}

	bill := m["bill"].(map[string]interface{})

	remarshalled, err := json.Marshal(bill["line_items"])
	if err != nil {
		return err
	}

	txs := []Tx{}

	err = json.Unmarshal(remarshalled, &txs)
	if err != nil {
		return err
	}

	b.Txs = txs

	return nil
}
