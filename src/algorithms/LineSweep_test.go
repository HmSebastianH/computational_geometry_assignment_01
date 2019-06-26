package algorithms

import (
	"bufio"
	"fmt"
	"geometry"
	. "geometry"
	_ "geometry"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)



func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestProfiling(t *testing.T) {
	startTime := time.Now()
	fileName := "10000_1"
	inputFile, err := os.Open(fmt.Sprintf("C:/Users/sebas/Desktop/hm/geometry/computational_geometry_assignment_01/data/s_%s.dat", fileName))
	check(err)

	outputFile, err := os.Create(fmt.Sprintf("C:/Users/sebas/Desktop/hm/geometry/computational_geometry_assignment_01/data/result_%s.dat", fileName))
	check(err)

	defer inputFile.Close()
	defer outputFile.Close()

	var data []*geometry.Line
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
		line := geometry.NewLine(currentLineIndex, geometry.Point{p0, p1}, geometry.Point{q0, q1})
		currentLineIndex++

		//line := Line{Point{p0, p1}, Point{q0, q1}}
		data = append(data, line)
	}

	fmt.Println("Time passed (Reading Data): ", time.Since(startTime))

	results := LineSweep(data)
	sort.Slice(results, func(i, j int) bool {
		if results[i].IndexA == results[j].IndexA {
			return results[i].IndexB < results[j].IndexB
		}
		return results[i].IndexA < results[j].IndexA
	})

	writer := bufio.NewWriter(outputFile)
	for _, result := range results {
		_, err = writer.WriteString(strconv.Itoa(result.IndexA) + "_" + strconv.Itoa(result.IndexB) + "\n")
		check(err)
	}
	check(writer.Flush())

	fmt.Println("Time passed: ", time.Since(startTime))
}

func TestLineSweepWithSimpleIntersection_1(t *testing.T) {
	lineId := 0
	var data []*Line
	data = appendLine(data, &lineId, 1, 0, 4, 2)
	data = appendLine(data, &lineId, 2, 2, 4, 0)

	intersections := LineSweep(data)
	checkResultSize(t, intersections, 1)
	checkIfResultContains(t, intersections, 0, 1)
}

func TestLineSweepWithNoIntersection_2(t *testing.T) {
	lineId := 0
	var data []*Line
	data = appendLine(data, &lineId, 1, 3, 3, 1)
	data = appendLine(data, &lineId, 2, 3, 4, 2)

	intersections := LineSweep(data)
	checkResultSize(t, intersections, 0)
}

func TestLineSweepWithTripleIntersection_3(t *testing.T) {
	lineId := 0
	var data []*Line
	data = appendLine(data, &lineId, 1, 2, 5, 2)
	data = appendLine(data, &lineId, 2, 1, 5, 4)
	data = appendLine(data, &lineId, 2, 3, 4, 1)
	data = appendLine(data, &lineId, 4, 4, 5, 3)

	intersections := LineSweep(data)
	checkResultSize(t, intersections, 4)
	checkIfResultContains(t, intersections, 0, 1)
	checkIfResultContains(t, intersections, 0, 2)
	checkIfResultContains(t, intersections, 1, 2)
	checkIfResultContains(t, intersections, 1, 3)
}

func TestLineSweepWithVerticalOverlap_4(t *testing.T) {
	lineId := 0
	var data []*Line
	data = appendLine(data, &lineId, 2, 1, 2, 3)
	data = appendLine(data, &lineId, 2, 2, 2, 4)

	intersections := LineSweep(data)
	checkResultSize(t, intersections, 1)
	checkIfResultContains(t, intersections, 0, 1)
}

func TestLineSweepWithPointOnLine_5(t *testing.T) {
	lineId := 0
	var data []*Line
	data = appendLine(data, &lineId, 2, 2, 5, 2)
	data = appendLine(data, &lineId, 4, 2, 4, 2)

	intersections := LineSweep(data)
	checkResultSize(t, intersections, 1)
	checkIfResultContains(t, intersections, 0, 1)
}

