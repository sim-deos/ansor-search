package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	search "github.com/sim-deos/ansor-search/search"
)

func main() {
	key := "AIzaSyDRh6Xh8iWIYcvbNcJFfDu7AUfoC0f8wWw"
	cx := "e5bf4846021224002"
	searcher := search.NewSearcher(key, cx)

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
