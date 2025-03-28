package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	gdenv "github.com/joho/godotenv"
	search "github.com/sim-deos/ansor-search/search"
)

func main() {
	err := gdenv.Load()
	if err != nil {
		log.Fatal("error loading .env files: ", err)
	}

	searcher := search.NewSearcher(os.Getenv("GSE_KEY"), os.Getenv("GSE_ENG"))

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("What would you like to search?")
	scanner.Scan()
	query := scanner.Text()

	res, err := searcher.Search(query)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range res.Items {
		fmt.Println(item.Title, ": ", item.Link)
	}
}
