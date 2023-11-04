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

const (
	FLAG_ADD    = "add"
	FLAG_DELETE = "d"
	FLAG_LIST   = "l"
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

	flag.String(FLAG_ADD, "", "-add add a named query")
	flag.String(FLAG_DELETE, "", "-del delete query")
	flag.Bool(FLAG_LIST, false, "-l list all queries")
	flag.Parse()

	positional := flag.Args()

	// usage: <queryname>
	if flag.NFlag() == 0 && len(positional) == 1 {
		querylist.Display(ql, positional[0])
	}

	flag.Visit(func(f *flag.Flag) {
		name := f.Name

		if name == FLAG_ADD && len(positional) != 1 {
			panic("usage: -add <query name> [query]")
		}

		if name == FLAG_ADD {
			ql.Add(f.Value.String(), positional[0])
			querylist.Flush(ql, file)
			querylist.Display(ql, f.Value.String())
		}

		if name == FLAG_DELETE {
			ql.Delete(f.Value.String())
			querylist.Flush(ql, file)
		}

		if name == FLAG_LIST {
			querylist.Display(ql, "")
		}

	})
}
