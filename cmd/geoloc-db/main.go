package main

import (
	"bytes"
	"container/heap"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"

	"github.com/stephen-soltesz/geoloc/topk"

	"github.com/stephen-soltesz/geoloc/model"

	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"
	"github.com/oschwald/geoip2-golang"
)

var (
	geofile flagx.FileBytes
	input   flagx.FileBytes
)

func init() {
	flag.Var(&geofile, "geolite2", "File to read for geolite2 db.")
	flag.Var(&input, "file", "File to read contents.")
	log.SetOutput(os.Stderr)
}

func findClosestN(city *geoip2.City, sites []*model.Site, n int) []*model.Site {
	nearest := make([]*model.Site, 0, n)
	for i := range sites {
		site := sites[i]
		a := s2.LatLng{
			Lat: s1.Angle(site.Lat * math.Pi / 180.0),
			Lng: s1.Angle(site.Lon * math.Pi / 180.0),
		}
		b := s2.LatLng{
			Lat: s1.Angle(city.Location.Latitude * math.Pi / 180.0),
			Lng: s1.Angle(city.Location.Longitude * math.Pi / 180.0),
		}
		sites[i].Distance = (float64)(a.Distance(b))
	}
	h := &topk.SiteDistance{}
	heap.Init(h)
	for i := 0; i < n; i++ {
		heap.Push(h, sites[i])
	}
	for i := range sites[n:] {
		if h.MinDistance() >= sites[i].Distance {
			heap.Pop(h)
			heap.Push(h, sites[i])
		}
	}
	fmt.Println("len", h.Len())
	nearest = append(nearest, heap.Pop(h).(*model.Site))
	for i := 0; i < (n - 1); i++ {
		nearest = append(nearest, heap.Pop(h).(*model.Site))
	}
	fmt.Println("len", h.Len())
	/*
		TODO: run benchmark to compare heap with simple sort.
		sort.Slice(*sites, func(i, j int) bool {
			return (*sites)[i].Distance < (*sites)[j].Distance
		})
	*/
	return nearest
}

func main() {
	flag.Parse()

	reader, err := geoip2.FromBytes(geofile)
	rtx.Must(err, "Failed to read geolite2 file")

	buf := bytes.NewBuffer(input)
	r := csv.NewReader(buf)
	sites := model.LoadSitesInfo("mlab-site-stats.json")

	for rec, err := r.Read(); err == nil; rec, err = r.Read() {
		// ts := rec[0]
		ip := rec[1]
		site := rec[2]

		city, err := reader.City(net.ParseIP(ip))
		if err != nil {
			log.Println(err)
			continue
		}

		match := findClosestN(city, sites, 1)

		fmt.Println(ip, city.Location.Latitude, city.Location.Longitude, site, match[0].Name, match[0].Distance)
	}
	// fmt.Println("ok", geoip2.City{})
}
