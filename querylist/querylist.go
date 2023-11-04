package querylist

type QueryList struct {
	m map[string]string
}

func New() *QueryList {
	return &QueryList{m: map[string]string{}}
}

type query struct {
	Key string
	Val string
}

func (q *QueryList) Items() []*query {
	var items []*query
	for k := range q.m {
		items = append(items, q.Get(k))
	}
	return items
}

func (q *QueryList) Get(key string) *query {
	val, ok := q.m[key]
	if ok {
		return &query{Key: key, Val: val}
	}
	return nil
}

func (q *QueryList) Add(key, val string) {
	q.m[key] = val
}

func (q *QueryList) Delete(key string) {
	delete(q.m, key)
}
