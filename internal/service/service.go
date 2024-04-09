package service

import (
	"log"
	"sync"

	"github.com/reonardoleis/manager/internal/models"
)

type Provider interface {
	Insert(v *models.Tx) error
}

type Bank interface {
	GetBill(url string) (*models.Bill, error)
	Authorize() error
	Revoke() error
}

type Service struct {
	provider Provider
	bank     Bank
	mapper   *models.Mapper
}

func New(provider Provider, bank Bank, mapper *models.Mapper) Service {
	return Service{
		provider: provider,
		bank:     bank,
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

func (s Service) GetBill(url string) (*models.Bill, error) {
	return s.bank.GetBill(url)
}
