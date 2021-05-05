package dot

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type (
	UUIDFactory interface {
		NewUUID() string
	}
	uuidFactory struct {
	}
)

func NewUUIDFactory() UUIDFactory { return &uuidFactory{} }

func (uuidFactory) NewUUID() string { return uuid.New().String() }

type (
	// StringBuilder is an utility for building string.
	// String returns a built string.
	StringBuilder interface {
		fmt.Stringer
		// Write appends string.
		Write(x string) (int, error)
		// WriteLine appends string, and a newline.
		WriteLine(x string) (int, error)
		// NewLine appends a newline.
		NewLine() error
		// Indent appends an indentation.
		Indent() error
		// IndentN invokes Indent n times.
		// No operation for a negative number.
		IndentN(n int) error
	}

	stringBuilder struct {
		b *strings.Builder
	}
)

const (
	newline = '\n'
	indent  = "  "
)

// NewStringBuilder creates a new StringBuilder.
func NewStringBuilder() StringBuilder {
	return &stringBuilder{
		b: &strings.Builder{},
	}
}

func (s stringBuilder) Write(x string) (int, error) { return s.b.WriteString(x) }
func (s stringBuilder) NewLine() error {
	_, err := s.b.WriteRune(newline)
	return err
}
func (s stringBuilder) Indent() error {
	_, err := s.Write(indent)
	return err
}
func (s stringBuilder) IndentN(n int) error {
	for i := 0; i < n; i++ {
		if err := s.Indent(); err != nil {
			return err
		}
	}
	return nil
}
func (s stringBuilder) String() string                  { return s.b.String() }
func (s stringBuilder) WriteLine(x string) (int, error) { return s.b.WriteString(x + string(newline)) }
