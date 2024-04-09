package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/reonardoleis/manager/internal/date"
	"github.com/reonardoleis/manager/internal/models"
	provider "github.com/reonardoleis/manager/internal/providers"
	"github.com/reonardoleis/manager/internal/service"
)

func main() {
	billFile, err := os.Open("bill.json")
	if err != nil {
		panic(err)
	}

	defer billFile.Close()

	mapperFile, err := os.Open("mappings.json")
	if err != nil {
		panic(err)
	}

	defer mapperFile.Close()

	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	billContents, err := io.ReadAll(billFile)
	if err != nil {
		panic(err)
	}

	mapperContents, err := io.ReadAll(mapperFile)
	if err != nil {
		panic(err)
	}

	configContents, err := io.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	mapper := &models.Mapper{}
	bill := &models.Bill{}
	config := &models.Config{}

	err = mapper.FromJSON(string(mapperContents))
	if err != nil {
		panic(err)
	}

	err = bill.FromJSON(string(billContents))
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

	txs := bill.TxsWithDate(window)

	question := "The following transactions will be added to Notion:\n" + bill.GetFormattedTitles(window, mapper) + "\nDo you want to proceed? (y/N): "

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

	notionProvider := provider.NewNotionProvider(mapper, config)
	service := service.New(notionProvider, mapper)

	err = service.Run(txs)
	if err != nil {
		panic(err)
	}
}
