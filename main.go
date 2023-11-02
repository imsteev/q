package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	QUERY_LIST_PATH = `querylist`
)

func main() {

	flag.String("add", "", "-add your query")
	flag.String("del", "", "-add your query")
	flag.Bool("list", false, "-list")

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "add" {
			fmt.Println(f.Value)

			positional := flag.Args()
			if len(positional) != 1 {
				log.Fatal("usage: -add <query name> [query]")
			}

			query := positional[0]

			fmt.Println("adding query: %s", query)
		}
		if f.Name == "del" {
			fmt.Println(f.Value)
		}
		if f.Name == "list" {
			bytes, err := os.ReadFile(QUERY_LIST_PATH)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bytes))
		}
	})

}
