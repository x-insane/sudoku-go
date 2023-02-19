package main

type Point struct {
	X, Y   int
	Value  int
	Acc    float64
	CanUse []int
}

type Step struct {
	X, Y  int
	Value int
	Acc   float64
}

type Steps struct {
	steps   []Step
	usedCnt *[10]int
}

func NewSteps(steps []Step, usedCnt *[10]int) Steps {
	return Steps{
		steps:   steps,
		usedCnt: usedCnt,
	}
}

func (s Steps) Len() int {
	return len(s.steps)
}

func (s Steps) Less(i, j int) bool {
	if s.steps[i].Acc == s.steps[j].Acc {
		// when acc equal pick the value that was used more times
		return s.usedCnt[s.steps[i].Value] > s.usedCnt[s.steps[j].Value]
	}
	return s.steps[i].Acc > s.steps[j].Acc
}

func (s Steps) Swap(i, j int) {
	s.steps[i].Acc, s.steps[j].Acc = s.steps[j].Acc, s.steps[i].Acc
}
