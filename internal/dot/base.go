// Package dot provides an interface for generating dot language text.
package dot

import (
	"fmt"
	"strings"
)

type (
	// Dot can be converted into dot language.
	Dot interface {
		// AsDot generates dot expression.
		AsDot() string
	}
	// Attrs is a list of the Attr.
	// AsDot generates string as attr_list.
	Attrs interface {
		Dot
		// Len returns the length of this.
		Len() int
		// Add appends Attr.
		Add(attr Attr) Attrs
		// Get returns Attr at the index i.
		// Returns false if out of index range.
		Get(i int) (Attr, bool)
	}
	// Node is a node of dot.
	Node interface {
		Dot
		// Name returns the name of the node.
		Name() string
		// Attrs returns the attributes of the node.
		Attrs() Attrs
	}
	// Nodes is a list of the Node.
	Nodes interface {
		// Add appends Node.
		Add(node Node) Nodes
		// Len returns the length of this.
		Len() int
		// Get returns Node at the index i.
		// Returns false if out of index range.
		Get(i int) (Node, bool)
	}
	// Edge is an edge of dot.
	Edge interface {
		Dot
		// Start returns the start node of the edge.
		Start() Node
		// End returns the end node of the edge.
		End() Node
		// Attrs returns the attributes of the edge.
		Attrs() Attrs
	}
	// Edges is a list of the Edge.
	Edges interface {
		// Add appends Edge.
		Add(edge Edge) Edges
		// Len returns the length of this.
		Len() int
		// Get returns Edge at the index i.
		// Returns false if out of index range.
		Get(i int) (Edge, bool)
	}
	// Digraph is a digraph of dot.
	Digraph interface {
		Dot
		// Attrs returns the attributes of this.
		Attrs() Attrs
		// Nodes returns the nodes of this.
		Nodes() Nodes
		// Edges returns the edges of this.
		Edges() Edges
	}
)

type (
	attrs struct {
		v []Attr
	}
)

// NewAttrs creates a new Attrs.
func NewAttrs() Attrs {
	return &attrs{
		v: []Attr{},
	}
}

func (s attrs) AsDot() string {
	if s.Len() == 0 {
		return ""
	}
	v := make([]string, s.Len())
	for i, x := range s.v {
		v[i] = x.AsDot()
	}
	return fmt.Sprintf("[%s]", strings.Join(v, " "))
}
func (s attrs) Len() int { return len(s.v) }
func (s attrs) Get(i int) (Attr, bool) {
	if i < 0 || i >= s.Len() {
		return nil, false
	}
	return s.v[i], true
}
func (s *attrs) Add(attr Attr) Attrs {
	s.v = append(s.v, attr)
	return s
}

type (
	node struct {
		name  string
		attrs Attrs
	}
)

// NewNode creates a new Node.
// Returns an error if name is empty.
func NewNode(name string) (Node, error) {
	if name == "" {
		return nil, ErrNodeNameEmpty
	}
	return &node{
		name:  name,
		attrs: NewAttrs(),
	}, nil
}

func (s node) Name() string { return s.name }
func (s node) Attrs() Attrs { return s.attrs }
func (s node) AsDot() string {
	b := NewStringBuilder()
	b.Write(s.Name())
	if s.Attrs().Len() > 0 {
		b.Write(" ")
		b.Write(s.Attrs().AsDot())
	}
	return b.String()
}

type (
	nodes struct {
		v []Node
	}
)

// NewNodes creates a new Nodes.
func NewNodes() Nodes {
	return &nodes{
		v: []Node{},
	}
}

func (s *nodes) Add(node Node) Nodes {
	s.v = append(s.v, node)
	return s
}
func (s nodes) Len() int { return len(s.v) }
func (s nodes) Get(i int) (Node, bool) {
	if i < 0 || i >= s.Len() {
		return nil, false
	}
	return s.v[i], true
}

type (
	edge struct {
		start Node
		end   Node
		attrs Attrs
	}
)

// NewEdge creates a new Edge.
// Returns an error if start and end are nil.
func NewEdge(start, end Node) (Edge, error) {
	if start == nil || end == nil {
		return nil, ErrInvalidEdge
	}
	return &edge{
		start: start,
		end:   end,
		attrs: NewAttrs(),
	}, nil
}

func (s edge) AsDot() string {
	b := NewStringBuilder()
	b.Write(fmt.Sprintf("%s -> %s", s.start.Name(), s.end.Name()))
	if s.attrs.Len() > 0 {
		b.Write(" ")
		b.Write(s.attrs.AsDot())
	}
	return b.String()
}
func (s edge) Start() Node  { return s.start }
func (s edge) End() Node    { return s.end }
func (s edge) Attrs() Attrs { return s.attrs }

type (
	edges struct {
		v []Edge
	}
)

// NewEdges creates a new Edges.
func NewEdges() Edges {
	return &edges{
		v: []Edge{},
	}
}

func (s *edges) Add(edge Edge) Edges {
	s.v = append(s.v, edge)
	return s
}
func (s edges) Len() int { return len(s.v) }
func (s edges) Get(i int) (Edge, bool) {
	if i < 0 || i >= s.Len() {
		return nil, false
	}
	return s.v[i], true
}

type (
	digraph struct {
		attrs Attrs
		nodes Nodes
		edges Edges
	}
)

// NewDigraph creates a new Digraph.
func NewDigraph() Digraph {
	return &digraph{
		attrs: NewAttrs(),
		nodes: NewNodes(),
		edges: NewEdges(),
	}
}

func (s digraph) AsDot() string {
	b := NewStringBuilder()
	b.WriteLine("digraph {")
	for i := 0; i < s.Attrs().Len(); i++ {
		x, _ := s.attrs.Get(i)
		b.Indent()
		b.WriteLine(x.AsDot())
	}
	for i := 0; i < s.Nodes().Len(); i++ {
		x, _ := s.nodes.Get(i)
		b.Indent()
		b.WriteLine(x.AsDot())
	}
	for i := 0; i < s.Edges().Len(); i++ {
		x, _ := s.edges.Get(i)
		b.Indent()
		b.WriteLine(x.AsDot())
	}
	b.Write("}")
	return b.String()
}
func (s digraph) Attrs() Attrs { return s.attrs }
func (s digraph) Nodes() Nodes { return s.nodes }
func (s digraph) Edges() Edges { return s.edges }
