package roughfa_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/berquerant/roughfa"
	"github.com/stretchr/testify/assert"
)

func configureRenderTestcase(t *testing.T) {
	if testing.Short() {
		t.Skip("this test is skipped in short mode.")
	}
	if _, err := exec.LookPath(os.Getenv("DOTCOMMAND")); err != nil {
		t.Errorf("Please fix envvar DOTCOMMAND: %v", err)
	}
	t.Parallel()
}

type rendererRenderTestcase struct {
	name     string
	source   string
	filename string
	err      error
}

func (s rendererRenderTestcase) test(t *testing.T) {
	configureRenderTestcase(t)
	p := os.Getenv("PROJECT")
	filename := fmt.Sprintf("%s/tmp/%s", p, s.filename)
	t.Logf("render to %s", filename)
	err := roughfa.NewRenderer().
		DotCommand(os.Getenv("DOTCOMMAND")).
		Source(s.source).
		Filename(filename).
		Render()
	assert.Equal(t, s.err, err)
}

func TestRendererRender(t *testing.T) {
	for _, tc := range []*rendererRenderTestcase{
		{
			name:     "no dot source",
			filename: "no-dot.png",
			err:      roughfa.ErrNoDotSource,
		},
		{
			name:     "small graph",
			filename: "small-graph.png",
			source: `digraph {
  node [shape=circle]
  A -> B [label="ab"]
  B -> C [label="bc"]
  C -> A [label="ca"]
}`,
		},
	} {
		t.Run(tc.name, tc.test)
	}
}
