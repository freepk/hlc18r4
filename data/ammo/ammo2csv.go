package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	//"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No input file")
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Println("query_path;query_id;param_name;param_value")
	s := bufio.NewScanner(f)
	for s.Scan() {
		s.Scan()
		b := []byte(s.Text())
		b = b[13:(len(b) - 9)]
		u, err := url.ParseRequestURI(string(b))
		if err != nil {
			log.Fatal(err)
		}
		q := u.Query()
		id, ok := q["query_id"]
		if !ok {
			log.Fatal("No query_id")
		}
		for k, v := range q {
			if k == "query_id" {
				continue
			}
			fmt.Printf("%s;%s;%s;%s\n", u.EscapedPath(), id, k, v)
		}
		s.Scan()
		s.Scan()
		s.Scan()
		s.Scan()
		s.Scan()
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}
