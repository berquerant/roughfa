package dot

type (
	baseFaDotBuilder struct {
		startStates  []string
		states       []string
		acceptStates []string
		transitions  map[string]map[rune][]string
	}
)

func (baseFaDotBuilder) newAcceptState(state string) (Node, error) {
	n, err := NewNode(state)
	if err != nil {
		return nil, err
	}
	a, err := NewAttrBuilder().Name("shape").Value("doublecircle").Build()
	if err != nil {
		return nil, err
	}
	n.Attrs().Add(a)
	return n, nil
}

func (baseFaDotBuilder) newNormalState(state string) (Node, error) {
	n, err := NewNode(state)
	if err != nil {
		return nil, err
	}
	a, err := NewAttrBuilder().Name("shape").Value("circle").Build()
	if err != nil {
		return nil, err
	}
	n.Attrs().Add(a)
	return n, nil
}

func (baseFaDotBuilder) newEdge(start, end Node, label string) (Edge, error) {
	e, err := NewEdge(start, end)
	if err != nil {
		return nil, err
	}
	a, err := NewAttrBuilder().Name("label").Value(label).Build()
	if err != nil {
		return nil, err
	}
	e.Attrs().Add(a)
	return e, nil
}

func (baseFaDotBuilder) newStartStateFeature(g Digraph, startState Node) error {
	b := NewStringBuilder()
	isAlpha := func(r rune) bool {
		return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z'
	}
	// not duplicated but valid in dot
	for _, r := range NewUUIDFactory().NewUUID() {
		if !isAlpha(r) {
			continue
		}
		b.Write(string(r))
	}
	p, err := NewNode(b.String())
	if err != nil {
		return err
	}
	a, err := NewAttrBuilder().Name("shape").Value("point").Build()
	if err != nil {
		return err
	}
	p.Attrs().Add(a)
	g.Nodes().Add(p)
	e, err := NewEdge(p, startState)
	if err != nil {
		return err
	}
	g.Edges().Add(e)
	return nil
}

func (s baseFaDotBuilder) Build() (Dot, error) {
	var (
		g       = NewDigraph()
		nodeMap = map[string]Node{}
	)
	{
		// append graph attributes
		a, err := NewAttrBuilder().Name("rankdir").Value("LR").Build()
		if err != nil {
			return nil, err
		}
		g.Attrs().Add(a)
	}

	// append states
	for _, x := range s.acceptStates {
		n, err := s.newAcceptState(x)
		if err != nil {
			return nil, err
		}
		nodeMap[x] = n
	}
	for _, x := range s.states {
		if _, ok := nodeMap[x]; ok {
			// exclude accept states
			continue
		}
		n, err := s.newNormalState(x)
		if err != nil {
			return nil, err
		}
		nodeMap[x] = n
	}

	// apply start states feature
	for _, startState := range s.startStates {
		n, ok := nodeMap[startState]
		if !ok {
			return nil, ErrMissingState
		}
		if err := s.newStartStateFeature(g, n); err != nil {
			return nil, err
		}
	}

	for _, n := range nodeMap {
		g.Nodes().Add(n)
	}

	// append edges
	for st, v := range s.transitions {
		start, ok := nodeMap[st]
		if !ok {
			return nil, ErrInvalidEdge
		}
		for to, eds := range v {
			for _, ed := range eds {
				end, ok := nodeMap[ed]
				if !ok {
					return nil, ErrInvalidEdge
				}
				e, err := s.newEdge(start, end, string(to))
				if err != nil {
					return nil, err
				}
				g.Edges().Add(e)
			}
		}
	}
	return g, nil
}
