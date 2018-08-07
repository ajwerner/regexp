package regexp

import (
	"fmt"
	"unicode/utf8"
)

type parser struct {
	in    string
	pos   int
	width int
}

func (p *parser) next() (r rune) {
	r, p.width = utf8.DecodeRuneInString(p.in[p.pos:])
	p.pos += p.width
	return
}

func (p *parser) backup() {
	p.pos -= p.width
	p.width = 0
}

func (p *parser) peek() (r rune) {
	r = p.next()
	p.backup()
	return
}

// expr is a graph of nodes with a start and an end
type expr struct{ start, end nodeBuilder }

func (e expr) isEmpty() bool { return e == expr{} }

func parse(in string) (start, end node, err error) {
	p := parser{in: in}
	e, err := parseClause(&p)
	if err != nil {
		return
	}
	if p.pos != len(in) {
		p.backup()
		err = fmt.Errorf("illegal %c at pos %d", p.peek(), p.pos)
		return
	}
	term := &terminalNode{}
	e = concatNode(e, term)
	return e.start, e.end, nil
}

// parseClause parses a sequence of terms connected by '|' or concatenation.
func parseClause(p *parser) (e expr, err error) {
	for {
		switch r := p.peek(); r {
		case utf8.RuneError, ')':
			return
		case '|':
			if e, err = parsePipe(p, e); err != nil {
				return
			}
		default:
			var next expr
			if next, err = parseTerm(p); err != nil {
				return
			}
			e = concat(e, next)
		}
	}
}

// parsePipe returns an expr which represents the or of lhs and the next clause
func parsePipe(p *parser, lhs expr) (e expr, err error) {
	p.next()
	if lhs.isEmpty() {
		err = fmt.Errorf("invalid empty expression to the left of %c at pos %d %#v %#v", p.peek(), p.pos, lhs.start, lhs.end)
		return
	}
	if e, err = parseClause(p); err != nil {
		return
	}
	e.start.addEpsilon(lhs.start)
	n := &pipeNode{}
	e = concatNode(e, n)
	lhs = concatNode(lhs, n)
	return
}

// parseTerm parses the next single expr to add to a concatenation
func parseTerm(p *parser) (e expr, err error) {
	var pf func(*parser) (expr, error)
	switch r := p.peek(); r {
	case '?', '+', '*':
		return expr{}, fmt.Errorf("Invalid character %c at pos %d", r, p.pos)
	case '(':
		pf = parseSubexp
	case '\\':
		pf = parseEscape
	case '.':
		pf = parseDot
	default:
		pf = parseLiteral
	}
	if e, err = pf(p); err != nil {
		return
	}
	switch r := p.peek(); r {
	case '?', '+', '*':
		return parseMeta(p, e)
	}
	return
}

func parseEscape(p *parser) (expr, error) {
	p.next()
	switch r := p.peek(); r {
	case '?', '*', '+', '(', ')', '|':
		return parseLiteral(p)
	default:
		return expr{}, fmt.Errorf("unknown escape sequence \\%c", r)
	}
}

func parseSubexp(p *parser) (e expr, err error) {
	start := p.pos
	p.next()
	if e, err = parseClause(p); err != nil {
		return
	}
	if r := p.next(); r != ')' {
		err = fmt.Errorf("unterminated subexp starting at pos %d", start)
	}
	lp := &leftParenNode{}
	rp := &rightParenNode{}
	e = concatNode(concat(expr{lp, lp}, e), rp)
	return
}

func parseDot(p *parser) (expr, error) {
	p.next()
	n := &dotNode{}
	return expr{n, n}, nil
}

func parseLiteral(p *parser) (expr, error) {
	n := &literalNode{r: p.next()}
	return expr{n, n}, nil
}

func parseMeta(p *parser, term expr) (e expr, err error) {
	mc := p.next()
	var n nodeBuilder
	switch mc {
	case '+':
		n = &plusNode{}
		n.addEpsilon(term.start)
	case '*':
		n = &starNode{}
		n.addEpsilon(term.start)
		term.start.addEpsilon(n)
	case '?':
		n = &qmarkNode{}
		term.start.addEpsilon(n)
	}
	e = concatNode(term, n)
	return
}

func concatNode(e expr, n nodeBuilder) expr {
	return concat(e, expr{n, n})
}

func concat(e expr, next expr) expr {
	if e.isEmpty() {
		return next
	}
	if next.isEmpty() {
		return e
	}
	if mb, ok := e.end.(matchBuilder); ok {
		mb.setNext(next.start)
	} else {
		e.end.addEpsilon(next.start)
	}
	e.end = next.end
	return e
}
