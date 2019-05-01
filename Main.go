package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)



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


func main() {

	startTime := time.Now()
	inputFile, err := os.Open("data/s_100000_1.dat")
	check(err)

	outputFile, err := os.Create("data/result_100000_1.dat")
	check(err)

	defer inputFile.Close()
	defer outputFile.Close()

	var data []*Line

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		var p0 float64
		var p1 float64
		var q0 float64
		var q1 float64


		_, err = fmt.Fscan(strings.NewReader(scanner.Text()), &p0, &p1, &q0, &q1)
		check(err)
		line := NewLine(Point{p0, p1}, Point{q0, q1})
		//line := Line{Point{p0, p1}, Point{q0, q1}}
		data = append(data, line)
	}
	check(scanner.Err())

	fmt.Println("Time passed (Reading Data): ", time.Since(startTime))


	sort.Slice(data, func(i, j int) bool {
		return data[i].start.x < data[j].start.x
	})

	fmt.Println("Time passed (Sorting Data): ", time.Since(startTime))


	totalHits := 0
	writer := bufio.NewWriter(outputFile)
	for iLineP := 0; iLineP < len(data); iLineP++ {
		lineP := data[iLineP]

		for iLineQ := iLineP+1; iLineQ < len(data); iLineQ++ {
			if iLineQ == iLineP {
				continue // skip itself
			}

			lineQ := data[iLineQ]
			if lineQ.start.x > lineP.end.x {
				// The other lines starts after this one ends, no possible overlap
				break
			}

			if lineP.isCrossedBy(*lineQ) {
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
