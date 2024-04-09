package models

import (
	"encoding/json"
	"errors"
	"strings"
)

type Place struct {
	Title    string `json:"title"`
	Category string `json:"category"`
}

type Mapper struct {
	Places           map[string]Place  `json:"places"`
	Category         map[string]string `json:"categories"`
	FallbackCategory string            `json:"fallbackCategory"`
	Months           []string          `json:"months"`
}

func (m *Mapper) FromJSON(v string) error {
	err := json.Unmarshal([]byte(v), &m)
	if err != nil {
		return err
	}
	return nil
}

func (m Mapper) GetCategory(place, appCategory string) string {
	for k, v := range m.Places {
		if strings.Contains(strings.ToLower(place), strings.ToLower(k)) {
			return v.Category
		}
	}

	c, ok := m.Category[appCategory]
	if ok {
		return c
	}

	return m.FallbackCategory
}

func (m Mapper) GetName(place string) string {
	for k, v := range m.Places {
		if strings.Contains(strings.ToLower(place), strings.ToLower(k)) {
			return v.Title
		}
	}

	return place
}

func (m Mapper) GetMonth(month int) (string, error) {
	if month < 0 || month > 11 {
		return "", errors.New("month out of range")
	}
	return m.Months[month], nil
}
