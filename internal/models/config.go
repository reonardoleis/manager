package models

import "encoding/json"

type Config struct {
	NotionKey        string `json:"notionKey"`
	NotionDatabaseID string `json:"notionDatabaseId"`
}

func (c *Config) FromJSON(v string) error {
	err := json.Unmarshal([]byte(v), &c)
	if err != nil {
		return err
	}
	return nil
}
