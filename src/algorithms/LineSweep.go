package algorithms

import (
	"events"
	"fmt"
	. "geometry"
	. "sweepLine"
)



func LineSweep(allLines []*Line) []MatchingIndices {
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
		if line.Start.X != line.End.X {
			eventQueue.Insert(events.NewLineStartEvent(*line))
			eventQueue.Insert(events.NewLineEndEvent(*line))
		} else {
			fmt.Print(".")
		}
		// TODO: it might be cheaper to insert the end event at the insert event,
		// because at that point the event tree will probably be smaller
	}
	fmt.Println();


	// TODO Next: Comp von events anpassen:
	//  - intersection points an der gleichen stelle müssen nach einander kommen
	//  - Bei anderen events soll line id für eindeutige ordnung verwendet werden
	//  -> Delete anpassen das es die richtige linie findet (CCW auf linien end punkt, sobald ccw = 0 id suche)

	// TODO: This panics
	if !eventQueue.AssertOrder() {
		panic("Sanity check failed")
	}
	allIntersections := make([]MatchingIndices, 0)

	sweepLine := NewSweepLine()
	currentEvent := eventQueue.Pop()
	for currentEvent != nil {
		// Handle the different events:
		switch currentEvent.(type) {
			case *events.LineStartEvent:
				event := currentEvent.(*events.LineStartEvent)
				insertedNode := sweepLine.Insert(event.Line)
				leftNode := insertedNode.Left()
				// TODO: Make this a loop, to catch potential multiple intersections
				// Loop condition should be that the ccw changed
				if leftNode != nil {
					var intersection Point
					if event.Line.GetIntersectionWith(leftNode.Value, &intersection) {
						// TODO: If intersection is a overlap, do not insert a Event
						eventQueue.Insert(events.NewIntersectionEvent(intersection, event.Line, leftNode.Value))
					}
				}
				rightNode := insertedNode.Left()
				if rightNode != nil {
					var intersection Point
					if event.Line.GetIntersectionWith(rightNode.Value, &intersection) {
						eventQueue.Insert(events.NewIntersectionEvent(intersection, event.Line, rightNode.Value))
					}
				}

			case *events.LineEndEvent:
				// TODO:
				//   1. Find the node to be deleted
				//   2. Do same checks as on insertion with left / right neighbrors
				//   3. delete node
				// Get -> Right neighbor, check against left neighbor
				event := currentEvent.(*events.LineEndEvent)
				lineNode := sweepLine.FindWithEndPoint(event.Line)

				leftNode := lineNode.Left()
				rightNode := lineNode.Right()
				if leftNode != nil && rightNode != nil {
					var intersection Point
					if leftNode.Value.GetIntersectionWith(rightNode.Value, &intersection) {
						eventQueue.Insert(events.NewIntersectionEvent(intersection, leftNode.Value, rightNode.Value))
					}
				}

				sweepLine.Delete(lineNode)
			case *events.VerticalLineEvent:
				event := currentEvent.(*events.VerticalLineEvent)
				var _ = event
				// TODO:
				//  1. Collect all Vertical Line events together
				//  2. Sort them by y to compare to each other (look for overlaps)
				//  3. Check if ccw of start and end are different for any lines in the sweep line
				//  ! Do not add any of these lines to the sweep line

			case *events.IntersectionEvent:
				event := currentEvent.(*events.IntersectionEvent)
				var _ = event
				allIntersections =
					append(allIntersections, MatchingIndices{event.LineA.Index, event.LineB.Index})
				// TODO:
				//  1. Find all Intersection events on the same spot
				//  2. Reverse the order of all lines affected


			default:
				panic("Unknown event")
		}

		currentEvent = eventQueue.Pop()
	}

	fmt.Println("Done. Intersects: ", len(allIntersections))

	return allIntersections
}
