package main

import (
	"sort"
)

type (
	HitsPage struct {
		Page string
		Hit  int
	}

	AllHits []HitsPage
)

func (a AllHits) Len() int           { return len(a) }
func (a AllHits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AllHits) Less(i, j int) bool { return a[i].Hit < a[j].Hit }

func TopHits(mapUrls map[string]int, top int) AllHits {

	ah := make(AllHits, len(mapUrls))
	i := 0
	for url, hit := range mapUrls {
		ah[i] = HitsPage{Page: url, Hit: hit}
		i++
	}

	sort.Sort(sort.Reverse(ah))

	return ah[:top]

}
