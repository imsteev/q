package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	QUERY_LIST_PATH = `querylist`
)

func main() {

	flag.String("add", "", "-add add a named query")
	flag.Bool("clear", false, "-clear clear all queries")
	flag.String("del", "", "-del delete query")
	flag.Bool("list", false, "-list list all queries")

	flag.Parse()
	positional := flag.Args()

	ql, qlf, err := parseQueryList()
	if err != nil {
		log.Fatal(err)
	}

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "add" && len(positional) != 1 {
			log.Fatal("usage: -add <query name> [query]")
		}
		if f.Name == "add" {
			ql.add(f.Value.String(), positional[0])
			ql.flush(qlf)
		}
		if f.Name == "clear" {
			ql = QueryList{}
			ql.flush(qlf)
		}
		if f.Name == "del" {
			// empty string means deleted
			ql.delete(f.Value.String())
			ql.flush(qlf)
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
	fmt.Println(q)
}

func (q QueryList) print() {
	for k, v := range q {
		fmt.Println(k, v)
	}
}

func (q QueryList) flush(f *os.File) error {
	var qj qjson
	for k, v := range q {
		qj.Queries = append(qj.Queries, query{key: k, val: v})
	}

	bytes, err := json.Marshal(qj)
	if err != nil {
		return err
	}
	if err := f.Truncate(0); err != nil {
		return err
	}
	_, err = f.Write(bytes)
	return err
}

type query struct {
	key string
	val string
}
type qjson struct {
	Queries []query `json:"queries"`
}

func parseQueryList() (QueryList, *os.File, error) {
	f, err := os.OpenFile(QUERY_LIST_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bytes, err := os.ReadFile(QUERY_LIST_PATH)
	if err != nil {
		log.Fatal(err)
	}

	ql := QueryList{}

	if len(bytes) != 0 {
		var qj qjson
		if err := json.Unmarshal(bytes, &qj); err != nil {
			return nil, nil, err
		}

		for _, q := range qj.Queries {
			ql.add(q.key, q.val)
		}
	}

	return ql, f, nil
}
