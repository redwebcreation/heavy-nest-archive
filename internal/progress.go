package internal

import (
	"fmt"
	"strings"
)

type Progress struct {
	RenderedOnce bool
	Label        string
	Current      int
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

	if p.RenderedOnce {
		bar = fmt.Sprintf("\033[1A\033[K%d/100 %s\033[0m", (p.Current*100)/progressBarWidth, bar)

		if p.Label != "" {
			bar += " " + p.Label[0:min(len(p.Label), TermWidth-len(bar))]
		}
	}

	fmt.Println(bar)
	p.RenderedOnce = true
	return p
}

func (p *Progress) WithLabel(label string) *Progress {
	p.Label = label
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

func min(x int, y int) int {
	if x < y {
		return x
	}

	return y
}
