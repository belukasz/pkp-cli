/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/belukasz/pkp-cli/app"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"sort"
	"sync"
	"time"
)

type Connections []scrapper.Connection

type DayScrape struct {
	Connections
	Date string
}

func (ds *DayScrape) RealDate() time.Time {
	t, err := time.Parse("2-01-2006_15:04", ds.Date)
	if err != nil {
		panic(err)
	}
	return t
}

type ScrapeData struct {
	From string
	To   string
	Date string
}

func printTable(scrapped []DayScrape) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Date", "Depart", "Arrival", "Price", "Train", "Duration"})
	for _, ds := range scrapped {
		for _, conn := range ds.Connections {
			date := ds.RealDate().Format("Mon Jan02")
			t.AppendRow(table.Row{date, conn.DepartureTime, conn.ArrivalTime, conn.Price, conn.Name, conn.Duration})
		}
		t.AppendSeparator()
	}
	t.Render()
}

func main() {
	count := 30
	outputs := make(chan DayScrape)
	inputs := make(chan ScrapeData)
	now := time.Now()

	maxWorkers := 5

	wg := &sync.WaitGroup{}
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(inputs, outputs, wg)
	}

	producerGroup := &sync.WaitGroup{}
	for i := 0; i < count; i++ {
		producerGroup.Add(1)
		go func() {
			defer producerGroup.Done()
			t := now.AddDate(0, 0, i)
			date := fmt.Sprintf("%d-%d-%d_06:00", t.Day(), t.Month(), t.Year())
			inputs <- ScrapeData{"krakow", "warszawa", date}
		}()
	}

	go workerMonitor(outputs, inputs, wg, producerGroup)
	
	// scrapped := []DayScrape{}
	scrapped := []DayScrape{}
	for dayScrape := range outputs {
		scrapped = append(scrapped, dayScrape)
	}
	sort.Slice(scrapped, func(i, j int) bool {
		return scrapped[j].RealDate().After(scrapped[i].RealDate())
	})
	printTable(scrapped)

}
func workerMonitor(results chan<- DayScrape, inputs chan<- ScrapeData, wg *sync.WaitGroup, pg *sync.WaitGroup) {
	pg.Wait()
	close(inputs)
	wg.Wait()
	close(results)
}

func worker(jobs <-chan ScrapeData, results chan<- DayScrape, wg *sync.WaitGroup) {
	defer wg.Done()
	for scrapeData := range jobs {
		con, _ := scrapper.ScrapeOneDay(scrapeData.From, scrapeData.To, scrapeData.Date)
		results <- DayScrape{con, scrapeData.Date}
	}
}
