package main

import (
	"sort"
	"strconv"
)

type (
	HitsPage struct {
		Page string
		Hit  int
	}

	AllHits []HitsPage

	Point struct {
		X int `json:"x"`
		Y int `json:"y"`
		R int `json:"r"`
	}
)

func (a AllHits) Len() int           { return len(a) }
func (a AllHits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AllHits) Less(i, j int) bool { return a[i].Hit < a[j].Hit }

func TopHits(mapUrls map[Url]int, top int) AllHits {

	ah := make(AllHits, len(mapUrls))
	i := 0
	for url, hit := range mapUrls {
		ah[i] = HitsPage{Page: string(url), Hit: hit}
		i++
	}

	sort.Sort(sort.Reverse(ah))

	return ah[:top]

}

func PreparePoints(mapDateTime map[Month]HoursHits) []Point {

	points := make([]Point, 0, len(mapDateTime)*24)

	for month, hourhits := range mapDateTime {
		y, _ := strconv.Atoi(string(month))
		for time, hits := range hourhits {

			r := hits * 15 / 1000
			if r == 0 {
				continue
			}

			x, _ := strconv.Atoi(string(time))
			points = append(points, Point{
				X: x,
				Y: y,
				R: r,
			})
		}
	}

	return points

}
