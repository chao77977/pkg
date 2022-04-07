package functool

// StringSet is an abstract data type with uniq value.
type StringSet map[string]struct{}

func NewStringSet() StringSet {
	return make(StringSet)
}

func NewStringSetWithValue(vs []string) StringSet {
	ss := NewStringSet()
	for _, s := range vs {
		if s != "" {
			ss.Add(s)
		}
	}

	return ss
}

func (s StringSet) Add(x string) {
	s[x] = struct{}{}
}

func (s StringSet) Remove(x string) {
	if s.Contains(x) {
		delete(s, x)
	}
}

func (s StringSet) Contains(x string) bool {
	_, ok := s[x]
	return ok
}

func (s StringSet) Len() int {
	return len(s)
}

func (s StringSet) List() []string {
	list := make([]string, s.Len())

	i := 0
	for k, _ := range s {
		list[i] = k
		i++
	}

	return list
}

func (s StringSet) Union(x StringSet) {
	if x.Len() != 0 {
		for k := range x {
			s.Add(k)
		}
	}
}
