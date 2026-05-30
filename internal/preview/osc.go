package preview

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Ryannnkl/ghostty-theme-tui/internal/ghostty"
)

const (
	esc = "\x1b"
	bel = "\x07"
)

type Previewer struct {
	out io.Writer
}

func New() Previewer {
	return Previewer{out: os.Stdout}
}

func NewWithWriter(out io.Writer) Previewer {
	return Previewer{out: out}
}

func (p Previewer) Apply(colors ghostty.ThemeColors) error {
	if colors.Foreground != "" {
		if _, err := fmt.Fprintf(p.out, "%s]10;%s%s", esc, colors.Foreground, bel); err != nil {
			return err
		}
	}
	if colors.Background != "" {
		if _, err := fmt.Fprintf(p.out, "%s]11;%s%s", esc, colors.Background, bel); err != nil {
			return err
		}
	}
	if colors.Cursor != "" {
		if _, err := fmt.Fprintf(p.out, "%s]12;%s%s", esc, colors.Cursor, bel); err != nil {
			return err
		}
	}

	indexes := make([]int, 0, len(colors.Palette))
	for index := range colors.Palette {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)

	for _, index := range indexes {
		if _, err := fmt.Fprintf(p.out, "%s]4;%d;%s%s", esc, index, colors.Palette[index], bel); err != nil {
			return err
		}
	}

	return nil
}
