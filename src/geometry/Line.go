package geometry

import (
	"fmt"
	"math"
)

type Line struct {
	Index int
	Start Point
	End   Point
}

func NewLine(index int, p, q Point) *Line {
	obj := new(Line)
	obj.Index = index
	if (p.X == q.X) {
		// Vertical lines should go from top to bottom
		if p.Y < q.Y {
			obj.Start, obj.End = q, p
		} else {
			obj.Start, obj.End = p, q
		}

	} else if (p.X <= q.X) {
		obj.Start, obj.End = p, q
	} else {
		obj.Start, obj.End = q, p
	}
	return obj
}

func (p *Line) String() string {
	return fmt.Sprintf("%s->%s[%d]", p.Start.String(), p.End.String(), p.Index)
}

func (p *Line) IsCrossedBy(q Line) bool {
	if p.isPoint() {
		if q.isPoint() {
			return p.Start.equals(q.Start)
		} else {
			return q.hasPoint(p.Start)
		}
	}
	if q.isPoint() {
		return p.hasPoint(q.Start)
	}

	// At this point we are sure that neither line is a point
	// a comparision of the Ccw value therefore should be enough
	// to determine crossing of lines

	// First check that the points of q are at different sites of p
	ccwA := Ccw(*p, q.Start)
	ccwB := Ccw(*p, q.End)

	if ccwA*ccwB > 0 {
		return false
	}

	if ccwA == 0 && ccwB == 0 {
		return p.hasOverlapWith(q)
	}

	// The lines or not on one "pane", check the perspective of the other line
	ccwA = Ccw(q, p.Start)
	ccwB = Ccw(q, p.End)

	return ccwA * ccwB <= 0
}

func (p *Line) GetIntersectionWith(q Line) (*Point, bool) {
	// "Stolen" from Stackoverflow :)
	// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect

	//char get_line_intersection(float p0_x, float p0_y, float p1_x, float p1_y,
	//	float p2_x, float p2_y, float p3_x, float p3_y, float *i_x, float *i_y)
	//{
	s0_x := p.End.X - p.Start.X;
	s0_y := p.End.Y - p.Start.Y;
	s1_x := q.End.X - q.Start.X;
	s1_y := q.End.Y - q.Start.Y;

	s := (-s0_y * (p.Start.X - q.Start.X) + s0_x * (p.Start.Y - q.Start.Y)) / (-s1_x * s0_y + s0_x * s1_y);
	t := ( s1_x * (p.Start.Y - q.Start.Y) - s1_y * (p.Start.X - q.Start.X)) / (-s1_x * s0_y + s0_x * s1_y);

	if s >= 0 && s <= 1 && t >= 0 && t <= 1 {
		// Collision detected
		//if (i_x != NULL)
		//*i_x = p0_x + (t * s0_x);
		//if (i_y != NULL)
		//*i_y = p0_y + (t * s0_y);
		intersection := Point{p.Start.X + (t * s0_x), p.Start.Y + (t * s0_y)}
		return &intersection, true
	}
	if (math.IsNaN(s) || math.IsNaN(t)) && q.IsCrossedBy(*p) {
		return nil, true
	}

	return nil, false // No collision
}

func (a *Line) isPoint () bool {
	return (a.Start.X == a.End.X) && (a.Start.Y == a.End.Y)
}

func (a *Line) hasPoint (r Point) bool {
	if Ccw(*a, r) != 0 {
		return false
	}

	// The point is on one pane with the line
	x0, x1, x2 := a.Start.X, a.End.X, r.X
	if x0 == x1 {
		// This is a non monotonic line, use Y instead
		x0, x1, x2 = a.Start.Y, a.End.Y, r.Y
	}
	if x1 < x0 {
		// Swap them
		x0, x1 = x1, x0
	}

	return (x0 <= x2) &&  (x2 <= x1)

}

func (a *Line) hasOverlapWith( b Line) bool {
	// Check if line b starts or ends in line a
	if a.hasPoint(b.Start) || a.hasPoint(b.End) {
		return true
	}

	// Check if line a starts or ends in line b
	return b.hasPoint(a.Start) || b.hasPoint(a.End)
}

// CCW checks the position of point r relative to the line formed by p and q.
// The result will be negative if point r is to the right of p->q.
// The result will be positive if point r is to the left of p->q.
// The result will be 0 if point r is on the infinite line created by p->q
func Ccw(a Line, r Point) float64 {
	p := a.Start
	q := a.End
	return p.Y*r.X - q.Y*r.X + q.X*r.Y - p.X*r.Y - p.Y*q.X + p.X*q.Y
}


type MatchingIndices struct {
	IndexA, IndexB int
}

func NewMatchingIndices(indexA, indexB int) *MatchingIndices {
	obj := new(MatchingIndices)
	if (indexA <= indexB) {
		obj.IndexA, obj.IndexB = indexA, indexB
	} else {
		obj.IndexA, obj.IndexB = indexB, indexA
	}
	return obj
}