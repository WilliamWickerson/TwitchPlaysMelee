package AST

import (
	"sort"
)

const (
	MAXFRAMEWINDOW = 90;
)

type Duration interface {
	Frames() []int;
}

type duration struct {
	frames []int;
	next int;
}

type Frames struct {
	Start int;
	End int;
}

func NewDuration(frames []Frames) Duration {
	set := make(map[int]bool);
	for _,f := range frames {
		for i := max(f.Start,1); i <= min(f.End, MAXFRAMEWINDOW); i++ {
			set[i - 1] = true;
		}
	}
	arrFrames := make([]int, 0);
	for i,_ := range set {
		arrFrames = append(arrFrames, i);
	}
	sort.Ints(arrFrames);
	return duration{arrFrames, 0};
}

func (d duration) Frames() []int {
	return d.frames;
}

func min(a, b int) int {
	if a < b {
		return a;
	}
	return b;
}

func max(a, b int) int {
	if a > b {
		return a;
	}
	return b;
}
