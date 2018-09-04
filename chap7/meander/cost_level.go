package meander

import (
	"strings"
)

// int8: max number of variasions
type Cost int8

const (
	// iota starts from 0
	_ Cost = iota
	Cost1
	Cost2
	Cost3
	Cost4
	Cost5
)

type CostRange struct {
	From Cost
	To   Cost
}

var costStrings = map[string]Cost{
	"$":     Cost1,
	"$$":    Cost2,
	"$$$":   Cost3,
	"$$$$":  Cost4,
	"$$$$$": Cost5,
}

func (c Cost) String() string {
	for s, v := range costStrings {
		if c == v {
			return s
		}
	}
	return "invalid value"
}

func ParseCost(s string) Cost {
	return costStrings[s]
}

func (r CostRange) String() string {
	return r.From.String() + "..." + r.To.String()
}

func ParseCostRange(s string) *CostRange {
	segs := strings.Split(s, "...")
	return &CostRange{
		From: ParseCost(segs[0]),
		To:   ParseCost(segs[1]),
	}
}
