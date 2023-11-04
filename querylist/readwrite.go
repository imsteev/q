package querylist

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Display(q *QueryList, queryName string) {
	if queryName != "" {
		query, ok := q.m[queryName]
		if !ok {
			fmt.Println("❌ QUERY NOT FOUND")
		} else {
			fmt.Printf("✅ [%s]\t%s\n", queryName, query)
		}
	} else {
		fmt.Println("[QUERY]")
		for k := range q.m {
			fmt.Println(k)
		}
	}
}

// used for serialization
type medium struct {
	Queries []*query `json:"queries"`
}

// Parse reads contents from disk into in-memory QueryList.
func Parse(f *os.File) (*QueryList, error) {
	var ql = New()

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

	bytes, err := json.Marshal(toMedium(q))
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	return err
}

func toMedium(ql *QueryList) medium {
	var m medium
	for k, v := range ql.m {
		m.Queries = append(m.Queries, &query{Key: k, Val: v})
	}
	return m
}
