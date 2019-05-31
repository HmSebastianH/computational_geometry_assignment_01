package main

import (
	"fmt"
	"github.com/ross-oreto/go-tree"
)

type SweepEvent struct {
	x float64
	y float64
}

type IntersectionEvent struct {
	*SweepEvent
	lineA, lineB Line
}

type LineStartEvent struct {
	*SweepEvent
	line Line
}

type LineEndEvent struct {
	*SweepEvent
	line Line
}

func CompFloats(a, b float64) int8 {
	if a > b {
		return 1
	} else if a < b {
		return -1
	} else {
		return 0
	}
}

type SweepEntry struct {
	x float64
	y float64
	line Line
}

func (i SweepEntry) Comp(val tree.Val) int8 {
	v := val.(SweepEntry)
	r := CompFloats(i.y, v.y)
	if r == 0 {
		return CompFloats(i.x, v.x)
	}
	return r
}

func (i SweepEvent) Comp(val tree.Val) int8 {
	v := val.(SweepEvent)
	r := CompFloats(i.x, v.x)
	if r == 0 {
		return CompFloats(i.y, v.y)
	}
	return r
}


func LineSweep(allLines []Line) {
	// Assumptions about the data:
	// x-Koordinaten der Schnitt- und Endpunkte sind paarweise
	// verschieden
	// • Länge der Segmente > 0
	// • nur echte Schnittpunkte
	// • keine Linien parallel zur y-Achse
	// • keine Mehrfachschnittpunkte
	// • keine überlappenden Segmente


	eventQueue := tree.New()
	sweepLine := tree.New()

	for _, line := range allLines {
		eventQueue.Insert(LineStartEvent{&SweepEvent{line.start.x}, line})
		eventQueue.Insert(LineEndEvent{&SweepEvent{line.end.x}, line})
	}

	currentEvent := eventQueue.Head()
	for currentEvent != nil {
		eventQueue.Delete(currentEvent)
		currentEvent := eventQueue.Head()
	}

	eventQueue.Ascend(handleSweepEvent)
	fmt.Println(eventQueue.Values())
}
