package algorithms

import (
	"events"
	"fmt"
	. "geometry"
	"sort"
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
	eventQueue.PrintOut()
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
			for leftNode != nil {
				var intersection Point
				if event.Line.GetIntersectionWith(leftNode.Value, &intersection) && intersection.X >= event.GetX() {
					// TODO: Update the function to properly differentiate between overlaps / intersections
					eventQueue.Insert(events.NewIntersectionEvent(intersection, event.Line, leftNode.Value))
				} else {
					break
				}
				leftNode = leftNode.Left()
			}
			rightNode := insertedNode.Right()
			for rightNode != nil {
				var intersection Point
				if event.Line.GetIntersectionWith(rightNode.Value, &intersection) && intersection.X >= event.GetX() {
					// TODO: Update the function to properly differentiate between overlaps / intersections
					eventQueue.Insert(events.NewIntersectionEvent(intersection, event.Line, rightNode.Value))
				} else {
					break
				}
				rightNode = rightNode.Right()
			}

		case *events.LineEndEvent:
			// TODO:
			//   1. Find the node to be deleted
			//   2. Do same checks as on insertion with left / right neighbrors
			//   3. delete node
			// Get -> Right neighbor, check against left neighbor
			event := currentEvent.(*events.LineEndEvent)
			lineNode := sweepLine.FindWithReferencePoint(event.Line.Index, event.Line.End)

			// TODO: LineNode can be nil, it should never be nil because we have lines to delete
			leftNode := lineNode.Left()
			rightNode := lineNode.Right()
			if leftNode != nil && rightNode != nil {
				var intersection Point
				if leftNode.Value.IsCrossedBy(rightNode.Value) {
					println("Intersec!")
				}
				if leftNode.Value.GetIntersectionWith(rightNode.Value, &intersection) && intersection.X > event.GetX() {
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
			allIntersections = append(allIntersections, *NewMatchingIndices(event.LineA.Index, event.LineB.Index))

			involvedIds := make(map[int]struct{}) // alias a Set of ints
			involvedIds[event.LineA.Index] = struct{}{}
			involvedIds[event.LineB.Index] = struct{}{}

			nextEvent := eventQueue.Head()
			for nextEvent != nil {
				additionalEvent, ok := nextEvent.Value.(*events.IntersectionEvent)
				if ok && additionalEvent.Intersection == event.Intersection {
					// Because that s the best way to handle sets ....
					allIntersections = append(allIntersections,
						*NewMatchingIndices(additionalEvent.LineA.Index, additionalEvent.LineB.Index))
					involvedIds[additionalEvent.LineA.Index] = struct{}{}
					involvedIds[additionalEvent.LineB.Index] = struct{}{}
					eventQueue.Pop() // This Event will be "handled"
					nextEvent = eventQueue.Head()
				} else {
					break
				}
			}
			reverseLineOrder(event.Intersection, involvedIds, sweepLine)

		default:
			panic("Unknown event")
		}

		sweepLine.PrintOut()
		eventQueue.PrintOut()
		currentEvent = eventQueue.Pop()
	}

	fmt.Println("Done. Intersects: ", len(allIntersections))

	return allIntersections
}

// reverseLineOrder is used to handle the swapping of lines in the sweep line after a intersection occured.
// It receives a set of all affected line ids, which can then be collected from the sweep line itself.
// The nodes which contain those affected lines will then be reversed in their order.
// To calculate the desired order the ccw of the line endpoint is used.
func reverseLineOrder(intersection Point, lineIds map[int]struct{}, sweepLine *SweepLine) {
	// Assumptions made about the sweep line state:
	// - All Ids are contained in the sweep line
	// - All Ids contained in the lineIds are neighbouring each other
	// - The affected lines are in the correct order
	// If theses assumptions are not met the function will fail
	if len(lineIds) < 2 {
		// Nothing to do here here
		return
	}

	// Search for one of the lines using the lineId and intersection
	referenceId := -1
	for lineId := range lineIds {
		// Get a arbitrary id from the contained ids
		referenceId = lineId
		break
	}
	// TODO: This function might return nil currently
	referenceNode := sweepLine.FindWithReferencePoint(referenceId, intersection)
	if referenceNode == nil {
		panic("Invalid SweepLine state")
	}
	// Start building a slice off ordered nodes
	affectedNodes := make([]*Node, 0, len(lineIds))
	affectedNodes = append(affectedNodes, referenceNode)
	delete(lineIds, referenceId)

	// Add all items to the "left" of the reference Node
	leftNeighbor := referenceNode.Left()
	for leftNeighbor != nil {
		_, ok := lineIds[leftNeighbor.Value.Index]
		if !ok { // Not one of the affected line group anymore
			break
		}
		// We found one more affected line, prepend it
		affectedNodes = append([]*Node{leftNeighbor}, affectedNodes...)
		delete(lineIds, leftNeighbor.Value.Index)
		leftNeighbor = leftNeighbor.Left()
	}

	// Repeat the same with right neighbors
	rightNeighbor := referenceNode.Right()
	for rightNeighbor != nil {
		_, ok := lineIds[rightNeighbor.Value.Index]
		if !ok { // Not one of the affected line group anymore
			break
		}
		affectedNodes = append(affectedNodes, rightNeighbor)
		delete(lineIds, rightNeighbor.Value.Index)
		rightNeighbor = rightNeighbor.Right()
	}

	// Make sure all assumptions were met
	if len(lineIds) != 0 {
		panic("Bad Line Sweep structure")
	}

	newOrder := orderLinesByEndPoint(affectedNodes)
	fmt.Println(newOrder)
	for i,l := range newOrder {
		fmt.Println(i, l)
		n := affectedNodes[i]
		n.Value = l
	}
}

func orderLinesByEndPoint(data []*Node) []Line {
	result := make([]Line, 0, len(data))
	for _, n := range data {
		result = append(result, n.Value)
	}

	sort.Slice(result, func(i, j int) bool {
		ccw := Ccw(result[i], result[j].End)
		return ccw > 0
	})
	return result
}
