package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

type Property1 struct {
	Url   string
	Image string
	Name  string
	Price string
}

func main() {
	var properties []Property1
	c := colly.NewCollector()

	c.OnHTML("div.Box-sc-19dsmxk-0", func(e *colly.HTMLElement) {
		//initializing a new property instance
		var property Property1
		property.Name = e.ChildText("h2.eTZbBZ")
		property.Price = e.ChildText("span.kHpqEN")
		// scraping the data of interest
		property.Url = e.ChildText("")
		//property.Image = e.ChildAttr("img", "src")

		// adding the product instance with scraped data to the list of products
		properties = append(properties, property)
	})
	c.Visit("https://dubai.dubizzle.com/en/property-for-rent/residential/apartmentflat/?filters=(bedrooms%3E%3D1%20AND%20bedrooms%3C%3D2)%20AND%20(neighborhoods.ids%3D191)")
	for _, p := range properties {
		fmt.Println(p.Name)
		fmt.Println(p.Price)
	}
}
