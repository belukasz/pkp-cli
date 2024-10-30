package scrapper

import (
	"os"
	"strings"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PrintTable(scrapped []DayScrape, workdays bool, limit int) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Date", "Depart", "Arrival", "Price", "Train", "Duration"})
	for _, ds := range scrapped {
		for i, conn := range ds.Connections {
			date := ds.RealDate().Format("Mon Jan02")
			day := strings.Split(date, " ")[0]
			shouldDisplay := (!workdays || !(workdays && (day == "Sun" || day == "Sat")))
			if i+1 < limit && shouldDisplay {
				t.AppendRow(table.Row{date, conn.DepartureTime, conn.ArrivalTime, conn.Price, conn.Name, conn.Duration})
			}
		}
		t.AppendSeparator()
	}
	t.Render()
}
