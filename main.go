package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"q/querylist"
)

const (
	FLAG_DELETE = "d"
)

func main() {

	filePath := os.Getenv("QUERYLIST_FILE_PATH")
	if filePath == "" {
		log.Fatal("[MISSING ENV] QUERYLIST_FILE_PATH is not set.")
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ql, err := querylist.Load(file)
	if err != nil {
		panic(err)
	}

	flag.Func(FLAG_DELETE, "-d <queryname>", func(queryName string) error {
		if ql.Delete(queryName) {
			querylist.Flush(ql, file)
			fmt.Printf("Deleted query: %q\n", queryName)
		} else {
			fmt.Printf("❌ Could not find query: %q\n", queryName)
		}
		return nil
	})

	flag.Parse()

	// handle when only positional args exist
	if flag.NFlag() == 0 {
		args := flag.Args()

		// usage: q
		if len(args) == 0 {
			querylist.DisplayAll(ql)
		}

		// usage: q [key]
		if len(args) == 1 {
			querylist.DisplayQuery(ql, args[0])
		}

		// usage: q [key] [value]
		if len(args) == 2 {
			ql.Add(args[0], args[1])
			querylist.Flush(ql, file)
			querylist.DisplayQuery(ql, args[0])
		}

		return
	}

}
