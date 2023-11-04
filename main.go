package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"q/querylist"
)

var QUERYLIST_FILE_PATH string

func init() {
	QUERYLIST_FILE_PATH = os.Getenv("QUERYLIST_FILE_PATH")
	if QUERYLIST_FILE_PATH == "" {
		log.Fatal("[MISSING ENV] QUERYLIST_FILE_PATH is not set.")
	}
}

const FLAG_DELETE = "d"

func main() {

	file, err := os.OpenFile(QUERYLIST_FILE_PATH, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ql, err := querylist.Load(file)
	if err != nil {
		panic(err)
	}

	flag.String(FLAG_DELETE, "", "-del delete query")
	flag.Parse()
	positional := flag.Args()

	if flag.NFlag() == 0 {

		// usage: q
		if len(positional) == 0 {
			querylist.DisplayAll(ql)
		}

		// usage: q [key]
		if len(positional) == 1 {
			querylist.DisplayQuery(ql, positional[0])
		}

		// usage: q [key] [value]
		if len(positional) == 2 {
			ql.Add(positional[0], positional[1])
			querylist.Flush(ql, file)
			querylist.DisplayQuery(ql, positional[0])
		}

		return
	}

	if flg := flag.Lookup(FLAG_DELETE); flg != nil {
		queryName := flg.Value.String()
		if ql.Delete(queryName) {
			querylist.Flush(ql, file)
			fmt.Printf("Deleted query: %q\n", queryName)
		} else {
			fmt.Printf("❌ Could not find query: %q\n", queryName)
		}
	}
}
