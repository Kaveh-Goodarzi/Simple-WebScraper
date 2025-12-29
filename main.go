package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type Product struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Price string `json:"price"`
	Image string `json:"img_url"`
}

type Scraper struct {
	collector   *colly.Collector
	products    []Product
	visitedURLs sync.Map
	outputFile  string
	startTime   time.Time
}

func NewScraper() *Scraper {
	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})

	return &Scraper{
		collector:  c,
		products:   make([]Product, 0),
		outputFile: "products.csv",
		startTime:  time.Now(),
	}
}

func (s *Scraper) SetupCallbacks() {
	s.collector.OnHTML("li.product", s.handleProduct)
	s.collector.OnHTML("a.next", s.handlePagination)
	s.collector.OnScraped(s.handleScraped)
	s.collector.OnError(s.handleError)
}

// Change to get internet access
func (s *Scraper) handleProduct(e *colly.HTMLElement) {
	product := extractProductFromHTML(
		e.ChildText(".product-name"),
		e.ChildText(".price"),
		e.ChildAttr("img", "src"),
		e.ChildAttr("a", "href"),
	)

	s.products = append(s.products, product)
	log.Printf("Scraped: %s - %s", product.Name, product.Price)
}


func (s *Scraper) handlePagination(e *colly.HTMLElement) {
	nextPage := e.Attr("href")

	if _, found := s.visitedURLs.Load(nextPage);!found {
		s.visitedURLs.Store(nextPage, struct{}{})
		fmt.Printf("Moving to page: %s\n", nextPage)
		e.Request.Visit(nextPage)
	}
}

func (s *Scraper) handleScraped(r *colly.Response) {
	s.saveToCSV()
	s.printSummary()
}

func (s *Scraper) handleError(r *colly.Response, err error) {
	log.Printf("Error scraping %s: %v", r.Request.URL, err)
}

func (s *Scraper) saveToCSV() error {
	file, err := os.Create(s.outputFile)
	if err!= nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Name", "Price", "Image", "URL"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, product := range s.products {
		record := []string{product.Name, product.Price, product.Image, product.URL}
		
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	log.Printf("Data saved to %s", s.outputFile)
	return nil
}

func (s *Scraper) printSummary() {
	duration := time.Since(s.startTime)
	
	fmt.Println("\n" + "==================================================")
	fmt.Println("SCRAPING COMPLETE")
	fmt.Println("==================================================")
	fmt.Printf("Total products: %d\n", len(s.products))
	fmt.Printf("Time elapsed: %v\n", duration)
	fmt.Printf("Output file: %s\n", s.outputFile)
	fmt.Println("==================================================")
}

func (s *Scraper) Start(startURL string) error {
	s.SetupCallbacks()

	fmt.Printf("Starting scraper...\nTarget: %s\n", startURL)

	if err := s.collector.Visit(startURL); err != nil {
		return fmt.Errorf("failed to start scraping: %w", err)
	}

	s.collector.Wait()
	return nil
}

func main() {
	scraper := NewScraper()

	if err := scraper.Start("https://www.scrapingcourse.com/ecommerce"); err != nil {
		log.Fatal(err)
	}
}


// Add for testing with internet access
func extractProductFromHTML(name, price, image, url string) Product {
	return Product{
		Name:  name,
		Price: price,
		Image: image,
		URL:   url,
	}
}
