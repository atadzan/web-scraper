package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"time"
)

type Property struct {
	Url   string
	Image string
	Name  string
	Price string
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func main() {
	start := time.Now()
	var properties []Property
	// initializing the list of pages to scrape with an empty slice
	var pagesToScrape []string

	// the first page to scrape
	pageToScrape := "https://dubai.dubizzle.com/en/property-for-rent/residential/apartmentflat/?filters=(bedrooms%3E%3D1%20AND%20bedrooms%3C%3D2)%20AND%20(neighborhoods.ids%3D191)"

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	// current iteration
	i := 1
	// max page to scrape
	limit := 5

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	c.OnHTML("div.page-list", func(e *colly.HTMLElement) {
		// discovering a new page
		newPaginationLink := e.ChildAttr("a", "href")

		// if the page discovered is new
		if !contains(pagesToScrape, newPaginationLink) {
			// if the page discovered should be scraped
			if !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})
	c.OnHTML("div.Box-sc-19dsmxk-0", func(e *colly.HTMLElement) {
		//initializing a new property instance
		var property Property

		// scraping the data of interest
		//property.Url = e.ChildAttr("h2", "listing-title")
		//property.Image = e.ChildAttr("img", "src")
		property.Name = e.ChildText("h2")
		property.Price = e.ChildText("span")

		// adding the product instance with scraped data to the list of products
		properties = append(properties, property)
	})

	c.OnScraped(func(response *colly.Response) {
		// until there is stilla page to scrape
		if len(pagesToScrape) != 0 && i < limit {
			// getting the current page to scrape and removing it from the list
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			// incrementing the iteration counter
			i++

			c.Visit(pageToScrape)
		}
	})

	c.Visit(pageToScrape)

	// opening the CSV file
	file, err := os.Create("assets/advanced_products.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// initializing a file writer
	writer := csv.NewWriter(file)

	// writing the CSV headers
	headers := []string{
		"url",
		"image",
		"name",
		"price",
	}
	writer.Write(headers)

	// writing each Pokemon product as a CSV row
	for _, property := range properties {
		// converting a Property to an array of strings
		record := []string{
			property.Url,
			property.Image,
			property.Name,
			property.Price,
		}

		// adding a CSV record to the output file
		writer.Write(record)
	}
	defer writer.Flush()
	fmt.Println("Duration: ", time.Since(start))
}
