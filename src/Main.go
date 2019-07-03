package main

import (
	"algorithms"
	"bufio"
	"fmt"
	"geometry"
	_ "geometry"
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





func main() {
	fileName := "1000_10"
	filePath := fmt.Sprintf("data/s_%s.dat", fileName)

	fmt.Printf("Start reading: %s\n", filePath)
	startTime := time.Now()

	inputFile, err := os.Open(filePath)
	check(err)

	outputFile, err := os.Create(fmt.Sprintf("data/result_%s.dat", fileName))
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

	//fmt.Println("Start Line Sweep...")
	//results := algorithms.LineSweep(data)

	fmt.Println("Start Primitive Search...")
	results := algorithms.PrimitiveSearch(data)

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

	fmt.Println("Found ", len(results), " Intersections, Total Time passed: ", time.Since(startTime))
}