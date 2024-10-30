/*
Copyright © 2024 Lukasz Bombol
*/
package main

import (
	"github.com/belukasz/pkp-cli/app"
) 

func main() {
	scrapped := scrapper.ScrapeConnections(30, "06:00", "EIP", "krakow", "warszawa")
	scrapper.PrintTable(scrapped, true, 3)
}
