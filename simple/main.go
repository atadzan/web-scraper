package main

import (
	"encoding/csv"
	"github.com/gocolly/colly"
	"log"
	"os"
)

type PokemonProduct struct {
	Url   string
	Image string
	Name  string
	Price string
}

func main() {
	var pokemonProducts []PokemonProduct
	c := colly.NewCollector()

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
	c.Visit("https://scrapeme.live/shop/")

	// opening the CSV file
	file, err := os.Create("assets/simple_products.csv")
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
}
