package scrapper

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

func ScrapeConnections(days int, starthour string, trainType string, startStation string, arrivalStation string) []DayScrape {

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
	for i := 0; i < days; i++ {
		producerGroup.Add(1)
		go func() {
			defer producerGroup.Done()
			t := now.AddDate(0, 0, i)
			date := fmt.Sprintf("%d-%d-%d_%s", t.Day(), t.Month(), t.Year(), starthour)
			inputs <- ScrapeData{From: startStation, To: arrivalStation, Date: date}
		}()
	}
	go workerMonitor(outputs, inputs, wg, producerGroup)
	scrapped := []DayScrape{}
	for dayScrape := range outputs {

		scrapped = append(scrapped, dayScrape)
	}
	sort.Slice(scrapped, func(i, j int) bool {
		return scrapped[j].RealDate().After(scrapped[i].RealDate())
	})
	return scrapped
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
		con, _ := ScrapeOneDay(scrapeData.From, scrapeData.To, scrapeData.Date, "EIP")
		results <- DayScrape{Connections: con, Date: scrapeData.Date}
	}
}
