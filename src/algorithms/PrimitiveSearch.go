package algorithms

import (
	. "geometry"
	"sort"
	"sync"
)

func PrimitiveSearch(data []*Line) []*MatchingIndices {

	// Sort the data by x coordinates initially
	sort.Slice(data, func(i, j int) bool {
		return data[i].Start.X < data[j].Start.X
	})

	// Create a max size chanel buffer for concurrency
	ch := make(chan *MatchingIndices, 100000)
	wg := sync.WaitGroup{}

	for iLineP := 0; iLineP < len(data); iLineP++ {

		wg.Add(1)
		go findOverlapsForLine(iLineP, &data, ch, &wg)
	}

	wg.Wait()
	close(ch)
	var results []*MatchingIndices
	for match := range ch {
		results = append(results, match)
	}
	return results
}

func findOverlapsForLine(lineIndex int, allLines *[]*Line,ch chan *MatchingIndices, wg *sync.WaitGroup) {
	defer wg.Done()
	lineP := (*allLines)[lineIndex]
	for iLineQ := lineIndex+1; iLineQ < len(*allLines); iLineQ++ {

		lineQ := (*allLines)[iLineQ]
		if lineQ.Start.X > lineP.End.X {
			// The other lines starts after this one ends, no possible overlap
			return
		}

		if lineP.IsCrossedBy(*lineQ) {
			ch <- NewMatchingIndices(lineP.Index, lineQ.Index)
		}
	}
}