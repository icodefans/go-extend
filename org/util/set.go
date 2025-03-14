package util

type Set map[string]struct{}

func (s Set) Append(k string) {
	s[k] = struct{}{}
}

func (s Set) Remove(k string) {
	delete(s, k)
}

func (s Set) Exist(k string) bool {
	_, ok := s[k]
	return ok
}
