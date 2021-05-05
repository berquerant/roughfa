package roughfa

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/berquerant/roughfa/internal/dot"
)

type (
	// Renderer renders Dot.
	Renderer interface {
		// Source sets a source string.
		// Required.
		Source(source string) Renderer
		// Filename sets a filename to output.
		// Required.
		Filename(filename string) Renderer
		// DotCommand sets a command to render dot.
		// Default is dot.
		DotCommand(dotCommand string) Renderer
		// Render renders dot into filename.
		// Use Graphviz to render.
		Render() error
		// RenderWithContext does Render with context.
		RenderWithContext(ctx context.Context) error
	}

	renderer struct {
		source     string
		filename   string
		dotCommand string
	}
)

// NewRenderer creates a new Renderer.
func NewRenderer() Renderer {
	return &renderer{
		dotCommand: "dot",
	}
}

func (s *renderer) Source(source string) Renderer {
	s.source = source
	return s
}
func (s *renderer) Filename(filename string) Renderer {
	s.filename = filename
	return s
}
func (s *renderer) DotCommand(dotCommand string) Renderer {
	s.dotCommand = dotCommand
	return s
}
func (s renderer) Render() error { return s.RenderWithContext(context.Background()) }
func (s renderer) RenderWithContext(ctx context.Context) error {
	if s.source == "" {
		return dot.ErrNoDotSource
	}
	tmpfile, err := ioutil.TempFile("", "*.dot")
	if err != nil {
		return err
	}
	defer func() {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
	}()
	if _, err := tmpfile.WriteString(s.source); err != nil {
		return err
	}
	t := fmt.Sprintf("-T%s", s.target())
	cmd := exec.CommandContext(ctx, s.dotCommand, t, tmpfile.Name(), "-o", s.filename)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	errStr, _ := ioutil.ReadAll(stderr)
	if err := cmd.Wait(); err != nil {
		return &RenderError{
			Err:    err,
			Stderr: string(errStr),
		}
	}
	return nil
}
func (s renderer) target() string { return strings.TrimLeft(filepath.Ext(s.filename), ".") }
