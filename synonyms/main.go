package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalf("failed to search synonyms of %q: %v \n", word, err)
		}
		if len(syns) == 0 {
			log.Fatalf("no synonyms of %q: \n", word)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
