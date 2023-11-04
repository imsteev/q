package main

import (
	"flag"
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

	flag.String("add", "", "-add add a named query")
	flag.String("del", "", "-del delete query")
	flag.String("view", "", "-view view query")
	flag.Bool("all", false, "-all list all queries")
	flag.Parse()
	positional := flag.Args()

	flag.Visit(func(f *flag.Flag) {
		name := f.Name

		if name == "add" && len(positional) != 1 {
			panic("usage: -add <query name> [query]")
		}
		if name == "add" {
			ql.Add(f.Value.String(), positional[0])
			querylist.Flush(ql, file)
			querylist.Display(ql, f.Value.String())
		}
		if name == "all" {
			querylist.Display(ql, "")
		}
		if name == "del" {
			ql.Delete(f.Value.String())
			querylist.Flush(ql, file)
		}
		if name == "view" {
			querylist.Display(ql, f.Value.String())
		}

	})
}
