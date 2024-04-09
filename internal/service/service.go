package service

import (
	"log"
	"sync"

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
	wg := sync.WaitGroup{}
	wg.Add(len(txs))
	for idx, tx := range txs {
		go func(idx int, tx models.Tx) {
			defer wg.Done()

			err := s.provider.Insert(&tx)
			if err != nil {
				log.Println("error inserting tx", tx.Title, idx, err)
				return
			}

			log.Println(idx+1, "of", len(txs), "inserted")
		}(idx, tx)
	}

	wg.Wait()
	return nil
}
