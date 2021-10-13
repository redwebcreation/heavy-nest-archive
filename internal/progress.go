package internal

import (
	"fmt"
	"strings"
)

type Progress struct {
	RenderedAtLeastOnce bool
	Prefix              string
	Suffix              string
	Current             int
}

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
		bar = fmt.Sprintf("\033[1A\033[K%s%d/100 %s\033[0m", p.Prefix, (p.Current*100)/progressBarWidth, bar)

		if p.Suffix != "" {
			bar += " " + p.Suffix
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
