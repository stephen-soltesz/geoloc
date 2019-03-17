package main

import (
	"container/heap"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/m-lab/go/rtx"

	"github.com/stephen-soltesz/geoloc/topk"

	"github.com/stephen-soltesz/geoloc/model"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/m-lab/go/flagx"
)

var (
	geofile flagx.FileBytes
	input   flagx.FileBytes
	Xcolors = []string{
		"#F44336", // red
		"#E91E63", // pink
		"#E040FB", // Purple
		"#7C4DFF", // Purple
		// "#536DFE", // blue
		"#2196F3", // blue
		"#00CBD4", // cyan
		// "#009688", // teal
		"#4CAF50", // green
		"#FFC107", // amber
		"#FF5722", // deep orange
	}
	Ycolors = []string{
		"#FF0029",
		"#00d2D5",
		// "#FF7F00",
		"#B3E900",
		"#F981BF",
		"#FB8072",
		"#FCCDE5",
		"#FFED6F",
		"#00D067",
		"#CDDC39", // lime
	}
	colors = []int{
		// "#ff80ab", //
		//	"#FF1744", // light red
		//"#d500f9", // bright purple
		// "#2979ff",
		//"#00b0ff", // blue
		//"#18ffff", // light blue.
		//"#00e676", // light green.
		//"#ffea00", // yellow.
		// ["#ff3d00", // red
		//"#c6ff00", // yellow-green.
		//"#00ff00", // green.

		/*
			"#FF1744", // red
			"#d500f9", // pink
			"#798dff", // purple
			"#40c4ff", // blue
			"#00e676", // green
			"#eeff41", // yellow
			"#ffab00", // orange
			"#efebe9", // white
		*/
		// "#000000", // black
		//"#FF0000", // red
		//"#FCA581", // orange
		//"#3FF52E", // green
		// "#0069FA", // blue
		//"#F951FC", // pink
		//"#9A81DF", // purple
		// "#835AFA", // purple
		//"#FFFF00", // yellow
		//"#00FFFF", // cyan
		//"#FFFFFF", // white
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
		/*
			"a",
			"b",
			"c",
			"d",
			"e",
			"f",
			"g",
			"h",
			"i",
			"j",
			"k",
			"l",
			"m",
			"n",
			"o",
			"p",
			"q",
			"r",
			"s",
			"t",
		*/

		//		"#68AFF9", // blue
		//	"#3EC300", // green
		//"#00FFFF", // light blue
		//"#FF00FF", // purple
		//"#FFFC31", // yellow
		//"#FFFFFF", // white
	}
)

func init() {
	flag.Var(&input, "file", "File to read contents.")
	log.SetOutput(os.Stderr)
}

func findClosestN(one *model.Site, sites []*model.Site, n int) []*model.Site {
	nearest := make([]*model.Site, 0, n)
	for i := range sites {
		site := sites[i]
		a := s2.LatLng{
			Lat: s1.Angle(site.Lat * math.Pi / 180.0),
			Lng: s1.Angle(site.Lon * math.Pi / 180.0),
		}
		b := s2.LatLng{
			Lat: s1.Angle(one.Lat * math.Pi / 180.0),
			Lng: s1.Angle(one.Lon * math.Pi / 180.0),
		}
		sites[i].Distance = (float64)(a.Distance(b))
	}
	h := &topk.SiteDistance{}
	heap.Init(h)
	for i := 0; i < n; i++ {
		heap.Push(h, sites[i])
	}
	for _, site := range sites[n:] {
		//if h.MinDistance() >= site.Distance {
		//heap.Pop(h)
		heap.Push(h, site)
		//}
	}
	// fmt.Println("len", h.Len())
	nearest = append(nearest, heap.Pop(h).(*model.Site))
	for i := 0; i < (n - 1); i++ {
		nearest = append(nearest, heap.Pop(h).(*model.Site))
	}
	// fmt.Println("len", h.Len())
	/*
		TODO: run benchmark to compare heap with simple sort.
		sort.Slice(*sites, func(i, j int) bool {
			return (*sites)[i].Distance < (*sites)[j].Distance
		})
	*/
	return nearest
}

func colorUniquely(metroColor map[string]int, metroNames []string) {
	haveColors := map[int]string{}
	needColors := []string{}
	// Find which metros already have colors, and which need new colors.
	for _, metro := range metroNames {
		if c, ok := metroColor[metro]; !ok {
			needColors = append(needColors, metro)
		} else {
			haveColors[c] = metro
		}
	}
	// For each metro needing a color, look for a unique color to assign.
	for _, metro := range needColors {
		for _, color := range colors {
			// Find first color that is not found in haveColors.
			_, mc := metroColor[metro]
			_, hc := haveColors[color]
			if !mc && !hc {
				haveColors[color] = metro
				metroColor[metro] = color
				break
			}
		}
	}
	return
}

func main() {
	flag.Parse()

	sites := model.LoadSitesInfo("mlab-site-stats.json")
	metros := map[string]*model.Site{}

	// Collect one site from each metro.
	// fmt.Println("Total:", len(sites))
	for _, site := range sites {
		if _, ok := metros[site.Metro[1]]; !ok {
			metros[site.Metro[1]] = site
		}
	}

	// Consolidate into a short list.
	fsites := []*model.Site{}
	for _, site := range metros {
		fsites = append(fsites, site)
	}
	// fmt.Println("Filtered:", len(fsites))

	// Find the 7 nearest metros to each metro. This is a conservative over estimate.
	nNearest := map[string][]string{}
	for _, site := range fsites {
		match := findClosestN(site, fsites, 8)
		names := []string{}
		for i := range match {
			names = append(names, match[i].Metro[1])
		}
		nNearest[site.Metro[1]] = names
		// fmt.Println("Nearest:", site.Metro[1], names)
	}

	// For each metro, color all neighbors with a distinct set of colors.
	colorMap := map[string]int{}
	for _, metroNames := range nNearest {
		// fmt.Println(colorMap)
		colorUniquely(colorMap, metroNames)
	}

	b, err := json.MarshalIndent(colorMap, "", "  ")
	rtx.Must(err, "")
	fmt.Print("metroColor = ", string(b), ";")
	// for metro, color := range colorMap {
	// fmt.Println(metro, color)
	// }

}
