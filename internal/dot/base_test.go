package dot_test

import (
	"testing"

	"github.com/berquerant/roughfa/internal/dot"
	"github.com/stretchr/testify/assert"
)

type (
	mockDot struct {
		v string
	}
	mockAttr struct {
		mockDot
	}
	mockNode struct {
		mockDot
		name  string
		attrs dot.Attrs
	}
	mockEdge struct {
		mockDot
		start dot.Node
		end   dot.Node
		attrs dot.Attrs
	}
)

func (s mockDot) AsDot() string     { return s.v }
func (mockAttr) Name() string       { return "" }
func (mockAttr) Value() string      { return "" }
func (s mockNode) Name() string     { return s.name }
func (s mockNode) Attrs() dot.Attrs { return s.attrs }
func (s mockEdge) Start() dot.Node  { return s.start }
func (s mockEdge) End() dot.Node    { return s.end }
func (s mockEdge) Attrs() dot.Attrs { return s.attrs }

func newMockAttr(v string) dot.Attr { return &mockAttr{mockDot{v: v}} }
func newMockNode(v, name string, attrs dot.Attrs) dot.Node {
	return &mockNode{
		mockDot: mockDot{v: v},
		name:    name,
		attrs:   attrs,
	}
}
func newMockEdge(v string, start, end dot.Node, attrs dot.Attrs) dot.Edge {
	return &mockEdge{
		mockDot: mockDot{v: v},
		start:   start,
		end:     end,
		attrs:   attrs,
	}
}

type digraphTestcase struct {
	name  string
	attrs []dot.Attr
	nodes []dot.Node
	edges []dot.Edge
	want  string
}

func (s digraphTestcase) test(t *testing.T) {
	d := dot.NewDigraph()
	for _, x := range s.attrs {
		d.Attrs().Add(x)
	}
	for _, x := range s.nodes {
		d.Nodes().Add(x)
	}
	for _, x := range s.edges {
		d.Edges().Add(x)
	}
	assert.Equal(t, s.want, d.AsDot())
}

func TestDigraph(t *testing.T) {
	for _, tc := range []*digraphTestcase{
		{
			name: "empty",
			want: `digraph {
}`,
		},
		{
			name: "all",
			attrs: []dot.Attr{
				newMockAttr("a1"),
			},
			nodes: []dot.Node{
				newMockNode("n1", "", nil),
				newMockNode("n2", "", nil),
			},
			edges: []dot.Edge{
				newMockEdge("n1n2", nil, nil, nil),
			},
			want: `digraph {
  a1
  n1
  n2
  n1n2
}`,
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

type edgeTestcase struct {
	name  string
	start dot.Node
	end   dot.Node
	attrs []dot.Attr
	want  string
	err   error
}

func (s edgeTestcase) test(t *testing.T) {
	e, err := dot.NewEdge(s.start, s.end)
	assert.Equal(t, s.err, err)
	if err != nil {
		return
	}
	for _, x := range s.attrs {
		e.Attrs().Add(x)
	}
	assert.Equal(t, s.want, e.AsDot())
}

func TestEdge(t *testing.T) {
	for _, tc := range []*edgeTestcase{
		{
			name:  "invalid edge",
			start: newMockNode("", "n", nil),
			err:   dot.ErrInvalidEdge,
		},
		{
			name:  "no attrs",
			start: newMockNode("", "n1", nil),
			end:   newMockNode("", "n2", nil),
			want:  "n1 -> n2",
		},
		{
			name:  "with attrs",
			start: newMockNode("", "n1", nil),
			end:   newMockNode("", "n2", nil),
			attrs: []dot.Attr{
				newMockAttr("a1"),
			},
			want: "n1 -> n2 [a1]",
		},
	} {
		t.Run(tc.name, tc.test)
	}
}

type nodeTestcase struct {
	cname string
	name  string
	attrs []dot.Attr
	want  string
	err   error
}

func (s nodeTestcase) test(t *testing.T) {
	n, err := dot.NewNode(s.name)
	assert.Equal(t, s.err, err)
	if err != nil {
		return
	}
	for _, x := range s.attrs {
		n.Attrs().Add(x)
	}
	assert.Equal(t, s.want, n.AsDot())
}

func TestNode(t *testing.T) {
	for _, tc := range []*nodeTestcase{
		{
			cname: "empty name",
			err:   dot.ErrNodeNameEmpty,
		},
		{
			cname: "no attrs",
			name:  "n",
			want:  "n",
		},
		{
			cname: "with attrs",
			name:  "n",
			attrs: []dot.Attr{
				newMockAttr("a1"),
			},
			want: "n [a1]",
		},
	} {
		t.Run(tc.cname, tc.test)
	}
}

type attrsTestcase struct {
	name  string
	attrs []dot.Attr
	want  string
}

func (s attrsTestcase) test(t *testing.T) {
	a := dot.NewAttrs()
	for _, x := range s.attrs {
		a.Add(x)
	}
	assert.Equal(t, s.want, a.AsDot())
}

func TestAttrs(t *testing.T) {
	for _, tc := range []*attrsTestcase{
		{
			name: "empty",
		},
		{
			name: "an attr",
			attrs: []dot.Attr{
				newMockAttr("a1"),
			},
			want: "[a1]",
		},
		{
			name: "two attrs",
			attrs: []dot.Attr{
				newMockAttr("a1"),
				newMockAttr("a2"),
			},
			want: "[a1 a2]",
		},
	} {
		t.Run(tc.name, tc.test)
	}
}
