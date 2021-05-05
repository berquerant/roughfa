package dot

import "fmt"

type AttrType int

const (
	// WrappedAttrType is a type of Attr whose value is wrapped by ".
	WrappedAttrType AttrType = iota
	// RawAttrType is a type of Attr whose value is a raw string.
	RawAttrType
)

type (
	// Attr is an attribute of dot.
	Attr interface {
		Dot
		// Name returns the name of the attribute.
		Name() string
		// Value returns the value of the attribute.
		Value() string
	}

	attr struct {
		attrType AttrType
		name     string
		value    string
	}
)

func (s attr) Name() string  { return s.name }
func (s attr) Value() string { return s.value }
func (s attr) AsDot() string {
	switch s.attrType {
	case RawAttrType:
		return fmt.Sprintf("%s=%s", s.Name(), s.Value())
	case WrappedAttrType:
		return fmt.Sprintf(`%s="%s"`, s.Name(), s.Value())
	default:
		panic("unknown attr type")
	}
}

type (
	// AttrBuilder is a builder of Attr.
	AttrBuilder interface {
		// Name configures the name of Attr.
		// Required.
		Name(name string) AttrBuilder
		// Value configures the value of Attr.
		// Required.
		Value(value string) AttrBuilder
		// AttrType configures the type of Attr.
		AttrType(attrType AttrType) AttrBuilder
		// Build creates a new Attr.
		// Returns an error if some validations fails.
		Build() (Attr, error)
	}

	attrBuilder struct {
		name     string
		value    string
		attrType AttrType
	}
)

// NewAttrBuilder creates a new AttrBuilder.
func NewAttrBuilder() AttrBuilder { return &attrBuilder{} }

func (s *attrBuilder) Name(name string) AttrBuilder {
	s.name = name
	return s
}
func (s *attrBuilder) Value(value string) AttrBuilder {
	s.value = value
	return s
}
func (s *attrBuilder) AttrType(attrType AttrType) AttrBuilder {
	s.attrType = attrType
	return s
}
func (s *attrBuilder) Build() (Attr, error) {
	if s.name == "" {
		return nil, ErrAttrNameEmpty
	}
	return &attr{
		name:     s.name,
		value:    s.value,
		attrType: s.attrType,
	}, nil
}
