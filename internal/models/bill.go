package models

import (
	"encoding/json"
	"fmt"
	"strings"
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

func (b Bill) GetFormattedTitles(dates []string, mapper *Mapper) string {
	titles := []string{}

	txs := b.TxsWithDate(dates)

	for _, tx := range txs {
		titles = append(titles, fmt.Sprintf("- %s (R$ %.2f)", mapper.GetName(tx.Title), float64(tx.Amount)/100))
	}

	return strings.Join(titles, "\n")
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
