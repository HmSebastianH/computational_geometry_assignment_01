package main

type Line struct {
	start Point
	end   Point
}

func NewLine(p, q Point) *Line {
	obj := new(Line)
	if (p.x <= q.x) {
		obj.start, obj.end = p, q
	} else {
		obj.start, obj.end = q, p
	}
	return obj
}

func (p *Line) isCrossedBy(q Line) bool {
	if p.isPoint() {
		if q.isPoint() {
			return p.start.equals(q.start)
		} else {
			return q.hasPoint(p.start)
		}
	}
	if q.isPoint() {
		return p.hasPoint(q.start)
	}

	// At this point we are sure that neither line is a point
	// a comparision of the ccw value therefore should be enough
	// to determine crossing of lines

	// First check that the points of q are at different sites of p
	ccwA := ccw(*p, q.start)
	ccwB := ccw(*p, q.end)

	if ccwA*ccwB > 0 {
		return false
	}

	if ccwA == 0 && ccwB == 0 {
		return p.hasOverlapWith(q)
	}

	// The lines or not on one "pane", check the perspective of the other line
	ccwA = ccw(q, p.start)
	ccwB = ccw(q, p.end)

	return ccwA * ccwB <= 0
}

func (a *Line) isPoint () bool {
	return (a.start.x == a.end.x) && (a.start.y == a.end.y)
}

func (a *Line) hasPoint (r Point) bool {
	if ccw(*a, r) != 0 {
		return false
	}

	// The point is on one pane with the line
	x0, x1, x2 := a.start.x, a.end.x, r.x
	if x0 == x1 {
		// This is a non monotonic line, use y instead
		x0, x1, x2 = a.start.y, a.end.y, r.y
	}
	if x1 < x0 {
		// Swap them
		x0, x1 = x1, x0
	}

	return (x0 <= x2) &&  (x2 <= x1)

}

func (a *Line) hasOverlapWith( b Line) bool {
	// Check if line b starts or ends in line a
	if a.hasPoint(b.start) || a.hasPoint(b.end) {
		return true
	}

	// Check if line a starts or ends in line b
	return b.hasPoint(a.start) || b.hasPoint(a.end)
}