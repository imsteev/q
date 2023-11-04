package querylist

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func DisplayQuery(q *QueryList, key string) {
	if query := q.Get(key); query != nil {
		fmt.Printf("[%s]\t%s\n", query.Key, query.Val)
	} else {
		fmt.Println("❌ QUERY NOT FOUND")

	}
}

func DisplayAll(q *QueryList) {
	fmt.Println("Queries")
	i := 1
	for _, query := range q.queries {
		fmt.Printf("%d. %s\n", i, query.Key)
		i++
	}
}

// used for serialization
type medium struct {
	Queries []*query `json:"queries"`
}

// Load reads from disk into in-memory QueryList.
func Load(f *os.File) (*QueryList, error) {
	var ql = New()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
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
		ql.Add(q.Key, q.Val)
	}

	return ql, nil
}

// Flush writes the contents of in-memory QueryList onto disk.
func Flush(q *QueryList, f *os.File) error {
	if err := f.Truncate(0); err != nil {
		return err
	}

	// Set I/O offset to the start. Truncate does NOT automatically set it.
	f.Seek(0, 0)

	seen := map[string]bool{}

	var m medium
	for _, query := range q.queries {
		query := query // :woozy i don't trust myself
		if !seen[query.Key] {
			m.Queries = append(m.Queries, query)
			seen[query.Key] = true
		}
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	return err
}
