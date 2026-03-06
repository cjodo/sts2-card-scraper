package scraper

import (
	"fmt"
	"strings"

	"sts2/internal/models"

	"github.com/gocolly/colly/v2"
)

const relicURL = "https://sts2.untapped.gg/en/relics"

func CrawlRelics(c *colly.Collector) ([]models.Relic, error) {
	count := 0
	var relics []models.Relic

	c.OnHTML(`div[class$="__relic"]`, func(e *colly.HTMLElement) {

		relic := models.Relic{}

		relic.Name = strings.TrimSpace(e.ChildText("h2"))
		relic.Img = e.ChildAttr("img", "src")

		relic.Character = e.ChildText(`div[class*="__details"] span:nth-child(2)`)
		relic.Rarity = e.ChildText(`div[class*="__details"] span:nth-child(3)`)

		relic.Description = strings.TrimSpace(
			e.ChildText(`p[class*="__description"]`),
		)

		if relic.Name != "" {
			count += 1
			fmt.Printf("\rRelics scraped: %d", count)
			relics = append(relics, relic)
		}
	})

	err := c.Visit(relicURL)
	c.Wait()

	return relics, err
}