func TestLineSweepWithEndToStartIntersection_6(t *testing.T) {
	lineId := 0
	var data []*Line
	data = appendLine(data, &lineId, 2, 2, 4, 2)
	data = appendLine(data, &lineId, 4, 2, 5, 1)

	intersections := LineSweep(data)
	checkResultSize(t, intersections, 1)
	checkIfResultContains(t, intersections, 0, 1)
}

func TestLineSweepWithMultipleDifferentIntersections_7(t *testing.T) {
	lineId := 0
	var data []*Line
	data = appendLine(data, &lineId, 1, 1, 5, 3)
	data = appendLine(data, &lineId, 1, 2, 2, 1)
	data = appendLine(data, &lineId, 2, 2, 3, 1)
	data = appendLine(data, &lineId, 2, 3, 4, 1)

	intersections := LineSweep(data)
	checkResultSize(t, intersections, 3)
	checkIfResultContains(t, intersections, 0, 1)
	checkIfResultContains(t, intersections, 0, 2)
	checkIfResultContains(t, intersections, 0, 3)
}

func TestLineSweepWithOverLapAndIntersection_8(t *testing.T) {
	lineId := 0
	var data []*Line
	data = appendLine(data, &lineId, 1, 2, 4, 2)
	data = appendLine(data, &lineId, 2, 1, 4, 3)
	data = appendLine(data, &lineId, 2, 2, 6, 2)

	intersections := LineSweep(data)
	checkResultSize(t, intersections, 3)
	checkIfResultContains(t, intersections, 0, 1)
	checkIfResultContains(t, intersections, 0, 2)
	checkIfResultContains(t, intersections, 1, 2)
}

func TestLineSweepWithQuadIntersection_9(t *testing.T) {
	lineId := 0
	var data []*Line
	data = appendLine(data, &lineId, 1, 4, 3, 0)
	data = appendLine(data, &lineId, 1, 3, 3, 1)
	data = appendLine(data, &lineId, 1, 1, 3, 3)
	data = appendLine(data, &lineId, 1, 0, 3, 4)
	data = appendLine(data, &lineId, 2, 3.5, 4, 3.5)
	data = appendLine(data, &lineId, 2, 0.5, 4, 0.5)
	data = appendLine(data, &lineId, 3, 1, 3, 3 )

	intersections := LineSweep(data)
	checkResultSize(t, intersections, 10)
	checkIfResultContains(t, intersections, 0, 1)
	checkIfResultContains(t, intersections, 0, 2)
	checkIfResultContains(t, intersections, 0,3)
	checkIfResultContains(t, intersections, 0, 5)
	checkIfResultContains(t, intersections, 1, 2)
	checkIfResultContains(t, intersections, 1, 3)
	checkIfResultContains(t, intersections, 1, 6)
	checkIfResultContains(t, intersections, 2, 3)
	checkIfResultContains(t, intersections, 2, 6)
	checkIfResultContains(t, intersections, 3, 4)
}

func appendLine(data []*Line, id *int, x0, y0, x1, y1 float64) []*Line{
	data = append(data, NewLine(*id, Point{x0,y0}, Point{x1, y1}))
	*id = *id + 1
	return data
}

func checkResultSize(t *testing.T, result []MatchingIndices, expected int) {
	logResultContent(t, result)
	if len(result) != expected {
		t.Errorf("Unexpected ammount of intersections found, expected %d but was %d",
			expected, len(result))
	}
}

func checkIfResultContains(t *testing.T, result []MatchingIndices, indexA, indexB int) {
	for _,r := range result {
		if r.IndexA == indexA && r.IndexB == indexB {
			// We found the expected result
			return
		}
	}
	t.Errorf("Result slice did not contain the result %d_%d", indexA, indexB)
}

func logResultContent(t *testing.T, result []MatchingIndices) {
	for _,r := range result {
		t.Log("Result:", r.IndexA, r.IndexB)
	}
}