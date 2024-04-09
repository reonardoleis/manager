package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/reonardoleis/manager/internal/bank"
	"github.com/reonardoleis/manager/internal/date"
	"github.com/reonardoleis/manager/internal/models"
	provider "github.com/reonardoleis/manager/internal/providers"
	"github.com/reonardoleis/manager/internal/service"
)

func main() {
	godotenv.Overload(".env")

	mapperFile, err := os.Open(os.Getenv("CONFIG_PATH") + "/mappings.json")
	if err != nil {
		panic(err)
	}

	defer mapperFile.Close()

	configFile, err := os.Open(os.Getenv("CONFIG_PATH") + "/config.json")
	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	mapperContents, err := io.ReadAll(mapperFile)
	if err != nil {
		panic(err)
	}

	configContents, err := io.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	mapper := &models.Mapper{}

	config := &models.Config{}

	err = mapper.FromJSON(string(mapperContents))
	if err != nil {
		panic(err)
	}

	err = config.FromJSON(string(configContents))
	if err != nil {
		panic(err)
	}

	daysFlag := flag.Int("days", 0, "Get transactions from n days ago")

	flag.Parse()

	window := date.Days(*daysFlag)

	bank := bank.NewNubank(config.Cpf, config.Password, config.ClientSecret)
	defer bank.Revoke()

	err = bank.Authorize()
	if err != nil {
		panic(err)
	}

	provider := provider.NewNotionProvider(mapper, config)
	service := service.New(provider, bank, mapper)

	bill, err := service.GetBill(config.CreditCardURL)
	if err != nil {
		panic(err)
	}

	txs := provider.Filter(bill.TxsWithDate(window))
	question :=
		"The following transactions will be added to Notion:\n" +
			bill.GetFormattedTitles(window, mapper) +
			"\nDo you want to proceed? (y/N): "

	if len(txs) == 0 {
		log.Println("No transactions found")
		return
	}

	fmt.Print(question)

	var answer string
	fmt.Scanln(&answer)

	if answer != "y" {
		return
	}

	err = service.Run(txs)
	if err != nil {
		panic(err)
	}
}
