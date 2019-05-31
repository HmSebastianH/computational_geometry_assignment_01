package main

import (
	"bufio"
	"fmt"
	"os"
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

type MatchingIndices struct {
	indexA, indexB int
}

func NewMatchingIndices(indexA, indexB int) *MatchingIndices {
	obj := new(MatchingIndices)
	if (indexA <= indexB) {
		obj.indexA, obj.indexB = indexA, indexB
	} else {
		obj.indexA, obj.indexB = indexB, indexA
	}
	return obj
}


func main() {

	startTime := time.Now()
	inputFile, err := os.Open("data/s_1000_1.dat")
	check(err)

	outputFile, err := os.Create("data/result_1000_1.dat")
	check(err)

	defer inputFile.Close()
	defer outputFile.Close()

	var data []*Line
	currentLineIndex := 0

	scanner := bufio.NewScanner(inputFile)
	check(scanner.Err())

	for scanner.Scan() {
		var p0 float64
		var p1 float64
		var q0 float64
		var q1 float64


		_, err = fmt.Fscan(strings.NewReader(scanner.Text()), &p0, &p1, &q0, &q1)
		check(err)
		line := NewLine(currentLineIndex, Point{p0, p1}, Point{q0, q1})
		currentLineIndex++

		//line := Line{Point{p0, p1}, Point{q0, q1}}
		data = append(data, line)
	}

	fmt.Println("Time passed (Reading Data): ", time.Since(startTime))

	var results []*MatchingIndices


	writer := bufio.NewWriter(outputFile)
	for _, result := range results {
		_, err = writer.WriteString(strconv.Itoa(result.indexA) + "_" + strconv.Itoa(result.indexB) + "\n")
		check(err)
	}
	check(writer.Flush())

	fmt.Println("Time passed: ", time.Since(startTime))
}