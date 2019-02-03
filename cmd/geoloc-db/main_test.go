package main

import (
	"math"
	"testing"

	"github.com/golang/geo/s1"
	h3 "github.com/golang/geo/s2"
	h4 "github.com/kellydunn/golang-geo"
	h1 "github.com/umahmood/haversine"
	h2 "pault.ag/go/haversine"
)

type site struct {
	name string
	lat  float64
	lon  float64
}

var (
	sites = []site{
		{"yyz02", 43.6767, -79.6306},
		{"tyo03", 35.5522, 139.78},
		{"syd02", -33.9461, 151.177},
		{"sea07", 47.4489, -122.3094},
		{"par05", 48.8584, 2.349},
		{"ord01", 41.9786, -87.9047},
		{"mia05", 25.7833, -80.2667},
		{"lga07", 40.7667, -73.8667},
		{"lhr03", 51.4697, -0.4514},
		{"lax01", 33.9425, -118.4072},
		{"iad03", 38.9444, -77.4558},
		{"dfw07", 32.8969, -97.0381},
		{"bru04", 50.4974, 3.3528},
		{"atl04", 33.6367, -84.4281},
		{"arn03", 59.6519, 17.9186},
		{"ams03", 52.3086, 4.7639},
		{"acc02", 5.606, -0.1681},
	}
	city = site{"home1", 40.7012, -73.9436}
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Benchmark_h1(b *testing.B) {
	var m, k float64
	for n := 0; n < b.N; n++ {
		for _, site := range sites {
			a := h1.Coord{Lat: site.lat, Lon: site.lon}
			b := h1.Coord{Lat: city.lat, Lon: city.lon}
			m, k = h1.Distance(a, b)
		}
	}
	_ = m + k
	// fmt.Println("h1", k)
}

func Benchmark_h2(b *testing.B) {
	var k float64
	for n := 0; n < b.N; n++ {
		for _, site := range sites {
			a := h2.Point{Lat: site.lat, Lon: site.lon}
			b := h2.Point{Lat: city.lat, Lon: city.lon}

			k = float64(a.MetresTo(b))
		}
	}
	_ = k
	// fmt.Println("h2", k)
}

func Benchmark_h3(b *testing.B) {
	var k float64
	for n := 0; n < b.N; n++ {
		for _, site := range sites {
			a := h3.LatLng{Lat: s1.Angle(site.lat * math.Pi / 180.0), Lng: s1.Angle(site.lon * math.Pi / 180.0)}
			b := h3.LatLng{Lat: s1.Angle(city.lat * math.Pi / 180.0), Lng: s1.Angle(city.lon * math.Pi / 180.0)}
			k = (float64)(a.Distance(b))
		}
	}
	k = k * 6371
	// fmt.Println("h3", k)
}

func Benchmark_h4(b *testing.B) {
	var k float64
	for n := 0; n < b.N; n++ {
		for _, site := range sites {
			a := h4.NewPoint(site.lat, site.lon)
			b := h4.NewPoint(city.lat, city.lon)
			k = a.GreatCircleDistance(b)
		}
	}
	_ = k
	// fmt.Println("h4", k)
}
