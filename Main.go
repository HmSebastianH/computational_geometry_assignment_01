package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)


type Point struct {
	x float64
	y float64
}

func (a *Point) equals( b Point) bool{
	return (a.x == b.x) && (a.y == b.y)
}

type Line struct {
	start Point
	end   Point
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// CCW checks the position of point r relative to the line formed by p and q.
// The result will be negative if point r is to the right of p->q.
// The result will be positive if point r is to the left of p->q.
// The result will be 0 if point r is on the infinite line created by p->q
func ccw(a Line, r Point) float64 {
	p := a.start
	q := a.end
	return p.y*r.x - q.y*r.x + q.x*r.y - p.x*r.y - p.y*q.x + p.x*q.y
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

func main() {

	startTime := time.Now()
	inputFile, err := os.Open("data/s_1000_1.dat")
	check(err)

	outputFile, err := os.Create("data/result_1000_1.dat")
	check(err)

	defer inputFile.Close()
	defer outputFile.Close()

	var data []Line

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		var p0 float64
		var p1 float64
		var q0 float64
		var q1 float64

		_, err = fmt.Fscan(strings.NewReader(scanner.Text()), &p0, &p1, &q0, &q1)
		check(err)
		line := Line{Point{p0, p1}, Point{q0, q1}}
		data = append(data, line)
	}
	check(scanner.Err())

	totalHits := 0
	writer := bufio.NewWriter(outputFile)
	for iLineP := 0; iLineP < len(data); iLineP++ {
		lineP := data[iLineP]

		for iLineQ := iLineP+1; iLineQ < len(data); iLineQ++ {
			if iLineQ == iLineP {
				continue // skip itself
			}

			lineQ := data[iLineQ]
			if lineP.isCrossedBy(lineQ) {
				totalHits++
				_, err = writer.WriteString(strconv.Itoa(iLineP) + "_" + strconv.Itoa(iLineQ) + "\n")
				check(err)
			}
		}
	}

	check(writer.Flush())

	fmt.Println("Num crossed lines:  ", totalHits);
	fmt.Println("Time passed: ", time.Since(startTime))
}
