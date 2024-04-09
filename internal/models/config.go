package models

import "encoding/json"

type Config struct {
	NotionKey        string `json:"notionKey"`
	NotionDatabaseID string `json:"notionDatabaseId"`
	Cpf              string `json:"cpf"`
	Password         string `json:"password"`
	ClientSecret     string `json:"client_id"`
	CreditCardURL    string `json:"creditCardUrl"`
}

func (c *Config) FromJSON(v string) error {
	err := json.Unmarshal([]byte(v), &c)
	if err != nil {
		return err
	}
	return nil
}
