package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
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

	flag.String("add", "", "-add add a named query")
	flag.String("del", "", "-del delete query")

	flag.Bool("init", false, "-init setup default storage")
	flag.Bool("clear", false, "-clear clear all queries")
	flag.Bool("list", false, "-list list all queries")

	flag.Parse()
	positional := flag.Args()

	file, err := os.OpenFile(QUERYLIST_FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ql, err := parse(file)
	if err != nil {
		log.Fatal(err)
	}

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "add" && len(positional) != 1 {
			log.Fatal("usage: -add <query name> [query]")
		}
		if f.Name == "add" {
			ql.add(f.Value.String(), positional[0])
			ql.flush(file)
			ql.print()
		}
		if f.Name == "clear" {
			ql = QueryList{}
			ql.flush(file)
		}
		if f.Name == "del" {
			ql.delete(f.Value.String())
			ql.flush(file)
		}
		if f.Name == "init" {
			// create file
		}
		if f.Name == "list" {
			ql.print()
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

func (q QueryList) print() {
	fmt.Println("[QUERY NAME]")
	for k, _ := range q {
		fmt.Println(k)
	}
}

type query struct {
	Key       string
	Val       string
	Timestamp time.Time
}

type medium struct {
	Queries []query `json:"queries"`
}

// flush writes the contents of in-memory QueryList onto disk.
func (q QueryList) flush(f *os.File) error {
	if err := f.Truncate(0); err != nil {
		return err
	}

	var m medium
	for k, v := range q {
		m.Queries = append(m.Queries, query{Key: k, Val: v, Timestamp: time.Now()})
	}

	bytes, err := json.Marshal(m)
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
		return QueryList{}, err
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

	return ql, nil
}
