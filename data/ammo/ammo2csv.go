package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No input file")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanner.Scan()
		text, err := url.PathUnescape(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		text = strings.Replace(text, "/", " ", -1)
		text = strings.Replace(text, "?", " ", -1)
		text = strings.Replace(text, "&", " ", -1)
		fmt.Println(text)
		scanner.Scan()
		scanner.Scan()
		scanner.Scan()
		scanner.Scan()
		scanner.Scan()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
