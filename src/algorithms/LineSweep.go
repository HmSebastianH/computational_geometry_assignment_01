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
		} else {
			eventQueue.Insert(events.NewVerticalLineEvent(*line))
		}
	}

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
			// Add the line end event now too, this is done later for performance reasons
			eventQueue.Insert(events.NewLineEndEvent(event.Line))
			insertedNode := sweepLine.Insert(event.Line)
			intersecs := checkNeighboringIntersections(insertedNode, event.GetX(), eventQueue, func(n *Node) *Node {
				return n.Left()
			})
			allIntersections = append(allIntersections, intersecs...)
			intersecs = checkNeighboringIntersections(insertedNode, event.GetX(), eventQueue, func(n *Node) *Node {
				return n.Right()
			})
			allIntersections = append(allIntersections, intersecs...)
		case *events.LineEndEvent:
			// Get -> Right neighbor, check against left neighbor
			event := currentEvent.(*events.LineEndEvent)
			lineNode := sweepLine.FindWithReferencePoint(event.Line.Index, event.Line.End)
			if lineNode == nil {
				panic("Invalid sweep line state")
			}

			// Gather all nodes to the left of the deleted node which have the same endpoint ccw (read: same direction)
			// and all nodes to the right of the deleted node which have the same endpoint ccw
			// Check for intersections in the opposing direction for all those nodes.
			// TODO: delete the node before doing these checks
			// TODO: for complete coverage the check method should be used to cover more casses
			leftNode := lineNode.Left()
			rightNode := lineNode.Right()
			if leftNode != nil && rightNode != nil {
				intersection, isIntersec := leftNode.Value.GetIntersectionWith(rightNode.Value)
				if isIntersec {
					if intersection != nil && intersection.X > event.GetX() {
						eventQueue.Insert(events.NewIntersectionEvent(*intersection, leftNode.Value, rightNode.Value))
					} else {
						// Just add store the intersection but no event for overlaps
						allIntersections = append(allIntersections, *NewMatchingIndices(leftNode.Value.Index, rightNode.Value.Index))
					}
				}
			}

			sweepLine.Delete(lineNode)
		case *events.VerticalLineEvent:
			event := currentEvent.(*events.VerticalLineEvent)
			// TODO:
			//  1. Collect all Vertical Line events together
			//  2. Sort them by y to compare to each other (look for overlaps)
			//  3. Check if ccw of start and end are different for any lines in the sweep line
			//  ! Do not add any of these lines to the sweep line
			nextEvent := eventQueue.Head()
			nVerticalLines := []Line{event.Line}
			for nextEvent != nil {
				additionalEvent, ok := nextEvent.Value.(*events.VerticalLineEvent)
				if ok && additionalEvent.GetX() == event.GetX() {
					// Gather all Vertical Line events on the same X-Line
					nVerticalLines = append(nVerticalLines, additionalEvent.Line)
					eventQueue.Pop()
					nextEvent = eventQueue.Head()
				} else {
					break
				}
			}

			for iLine, vLine := range nVerticalLines {
				// First check the line against the remaining other vertical lines
				for _, other := range nVerticalLines[iLine+1:] {
					if vLine.IsCrossedBy(other) {
						allIntersections = append(allIntersections, *NewMatchingIndices(vLine.Index, other.Index))
					}
				}
				// Then search for vertical overlaps
				verticalMatches := sweepLine.FindVerticalIntersections(vLine)
				allIntersections = append(allIntersections, verticalMatches...)
			}

		case *events.IntersectionEvent:
			event := currentEvent.(*events.IntersectionEvent)
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
			affectedNodes := reverseLineOrder(event.Intersection, involvedIds, sweepLine)

			if affectedNodes != nil && len(affectedNodes) > 0 {
				// TODO: The checks might have to be repeated for lines with the same ccw as the outer nodes
				leftMostNode := affectedNodes[0]
				intersecs := checkNeighboringIntersections(leftMostNode, event.GetX(), eventQueue, func(n *Node) *Node {
					return n.Left()
				})
				allIntersections = append(allIntersections, intersecs...)
				rightMostNode := affectedNodes[len(affectedNodes)-1]
				intersecs = checkNeighboringIntersections(rightMostNode, event.GetX(), eventQueue, func(n *Node) *Node {
					return n.Right()
				})
				allIntersections = append(allIntersections, intersecs...)
			}
		default:
			panic("Unknown event")
		}

		sweepLine.PrintOut()
		eventQueue.PrintOut()
		currentEvent = eventQueue.Pop()
	}

	// Overlapping lines might be detected as intersections multiple times
	allIntersections = filterDuplicates(allIntersections)
	fmt.Println("Done. Intersects: ", len(allIntersections))

	return allIntersections
}

