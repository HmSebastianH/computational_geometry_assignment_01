package events

import . "geometry"

// In general every sweep event is associated with a x point at which it is triggered
type SweepEvent interface {
	GetX() float64
	getPriority() int8
}


type LineStartEvent struct {
	Line Line
}
func (t LineStartEvent) GetX() float64 {
	return t.Line.Start.X
}
func (t LineStartEvent) getPriority() int8 { return 1 }



type VerticalLineEvent struct {
	Line Line
}
func (t VerticalLineEvent) GetX() float64 {
	return t.Line.Start.X
}
func (t VerticalLineEvent) getPriority() int8 { return 2 }



type IntersectionEvent struct {
	Intersection Point
	LineA, LineB Line
}
func (t IntersectionEvent) GetX() float64 {
	return t.Intersection.X
}
func (t IntersectionEvent) getPriority() int8 { return 3 }



type LineEndEvent struct {
	line Line
}
func (t LineEndEvent) GetX() float64 {
	return t.line.End.X
}
func (t LineEndEvent) getPriority() int8 { return 4 }


func compFloats(a, b float64) int8 {
	if a > b {
		return 1
	} else if a < b {
		return -1
	} else {
		return 0
	}
}

func CompEvents(eventA, eventB SweepEvent) int8 {
	xComp := compFloats(eventA.GetX(), eventB.GetX())
	if xComp != 0 {
		return xComp
	}
	prioComp := eventB.getPriority() - eventA.getPriority()
	if prioComp != 0 {
		return prioComp / prioComp
	}

	return 0
}