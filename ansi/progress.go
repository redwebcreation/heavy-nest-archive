package ansi

import (
	"fmt"
	"strings"
)

type Progress struct {
	Total                int
	Current              int
	Bar                  string
	Prefix               string
	Suffix               string
	Width                int
	Label                string
	hasBeenRenderedOnce  bool
	hasRenderedLabelOnce bool
}

func (p *Progress) String() string {
	if p.Total == 0 {
		return ""
	}

	percent := float64(p.Current) / float64(p.Total)
	filled := int(float64(p.Width) * percent)
	empty := p.Width - filled

	bar := p.Prefix + strings.Repeat(p.Bar, filled) + strings.Repeat(" ", empty) + p.Suffix + "\n"

	if p.Label != "" {
		return bar + p.Label + "\n"
	}

	return bar
}

func NewProgress(total int, width int) *Progress {
	return &Progress{
		Total:   total,
		Current: 0,
		Bar:     "=",
		Prefix:  "[",
		Suffix:  "]",
		Width:   width,
	}
}

func (p *Progress) Set(current int) *Progress {
	p.Current = current
	p.Render()
	return p
}

func (p *Progress) WithLabel(label string) *Progress {
	p.Label = label
	return p
}

func (p *Progress) Finish() *Progress {
	return p.Set(p.Total)
}

func (p *Progress) Increment() *Progress {
	return p.Set(p.Current + 1)
}

func (p *Progress) Render() {
	if !p.hasBeenRenderedOnce {
		fmt.Print(p.String())
		p.hasBeenRenderedOnce = true
		return
	}

	if !PrintAnsi {
		fmt.Print(p.String())
		return
	}
	if p.Label != "" {
		CursorUp(2)
	} else {
		CursorUp(1)
	}

	fmt.Print(p.String())
}

func (p *Progress) SetPrefix(s string) *Progress {
	p.Prefix = s
	return p
}
