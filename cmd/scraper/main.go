package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sts2/internal"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const (
	untappedBaseURL = "https://sts2.untapped.gg"
	untappedCardsURL = "https://sts2.untapped.gg/en/cards"
	cardLinkSelector = `a.ugg-link-link[href^="/en/cards/"]`
)

var cards []internal.Card

func main() {
	c := colly.NewCollector(
		// Prevents redownloading pages even after restart
		colly.CacheDir("./card_cache"),
		colly.CacheExpiration(time.Hour*24),
		)
	if err := crawl(c); err != nil {
		log.Fatal(err)
	}

	c.Wait()

	if err := writeJSON("cards.json", cards); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Wrote", len(cards), "cards to cards.json")
}

func crawl(c *colly.Collector) error {
	// New collector to get card detail
	detailCollector := c.Clone()

	c.OnHTML(cardLinkSelector, func(h *colly.HTMLElement) {
		fmt.Println("Found card link: ", h.Attr("href"))
		cardUrl := untappedBaseURL+h.Attr("href")
		detailCollector.Visit(cardUrl)
	})

	detailCollector.OnHTML(".page-module-scss-module__kK2N2a__container", func(h *colly.HTMLElement) {
		// details container
		card, err := collectCardDetails(h)
		if err != nil {
			fmt.Println("Error collecting card data: ", err)
			return
		}

		cards = append(cards, card)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	return c.Visit(untappedBaseURL)
}
func collectCardDetails(h *colly.HTMLElement) (Card, error) {
	var card Card

	card.URL = h.Request.URL.String()

	card.Title = strings.TrimSpace(h.ChildText("h1"))
	if card.Title == "" {
		return Card{}, fmt.Errorf("title not found")
	}

	h.ForEach("label", func(_ int, el *colly.HTMLElement) {
		labelText := strings.TrimSpace(el.Text)

		switch labelText {

		case "DESCRIPTION":
			desc := strings.TrimSpace(el.DOM.Next().Text())
			card.Description = desc

		case "CARD DETAILS":
			el.DOM.Next().Find("div").Each(func(_ int, detail *goquery.Selection) {
				key := strings.TrimSpace(detail.Find("label").Text())
				value := strings.TrimSpace(detail.Find("span").Text())

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

		case "SOURCE":
			card.Source = strings.TrimSpace(el.DOM.Next().Text())
		}
	})

	card.Img = h.ChildAttr(`img[class*="cardImage"]`, "src")

	return card, nil
}

func writeJSON(filename string, data any) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}
