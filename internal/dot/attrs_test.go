package dot_test

import (
	"testing"

	"github.com/berquerant/roughfa/internal/dot"
	"github.com/stretchr/testify/assert"
)

type attrTestcase struct {
	cname    string
	name     string
	value    string
	attrType dot.AttrType
	want     string
	err      error
}

func (s attrTestcase) test(t *testing.T) {
	a, err := dot.NewAttrBuilder().
		Name(s.name).
		Value(s.value).
		AttrType(s.attrType).
		Build()
	assert.Equal(t, s.err, err)
	if err != nil {
		return
	}
	assert.Equal(t, s.want, a.AsDot())
}

func TestAttr(t *testing.T) {
	for _, tc := range []*attrTestcase{
		{
			cname: "empty name",
			err:   dot.ErrAttrNameEmpty,
		},
		{
			cname: "raw",
			name:  "color",
			value: "red",
			want:  `color="red"`,
		},
		{
			cname:    "wrapped",
			name:     "color",
			value:    "#000000",
			attrType: dot.WrappedAttrType,
			want:     `color="#000000"`,
		},
	} {
		t.Run(tc.cname, tc.test)
	}
}
