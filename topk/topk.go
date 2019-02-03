package topk

import (
	"github.com/stephen-soltesz/geoloc/model"
)

// SiteDistance is a thing.
type SiteDistance []*model.Site

func (h SiteDistance) Len() int {
	return len(h)
}
func (h SiteDistance) Less(i, j int) bool {
	return h[i].Distance < h[j].Distance
}
func (h SiteDistance) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *SiteDistance) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*model.Site))
}

func (h *SiteDistance) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *SiteDistance) MinDistance() float64 {
	return (*h)[0].Distance
}
