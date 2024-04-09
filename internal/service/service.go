package service

import "github.com/reonardoleis/manager/internal/models"

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
	for _, tx := range txs {
		err := s.provider.Insert(&tx)
		if err != nil {
			return err
		}
	}
	return nil
}
