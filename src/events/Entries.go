package events

import (
	"fmt"
	. "geometry"
)

// In general every sweep event is associated with a x point at which it is triggered
type SweepEvent interface {
	GetX() float64
	getPriority() int8
	CompareTo(e SweepEvent) int8
}


type LineStartEvent struct {
	Line Line
}
func (t LineStartEvent) String() string{
	return fmt.Sprintf("LineStart(%s)", t.Line.String())
}
func (t LineStartEvent) GetX() float64 {
	return t.Line.Start.X
}
func (t LineStartEvent) getPriority() int8 { return 1 }
func NewLineStartEvent(line Line) * LineStartEvent {
	return &LineStartEvent{line}
}
func (t LineStartEvent) CompareTo(e SweepEvent) int8 {
	defaultResult := defaultComp(t, e)
	if defaultResult != 0 {
		return defaultResult
	}
	lineStartEvent := e.(*LineStartEvent)
	indexDif := t.Line.Index - lineStartEvent.Line.Index

	return abs(indexDif)
}




type VerticalLineEvent struct {
	Line Line
}
func (t VerticalLineEvent) String() string{
	return fmt.Sprintf("VerticalLine(%s)", t.Line.String())
}
func (t VerticalLineEvent) GetX() float64 {
	return t.Line.Start.X
}
func (t VerticalLineEvent) getPriority() int8 { return 2 }
func (t VerticalLineEvent) CompareTo(e SweepEvent) int8 {
	defaultResult := defaultComp(t, e)
	if defaultResult != 0 {
		return defaultResult
	}
	verticalEvent := e.(*VerticalLineEvent)
	indexDif := t.Line.Index - verticalEvent.Line.Index

	return abs(indexDif)
}
func NewVerticalLineEvent(line Line) * VerticalLineEvent {
	return &VerticalLineEvent{line}
}



type IntersectionEvent struct {
	Intersection Point
	LineA, LineB Line
}
func (t IntersectionEvent) String() string{
	return fmt.Sprintf("IntersectionEvent(%s %d_%d)", t.Intersection.String(), t.LineA.Index, t.LineB.Index)
}
func (t IntersectionEvent) GetX() float64 {
	return t.Intersection.X
}
func (t IntersectionEvent) getPriority() int8 { return 3 }

func NewIntersectionEvent(intersection Point, lineA, lineB Line) * IntersectionEvent {
	return &IntersectionEvent{intersection, lineA, lineB}
}
func (t IntersectionEvent) CompareTo(e SweepEvent) int8 {
	defaultResult := defaultComp(t, e)
	if defaultResult != 0 {
		return defaultResult
	}
	intersectionEvent := e.(*IntersectionEvent)
	intersecA := t.Intersection
	intersecB := intersectionEvent.Intersection

	// X should be assumed to be equal already
	yComp := compFloats(intersecA.Y, intersecB.Y)
	if yComp != 0 {
		return yComp
	}

	indexA1, indexB1 := t.LineA.Index, t.LineB.Index
	indexA2, indexB2 := intersectionEvent.LineA.Index, intersectionEvent.LineB.Index
	if indexA1 < indexB1 {
		indexA1, indexB1 = indexB1, indexA1
	}

	if indexA2 < indexB2 {
		indexA2, indexB2 = indexB2, indexA2
	}

	indexDif := indexA1 - indexA2
	if indexDif != 0 {
		return abs(indexDif)
	}

	indexDif = indexB1 - indexB2
	return abs(indexDif)
}



type LineEndEvent struct {
	Line Line
}
func (t LineEndEvent) String() string{
	return fmt.Sprintf("LineEnd(%s)", t.Line.String())
}
func (t LineEndEvent) GetX() float64 {
	return t.Line.End.X
}
func (t LineEndEvent) getPriority() int8 { return 4 }
func NewLineEndEvent(line Line) * LineEndEvent {
	return &LineEndEvent{line}
}
func (t LineEndEvent) CompareTo(e SweepEvent) int8 {
	defaultResult := defaultComp(t, e)
	if defaultResult != 0 {
		return defaultResult
	}
	endEvent := e.(*LineEndEvent)
	indexDif := t.Line.Index - endEvent.Line.Index

	return abs(indexDif)
}


func abs(x int) int8 {
	if x < 0 {
		return -1
	} else if x > 0 {
		return 1
	} else {
		return 0
	}
}

func compFloats(a, b float64) int8 {
	if a > b {
		return 1
	} else if a < b {
		return -1
	} else {
		return 0
	}
}

func defaultComp(eventA, eventB SweepEvent) int8 {
	xComp := compFloats(eventA.GetX(), eventB.GetX())
	if xComp != 0 {
		return xComp
	}
	prioComp := eventA.getPriority() - eventB.getPriority()
	if prioComp < 0 {
		return -1
	} else if prioComp > 0 {
		return 1
	}

	return 0
}