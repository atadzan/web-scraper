package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"os"
	"time"
)

type PokemonProduct struct {
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
	var pokemonProducts []PokemonProduct
	// initializing the list of pages to scrape with an empty slice
	var pagesToScrape []string

	// the first page to scrape
	pageToScrape := "https://scrapeme.live/shop/page/1/"

	// initializing the list of pages discovered with a pageToScrape
	pagesDiscovered := []string{pageToScrape}

	// current iteration
	i := 1
	// max page to scrape
	limit := 5

	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64)"
	c.OnHTML("a.page-numbers", func(e *colly.HTMLElement) {
		// discovering a new page
		newPaginationLink := e.Attr("href")

		// if the page discovered is new
		if !contains(pagesToScrape, newPaginationLink) {
			// if the page discovered should be scraped
			if !contains(pagesDiscovered, newPaginationLink) {
				pagesToScrape = append(pagesToScrape, newPaginationLink)
			}
			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
		}
	})
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		//initializing a new pokemonProduct instance
		var pokemonProduct PokemonProduct

		// scraping the data of interest
		pokemonProduct.Url = e.ChildAttr("a", "href")
		pokemonProduct.Image = e.ChildAttr("img", "src")
		pokemonProduct.Name = e.ChildText("h2")
		pokemonProduct.Price = e.ChildText(".price")

		// adding the product instance with scraped data to the list of products
		pokemonProducts = append(pokemonProducts, pokemonProduct)
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
	for _, pokemonProduct := range pokemonProducts {
		// converting a PokemonProduct to an array of strings
		record := []string{
			pokemonProduct.Url,
			pokemonProduct.Image,
			pokemonProduct.Name,
			pokemonProduct.Price,
		}

		// adding a CSV record to the output file
		writer.Write(record)
	}
	defer writer.Flush()
	fmt.Println("Duration: ", time.Since(start))
}