func filterDuplicates(allIntersecs []MatchingIndices) []MatchingIndices {
	keys := make(map[MatchingIndices]struct{})
	list := make([]MatchingIndices, 0, len(allIntersecs))
	for _, entry := range allIntersecs {
		if _, value := keys[entry]; !value {
			keys[entry] = struct{}{}
			list = append(list, entry)
		}
	}
	return list
}

// reverseLineOrder is used to handle the swapping of lines in the sweep line after a intersection occured.
// It receives a set of all affected line ids, which can then be collected from the sweep line itself.
// The nodes which contain those affected lines will then be reversed in their order.
// To calculate the desired order the ccw of the line endpoint is used.
func reverseLineOrder(intersection Point, lineIds map[int]struct{}, sweepLine *SweepLine) []*Node {
	// Assumptions made about the sweep line state:
	// - All Ids are contained in the sweep line
	// - All Ids contained in the lineIds are neighbouring each other
	// - The affected lines are in the correct order
	// If theses assumptions are not met the function will fail
	if len(lineIds) < 2 {
		// Nothing to do here here
		return nil
	}

	idsToFind := make(map[int]struct{}, len(lineIds))
	for id := range lineIds {
		idsToFind[id] = struct{}{}
	}

	// Search for one of the lines using the lineId and intersection
	referenceId := -1
	for lineId := range idsToFind {
		// Get a arbitrary id from the contained ids
		referenceId = lineId
		break
	}
	// TODO: This function might return nil currently
	referenceNode := sweepLine.FindWithReferencePoint(referenceId, intersection)
	if referenceNode == nil {
		referenceNode = sweepLine.FindWithReferencePoint(referenceId, intersection)
		panic("Invalid SweepLine state")
	}
	// Start building a slice off ordered nodes
	affectedNodes := make([]*Node, 0, len(idsToFind))
	affectedNodes = append(affectedNodes, referenceNode)
	delete(idsToFind, referenceId)

	// Add all items to the "left" of the reference Node
	leftNeighbor := referenceNode.Left()
	for leftNeighbor != nil {
		_, ok := idsToFind[leftNeighbor.Value.Index]
		if !ok { // Not one of the affected line group anymore
			break
		}
		// We found one more affected line, prepend it
		affectedNodes = append([]*Node{leftNeighbor}, affectedNodes...)
		delete(idsToFind, leftNeighbor.Value.Index)
		leftNeighbor = leftNeighbor.Left()
	}

	// Repeat the same with right neighbors
	rightNeighbor := referenceNode.Right()
	for rightNeighbor != nil {
		_, ok := idsToFind[rightNeighbor.Value.Index]
		if !ok { // Not one of the affected line group anymore
			break
		}
		affectedNodes = append(affectedNodes, rightNeighbor)
		delete(idsToFind, rightNeighbor.Value.Index)
		rightNeighbor = rightNeighbor.Right()
	}

	// Make sure all assumptions were met
	if len(idsToFind) != 0 {
		panic("Bad Line Sweep structure")
	}

	newOrder := orderLinesByEndPoint(affectedNodes)
	for i, l := range newOrder {
		n := affectedNodes[i]
		n.Value = l
	}
	return affectedNodes
}

// This function is called for a node which was just inserted into the sweep line.
// It is checked for intersections against:
// - Lines directly above
// - Lines above lines with intersections (or above a ccw which had intersections)
// - Lines with the same endpoint ccw as the lines described above
// The same checks are done for lines directly below the inserted line
func checkNeighboringIntersections(n *Node, xThresh float64, eventQueue *events.EventQueue, iterFunc func(*Node) *Node) []MatchingIndices {
	leftN := iterFunc(n)
	ccwHadIntersec := true
	previousEndPoint := n.Value.End

	allIntersections := make([]MatchingIndices, 0)
	for leftN != nil {
		currentCcw := Ccw(leftN.Value, previousEndPoint)
		if currentCcw != 0 {
			if ccwHadIntersec {
				// We operate on a new ccw now
				ccwHadIntersec = false
			} else {
				// We found a ccw which does not contain any intersections
				break
			}
		}
		// Do the actual intersection check
		intersection, isIntersec := leftN.Value.GetIntersectionWith(n.Value)
		if isIntersec {
			if intersection != nil && intersection.X > xThresh {
				eventQueue.Insert(events.NewIntersectionEvent(*intersection, leftN.Value, n.Value))
			} else {
				// Just add store the intersection but no event for overlaps
				allIntersections = append(allIntersections, *NewMatchingIndices(leftN.Value.Index, n.Value.Index))
			}
			ccwHadIntersec = true
		}
		// Prep data for next iteration
		previousEndPoint = leftN.Value.End
		leftN = iterFunc(leftN)
	}
	return allIntersections
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
