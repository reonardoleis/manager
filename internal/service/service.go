package service

import (
	"log"

	"github.com/reonardoleis/manager/internal/models"
)

type Provider interface {
	Insert(v *models.Tx) error
}

type Service struct {
	provider Provider
	mapper   *models.Mapper
}

func New(provider Provider, mapper *models.Mapper) Service {
	return Service{
		provider: provider,
		mapper:   mapper,
	}
}

func (s Service) Run(txs []models.Tx) error {
	log.Println("Inserting", len(txs), "transactions")
	for idx, tx := range txs {
		err := s.provider.Insert(&tx)
		if err != nil {
			return err
		}

		log.Println(idx+1, "of", len(txs), "inserted")
	}
	return nil
}
