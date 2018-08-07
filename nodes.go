package regexp

////////////////////////////////////////////////////////////////////////////////
// NFA nodes interfaces
////////////////////////////////////////////////////////////////////////////////

// node is an NFA node used to evaluate a regexp
type node interface {
	epsilons() []node
}

// matchNode is a node which might match a rune moving the automaton to next.
type matchNode interface {
	node
	matches(rune) bool
	next() node
}

// nodeBuilder is a parser-exposed interface for building nodes.
type nodeBuilder interface {
	node
	addEpsilon(node)
}

// matchBuilder is a parser-exposed interface for building matchNodes.
type matchBuilder interface {
	matchNode
	setNext(node)
}

////////////////////////////////////////////////////////////////////////////////
// base implementations
////////////////////////////////////////////////////////////////////////////////

// nodeBase is the base implementation of nodeBuilder
type nodeBase []node

func (n *nodeBase) epsilons() []node  { return []node(*n) }
func (n *nodeBase) addEpsilon(e node) { *n = append(*n, e) }

// matchBase is the base implementation of matchBuilder
type matchBase struct {
	nodeBase
	n node
}

func (n *matchBase) next() node        { return n.n }
func (n *matchBase) setNext(next node) { n.n = next }

////////////////////////////////////////////////////////////////////////////////
// concrete node types
////////////////////////////////////////////////////////////////////////////////

type (
	starNode       struct{ nodeBase }
	qmarkNode      struct{ nodeBase }
	plusNode       struct{ nodeBase }
	terminalNode   struct{ nodeBase }
	leftParenNode  struct{ nodeBase }
	rightParenNode struct{ nodeBase }
	pipeNode       struct{ nodeBase }
	dotNode        struct{ matchBase }
	literalNode    struct {
		matchBase
		r rune
	}
)

////////////////////////////////////////////////////////////////////////////////
// match node implementations
////////////////////////////////////////////////////////////////////////////////

var _, _ matchBuilder = (*dotNode)(nil), (*literalNode)(nil)

func (n dotNode) matches(_ rune) bool     { return true }
func (n literalNode) matches(r rune) bool { return r == n.r }
