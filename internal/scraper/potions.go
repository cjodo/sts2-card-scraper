package scraper

import (
	"fmt"
	"strings"

	"sts2/internal/models"

	"github.com/gocolly/colly/v2"
)

const potionURL = "https://sts2.untapped.gg/en/potions"

func CrawlPotions(c *colly.Collector) ([]models.Potion, error) {
	count := 0
	var potions []models.Potion

	c.OnHTML(`div[class$="__relic"]`, func(e *colly.HTMLElement) {
		potion := models.Potion{}

		potion.Name = strings.TrimSpace(e.ChildText("h2"))
		potion.Img = e.ChildAttr("img", "src")

		potion.Character = e.ChildText(`div[class*="__details"] span:nth-child(2)`)
		potion.Rarity = e.ChildText(`div[class*="__details"] span:nth-child(3)`)

		// join all <span> in description to capture mechanics text
		descSpans := e.ChildTexts(`p[class*="__description"] span`)
		potion.Description = strings.Join(descSpans, " ")

		if potion.Name != "" {
			count++
			fmt.Printf("\rPotions scraped: %d", count)
			potions = append(potions, potion)
		}
	})

	if err := c.Visit(potionURL); err != nil {
		return nil, err
	}

	c.Wait()

	fmt.Println() // final newline after the live counter

	return potions, nil
}
