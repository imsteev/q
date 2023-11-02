package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

var QUERYLIST_FILE_PATH string

func init() {
	if envPath := os.Getenv("QUERYLIST_FILE"); envPath != "" {
		QUERYLIST_FILE_PATH = envPath
	} else {
		QUERYLIST_FILE_PATH = ".querylist"
	}
}

func main() {

	// Setup in-memory data structure
	file, err := os.OpenFile(QUERYLIST_FILE_PATH, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ql, err := parse(file)
	if err != nil {
		panic(err)
	}

	flag.String("add", "", "-add add a named query")
	flag.String("del", "", "-del delete query")
	flag.String("view", "", "-view view query")

	flag.Bool("clear", false, "-clear clear all queries")
	flag.Bool("all", false, "-all list all queries")

	flag.Parse()
	positional := flag.Args()

	flag.Visit(func(f *flag.Flag) {
		name := f.Name

		if name == "add" && len(positional) != 1 {
			panic("usage: -add <query name> [query]")
		}
		if name == "add" {
			ql.add(f.Value.String(), positional[0])
			ql.flush(file)
			ql.display(f.Value.String())
		}
		if name == "all" {
			ql.display("")
		}
		if name == "clear" {
			ql = QueryList{}
			ql.flush(file)
		}
		if name == "del" {
			ql.delete(f.Value.String())
			ql.flush(file)
		}
		if name == "view" {
			ql.display(f.Value.String())
		}

	})
}

type QueryList map[string]string

func (q QueryList) delete(key string) {
	delete(q, key)
}
func (q QueryList) add(key, val string) {
	q[key] = val
}

func (q QueryList) display(name string) {
	if name != "" {
		query, ok := q[name]
		if !ok {
			fmt.Println("❌ QUERY NOT FOUND")
		} else {
			fmt.Println("[QUERY NAME]")
			fmt.Println(name)
			fmt.Println("[QUERY]")
			fmt.Println(query)
		}
	} else {
		fmt.Println("[QUERY NAME]")
		for k := range q {
			fmt.Println(k)
		}
	}
}

type query struct {
	Key string
	Val string
}

// used for serialization
type medium struct {
	Queries []query `json:"queries"`
}

// flush writes the contents of in-memory QueryList onto disk.
func (q QueryList) flush(f *os.File) error {
	if err := f.Truncate(0); err != nil {
		return err
	}

	bytes, err := json.Marshal(toMedium(q))
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	return err
}

// parse reads contents from disk into in-memory QueryList.
func parse(f *os.File) (QueryList, error) {
	var ql = QueryList{}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return ql, err
	}

	// return early so we don't err out trying to parse nothing.
	if len(bytes) == 0 {
		return ql, nil
	}

	var m medium
	if err := json.Unmarshal(bytes, &m); err != nil {
		return nil, err
	}

	// construct the query list from serialized format
	for _, q := range m.Queries {
		ql.add(q.Key, q.Val)
	}

	return toQueryList(m), nil
}

func toMedium(ql QueryList) medium {
	var m medium
	for k, v := range ql {
		m.Queries = append(m.Queries, query{Key: k, Val: v})
	}
	return m
}

// construct the query list from serialized format
func toQueryList(m medium) QueryList {
	var ql = QueryList{}
	for _, q := range m.Queries {
		ql.add(q.Key, q.Val)
	}
	return ql
}
