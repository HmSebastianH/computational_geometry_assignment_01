package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
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
	check(scanner.Err())

	fmt.Println("Time passed (Reading Data): ", time.Since(startTime))


	sort.Slice(data, func(i, j int) bool {
		return data[i].start.x < data[j].start.x
	})

	fmt.Println("Time passed (Sorting Data): ", time.Since(startTime))


	// TODO: make chanel buffering "save"
	ch := make(chan *MatchingIndices, 100000)
	wg := sync.WaitGroup{}

	for iLineP := 0; iLineP < len(data); iLineP++ {

		wg.Add(1)
		go findOverlapsForLine(iLineP, &data, ch, &wg)
		//lineP := data[iLineP]
		/*for iLineQ := iLineP+1; iLineQ < len(data); iLineQ++ {
			lineQ := data[iLineQ]
			if lineQ.start.x > lineP.end.x {
				// The other lines starts after this one ends, no possible overlap
				break
			}

			if lineP.isCrossedBy(*lineQ) {
				results = append(results, NewMatchingIndices(lineP.index, lineQ.index))
			}
		}*/
	}

	wg.Wait()
	close(ch)
	var results []*MatchingIndices
	for match := range ch {
		results = append(results, match)
	}

	fmt.Println("Time passed (Calculating Matches): ", time.Since(startTime))
	fmt.Println("Num crossed lines:  ", len(results));

	sort.Slice(results, func(i, j int) bool {
		if results[i].indexA == results[j].indexA {
			return results[i].indexB < results[j].indexB
		}
		return results[i].indexA < results[j].indexA
	})

	writer := bufio.NewWriter(outputFile)
	for _, result := range results {
		_, err = writer.WriteString(strconv.Itoa(result.indexA) + "_" + strconv.Itoa(result.indexB) + "\n")
		check(err)
	}
	check(writer.Flush())

	fmt.Println("Time passed: ", time.Since(startTime))
}

func findOverlapsForLine(lineIndex int, allLines *[]*Line,ch chan *MatchingIndices, wg *sync.WaitGroup) {
	defer wg.Done()
	lineP := (*allLines)[lineIndex]
	for iLineQ := lineIndex+1; iLineQ < len(*allLines); iLineQ++ {

		lineQ := (*allLines)[iLineQ]
		if lineQ.start.x > lineP.end.x {
			// The other lines starts after this one ends, no possible overlap
			return
		}

		if lineP.isCrossedBy(*lineQ) {
			ch <- NewMatchingIndices(lineP.index, lineQ.index)
		}
	}
}