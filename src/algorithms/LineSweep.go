package algorithms

import (
	"events"
	. "geometry"
	. "sweepLine"
)



func LineSweep(allLines []Line) {
	// Assumptions about the data:
	// x-Koordinaten der Schnitt- und Endpunkte sind paarweise
	// verschieden
	// • Länge der Segmente > 0
	// • nur echte Schnittpunkte
	// • keine Linien parallel zur y-Achse
	// • keine Mehrfachschnittpunkte
	// • keine überlappenden Segmente

	// Extra handling of: Vertical lines, Points, Multiple Intersections (? merge intersection  events if in same point)
	// Vertical lines can be checked in tree with ccw too (Crosses all lines with different ccws on start and end point)


	// Create the event queue to work on
	eventQueue := events.NewEventQueue()

	// Fill it with all known start and end events
	for _, line := range allLines {
		eventQueue.Insert(events.LineStartEvent{line})
		eventQueue.Insert(events.LineEndEvent{line})
		// TODO: it might be cheaper to insert the end event at the insert event,
		// because at that point the event tree will probably be smaller
	}


	sweepLine := NewSweepLine()
	allIntersections := make([]MatchingIndices, 0)
	currentEvent := eventQueue.Pop()
	for currentEvent != nil {
		// Handle the different events:
		switch currentEvent.(type) {
			case events.LineStartEvent:
				event := currentEvent.(events.LineStartEvent)
				sweepLine.Insert(event.Line)

			case events.LineEndEvent:
				event := currentEvent.(events.LineEndEvent)
				var _ = event
				// Delete Elem
				// Get -> Right neighbor, check against left neighbor

			case events.VerticalLineEvent:
				event := currentEvent.(events.VerticalLineEvent)
				var _ = event

			case events.IntersectionEvent:
				event := currentEvent.(events.IntersectionEvent)
				var _ = event


		}

		currentEvent = eventQueue.Pop()
	}

	/*
	eventQueue.Ascend(handleSweepEvent)
	fmt.Println(eventQueue.Values())*/
}
