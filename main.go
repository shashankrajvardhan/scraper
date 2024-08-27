package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	fetchURL := "https://www.imdb.com/list/ls033609554"
	fileName := "disney-movie.csv"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("ERROR: Could not create file %q: %s\n", fileName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"sl. No", "Movie Name", "Release Year", "Certificate", "Genre", "Running time", "Rating", "Number of Votes", "Gross"})

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.OnHTML(`.lister-item-content`, func(e *colly.HTMLElement) {
		number := e.ChildText(".lister-item-index")
		name := e.ChildText(".lister-item-header a")
		year := e.ChildText(".lister-item-year")
		certificate := e.ChildText(".certificate")
		genre := e.ChildText(".genre")
		runtime := e.ChildText(".runtime")
		rating := e.ChildText(".ipl-rating-star.small .ipl-rating-star__rating")
		vote := e.ChildAttr("span[name=nv]", "data-value")
		gross := e.ChildText(".text-muted:contains('Gross') ~ span[name=nv]")

		writer.Write([]string{
			number,
			name,
			year,
			certificate,
			genre,
			runtime,
			rating,
			vote,
			gross,
		})
	})

	c.Visit(fetchURL)
	fmt.Println("End of scraping: ", fetchURL)
}
