package main

import (
	"fmt"
	"log"
	"sts2/internal/scraper"
	"sts2/internal/writer"

	"github.com/gocolly/colly/v2"
)

func newCollector(cache string) *colly.Collector {
	c := colly.NewCollector(
		colly.CacheDir(cache),
		)

	return c
}

func main() {
	cardCrawler := newCollector("./cache/cards")
	relicCrawler := newCollector("./cache/relics")

	fmt.Println("\nscraping cards\n")
	cards, err := scraper.CrawlCards(cardCrawler)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("\nscraping relics\n")
	relics, err := scraper.CrawlRelics(relicCrawler)
	if err != nil {
		log.Fatal(err)
	}

	writer.Write("./cards.json", cards)
	writer.Write("./relics.json", relics)

	fmt.Println("Cards:", len(cards))
	fmt.Println("Relics:", len(relics))
}
