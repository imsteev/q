package querylist

type QueryList struct {
	queries []*query
}

func New() *QueryList {
	return &QueryList{queries: []*query{}}
}

type query struct {
	Key string
	Val string
}

func (q *QueryList) Items() []*query {
	return q.queries
}

func (q *QueryList) Get(key string) *query {
	var lastSeen *query
	for _, query := range q.queries {
		if query.Key == key {
			lastSeen = query
		}
	}
	return lastSeen
}

func (q *QueryList) Add(key, val string) {
	q.queries = append(q.queries, &query{Key: key, Val: val})
}

func (q *QueryList) Delete(key string) bool {
	var filtered []*query
	for _, query := range q.queries {
		if query.Key != key {
			filtered = append(filtered, query)
		}
	}
	deleted := len(filtered) != len(q.queries)
	q.queries = filtered
	return deleted
}
