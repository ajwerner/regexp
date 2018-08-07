// Package regexp is a basic regexp library.
//
// The library was implemented to explore regexp parsing and NFA representation.
package regexp

// Regexp is a compiled regular expression.
type Regexp struct {
	start, term node
	pattern     string
}

// Compile compiles a pattern into Regexp.
func Compile(pattern string) (r *Regexp, err error) {
	r = &Regexp{pattern: pattern}
	if r.start, r.term, err = parse(pattern); err != nil {
		return nil, err
	}
	return
}

// MustCompile invokes Compile and panics if an error is returned.
func MustCompile(pattern string) *Regexp {
	r, err := Compile(pattern)
	if err != nil {
		panic(err)
	}
	return r
}

// MatchString returns true if the the passed input matches the pattern.
func (r *Regexp) MatchString(input string) bool {
	cur, next := nodeSet{}, nodeSet{}
	cur.add(r.start)
	for _, r := range input {
		for n := range cur {
			if m, ok := n.(matchNode); ok && m.matches(r) {
				next.add(m.next())
			}
		}
		if len(next) == 0 {
			return false
		}
		cur.clear()
		cur, next = next, cur
	}
	_, hasTerminal := cur[r.term]
	return hasTerminal
}

////////////////////////////////////////////////////////////////////////////////
// nodeSet
////////////////////////////////////////////////////////////////////////////////

// nodeSet is used to evaluate the NFA
type nodeSet map[node]struct{}

func (s nodeSet) add(n node) {
	if _, exists := s[n]; exists {
		return
	}
	s[n] = struct{}{}
	for _, e := range n.epsilons() {
		s.add(e)
	}
}

func (s nodeSet) clear() {
	for n := range s {
		delete(s, n)
	}
}
