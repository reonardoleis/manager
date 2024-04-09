package provider

import (
	"context"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/jomei/notionapi"
	"github.com/reonardoleis/manager/internal/models"
	"github.com/reonardoleis/manager/internal/service"
)

type NotionProvider struct {
	cli    *notionapi.Client
	mapper *models.Mapper
	config *models.Config
}

func NewNotionProvider(mapper *models.Mapper, config *models.Config) service.Provider {
	cli := notionapi.NewClient(notionapi.Token(config.NotionKey))
	return &NotionProvider{cli, mapper, config}
}

func (n NotionProvider) name(v *models.Tx) notionapi.TitleProperty {
	return notionapi.TitleProperty{
		Title: []notionapi.RichText{
			{
				Text: &notionapi.Text{
					Content: n.mapper.GetName(v.Title),
				},
			},
		},
	}
}

func (n NotionProvider) amount(v *models.Tx) notionapi.NumberProperty {
	return notionapi.NumberProperty{
		Number: -float64(v.Amount) / 100,
	}
}

func (n NotionProvider) status() notionapi.SelectProperty {
	currentMonth := time.Now().Month()

	month, err := n.mapper.GetMonth(int(currentMonth) - 1)
	if err != nil {
		panic(err)
	}

	return notionapi.SelectProperty{
		Select: notionapi.Option{
			Name: month,
		},
	}
}

func (n NotionProvider) category(v *models.Tx) notionapi.SelectProperty {
	category := n.mapper.GetCategory(v.Title, v.Category)
	return notionapi.SelectProperty{
		Select: notionapi.Option{
			Name: category,
		},
	}
}

func (n NotionProvider) method() notionapi.RichTextProperty {
	return notionapi.RichTextProperty{
		RichText: []notionapi.RichText{
			{
				Text: &notionapi.Text{
					Content: "c",
				},
			},
		},
	}
}

func (n NotionProvider) id(v *models.Tx) notionapi.RichTextProperty {
	return notionapi.RichTextProperty{
		RichText: []notionapi.RichText{
			{
				Text: &notionapi.Text{
					Content: v.ID,
				},
			},
		},
	}

}

func (n NotionProvider) exists(v *models.Tx) (bool, error) {
	r, err := n.cli.Database.Query(context.Background(), notionapi.DatabaseID(n.config.NotionDatabaseID), &notionapi.DatabaseQueryRequest{
		Filter: notionapi.PropertyFilter{
			Property: "ID",
			RichText: &notionapi.TextFilterCondition{
				Equals: v.ID,
			},
		},
	})
	if err != nil {
		return false, err
	}

	if len(r.Results) > 0 {
		return true, nil
	}

	return false, nil
}

func (n NotionProvider) Filter(txs []models.Tx) []models.Tx {
	out := []models.Tx{}

	wg := sync.WaitGroup{}
	wg.Add(len(txs))

	lock := sync.Mutex{}

	for idx, tx := range txs {
		go func(idx int, tx models.Tx) {
			defer wg.Done()

			exists, err := n.exists(&tx)
			if err != nil {
				log.Println("error checking if tx exists", tx.Title, idx, err)
				return
			}

			if !exists {
				lock.Lock()
				out = append(out, tx)
				lock.Unlock()
			}
		}(idx, tx)
	}

	wg.Wait()

	sort.Slice(out, func(i, j int) bool {
		return out[i].Idx < out[j].Idx
	})

	return out
}

func (n NotionProvider) Insert(v *models.Tx) error {
	_, err := n.cli.Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent: notionapi.Parent{
			DatabaseID: notionapi.DatabaseID(n.config.NotionDatabaseID),
		},
		Properties: notionapi.Properties{
			"Name":   n.name(v),
			"Valor":  n.amount(v),
			"Status": n.status(),
			"Tipo":   n.category(v),
			"Método": n.method(),
			"ID":     n.id(v),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
