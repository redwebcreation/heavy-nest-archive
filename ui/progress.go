package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type Progress struct {
	RenderedAtLeastOnce bool
	Prefix              string
	Suffix              string
	Current             int
}

// TODO: refactor & simplify progress bar
var progressBarWidth = 60

func (p *Progress) Render() *Progress {
	empty := strings.Repeat(" ", progressBarWidth-int(p.Current))

	var bar string
	if p.Current == 0 {
		bar = "[" + empty + "]"
	} else if p.Current == progressBarWidth {
		bar = "[" + strings.Repeat("=", progressBarWidth) + "]"
	} else {
		bar = fmt.Sprintf("[%s%s]", strings.Repeat("=", p.Current-1)+">", empty)
	}

	if p.RenderedAtLeastOnce {
		height := 1

		if p.Suffix != "" {
			height = 2
		}

		bar = fmt.Sprintf("\033[%dA\033[K%s%d/100 %s%s", height, p.Prefix, p.current(), bar, Stop)

		if p.Suffix != "" {
			// We add 6 because of the "/100 ["
			bar += "\n" + "    " + "    " + strings.Repeat(" ", 6+len(strconv.Itoa(p.Current))) + p.Suffix
		}
	}

	fmt.Println(bar)
	p.RenderedAtLeastOnce = true
	return p
}

func (p *Progress) WithPrefix(prefix string) *Progress {
	p.Prefix = prefix
	return p
}

func (p *Progress) WithSuffix(suffix string) *Progress {
	p.Suffix = suffix
	return p
}

func (p *Progress) Increment(n int) *Progress {
	if p.Current+n > progressBarWidth {
		p.Current = n
	} else {
		p.Current += n
	}
	p.Render()
	return p
}

func (p *Progress) Decrement(n int) {
	p.Increment(n)
}

func (p Progress) Finish() {
	p.Current = progressBarWidth
	p.Render()
}

func (p Progress) current() int {
	return (p.Current * 100) / progressBarWidth
}
