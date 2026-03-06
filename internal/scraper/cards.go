package scraper

import (
	"fmt"
	"strings"

	"sts2/internal/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const baseURL = "https://sts2.untapped.gg"

func CrawlCards(c *colly.Collector) ([]models.Card, error) {
	count := 0
	var cards []models.Card

	detail := c.Clone()

	seen := map[string]struct{}{}

	c.OnHTML(`a.ugg-link-link[href^="/en/cards/"]`, func(e *colly.HTMLElement) {
		url := baseURL + e.Attr("href")
		if _, ok := seen[url]; ok {
			return
		}
		seen[url] = struct{}{}
		detail.Visit(url)
	})

	detail.OnHTML(".page-module-scss-module__kK2N2a__container", func(e *colly.HTMLElement) {
		card, err := collectCard(e)
		if err != nil {
			return
		}
		count += 1
		fmt.Printf("\rCards scraped: %d", count)
		cards = append(cards, card)
	})

	// tier list pages containing links
	pages := []string{
		"/en/tier-list/cards/ironclad",
		"/en/tier-list/cards/silent",
		"/en/tier-list/cards/defect",
		"/en/tier-list/cards/necrobinder",
		"/en/tier-list/cards/regent",
		"/en/tier-list/cards/colorless",
	}

	for _, p := range pages {
		c.Visit(baseURL + p)
	}

	c.Wait()

	return cards, nil
}

func collectCard(e *colly.HTMLElement) (models.Card, error) {

	var card models.Card

	card.URL = e.Request.URL.String()
	card.Title = strings.TrimSpace(e.ChildText("h1"))

	e.ForEach("label", func(_ int, el *colly.HTMLElement) {

		switch strings.TrimSpace(el.Text) {

		case "DESCRIPTION":
			card.Description = strings.TrimSpace(el.DOM.Next().Text())

		case "CARD DETAILS":
			el.DOM.Next().Find("div").Each(func(_ int, d *goquery.Selection) {

				key := strings.TrimSpace(d.Find("label").Text())
				value := strings.TrimSpace(d.Find("span").Text())

				switch key {

				case "Character":
					card.Character = value

				case "Type":
					card.Type = value

				case "Cost":
					card.EnergyCost = value

				case "Rarity":
					card.Rarity = value
				}
			})
		}
	})

	card.Img = e.ChildAttr(`img[class*="cardImage"]`, "src")

	if card.Title == "" {
		return card, fmt.Errorf("invalid card")
	}

	return card, nil
}
