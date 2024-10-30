package scrapper

import (
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"strconv"
	"strings"
	"time"
)

var trainTypes = map[string]string{
	"IC":    "IC--EIP-IC",
	"EIP":   "EIP--EIP-IC",
	"EIPIC": "EIP-IC--EIP-IC",
	"ALL":   "all",
}

type Connections []Connection

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

type Connection struct {
	DepartureTime string
	ArrivalTime   string
	Price         float64
	Name          string
	Duration      string
}

func ScrapeOneDay(start string, end string, date string, trainType string) ([]Connection, error) {
	suffix := trainTypes[trainType]
	url := fmt.Sprintf("https://koleo.pl/rozklad-pkp/%s/%s/%s/all/%s", start, end, date, suffix)
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancelTimeout := context.WithTimeout(ctx, 15*time.Second)
	defer cancelTimeout()

	var renderedContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			renderedContent, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	)
	document, err := goquery.NewDocumentFromReader(strings.NewReader(renderedContent))
	if err != nil {
		return nil, err
	}

	var connections []Connection
	document.Find(".has-train-nr").Each(func(i int, s *goquery.Selection) {
		var (
			from     string
			to       string
			train    string
			brand    string
			price    float64
			duration string
		)
		fromQuery := ".medium-3.small-2.columns.time.from"
		timeToQuery := ".large-3.large-offset-1.medium-3.medium-offset-0.small-2.columns.time.to"
		trainNumberQuery := ".train-number-details"

		from = strings.TrimSpace(s.Find(fromQuery).First().Text())
		to = strings.TrimSpace(s.Find(timeToQuery).First().Text())

		train = strings.TrimSpace(s.Find(trainNumberQuery).First().Text())
		brand = strings.TrimSpace(s.Find(".brand-logo").First().Text())
		duration = strings.ReplaceAll(strings.TrimSpace(s.Find(".travel-time-value").First().Text()), "h", ":")
		priceParts := strings.TrimSpace(s.Find(".price-parts").First().Text())
		priceParts = strings.ReplaceAll(priceParts, "\u00a0", "")
		priceParts = strings.ReplaceAll(priceParts, "zł", "")
		priceParts = strings.ReplaceAll(priceParts, ",", ".")

		priceParsed, err := strconv.ParseFloat(priceParts, 64)
		if err != nil {
			panic(err)
		}
		price = priceParsed
		conn := Connection{
			DepartureTime: from,
			ArrivalTime:   to,
			Price:         price,
			Name:          brand + "-" + train,
			Duration:      duration,
		}
		connections = append(connections, conn)
	})
	return connections, nil
}
