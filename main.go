package main

import (
	"flag"
	"fmt"
	"os"
	"q/querylist"
)

var QUERYLIST_FILE_PATH string

func init() {
	if envPath := os.Getenv("QUERYLIST_FILE"); envPath != "" {
		QUERYLIST_FILE_PATH = envPath
	} else {
		QUERYLIST_FILE_PATH = ".querylist.json"
	}
}

const (
	FLAG_DELETE = "d"
)

func main() {

	// Setup in-memory data structure
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
			querylist.Display(ql, "")
		}

		// usage: q [key]
		if len(positional) == 1 {
			querylist.Display(ql, positional[0])
		}

		// usage: q [key] [value]
		if len(positional) == 2 {
			ql.Add(positional[0], positional[1])
			querylist.Flush(ql, file)
			querylist.Display(ql, positional[0])
		}
	}

	flag.Visit(func(f *flag.Flag) {
		if f.Name == FLAG_DELETE {
			ql.Delete(f.Value.String())
			querylist.Flush(ql, file)
			fmt.Printf("Deleted query: %q\n", f.Value.String())
		}

	})
}
