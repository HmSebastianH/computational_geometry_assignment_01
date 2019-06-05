package main

import (
	"bufio"
	"fmt"
	"geometry"
	_ "geometry"
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





func main() {

	startTime := time.Now()
	inputFile, err := os.Open("data/s_1000_1.dat")
	check(err)

	outputFile, err := os.Create("data/result_1000_1.dat")
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

	var results []*MatchingIndices


	writer := bufio.NewWriter(outputFile)
	for _, result := range results {
		_, err = writer.WriteString(strconv.Itoa(result.indexA) + "_" + strconv.Itoa(result.indexB) + "\n")
		check(err)
	}
	check(writer.Flush())

	fmt.Println("Time passed: ", time.Since(startTime))
}