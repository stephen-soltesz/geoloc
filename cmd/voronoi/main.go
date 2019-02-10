// Logic derived from https://rosettacode.org/wiki/Voronoi_diagram#Go
package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"strconv"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/m-lab/go/flagx"
	"github.com/stephen-soltesz/geoloc/model"
)

// type Converter interface {
// 	X() int64
// 	Y() int64
// }

// Point represents a latitude and longitude and is used to convert to X & Y coordinates.
type Point struct {
	Lat    float64
	Lon    float64
	Width  int
	Height int
}

// X converts pixel position of point.
func (p *Point) X() int {
	// return int(float64(p.Width) * ((p.Lon + 180.0) / 390.0))
	return int(float64(p.Width) * ((p.Lon + 180.0) / 360.0))
}

// Y converts pixel position of point.
func (p *Point) Y() int {
	// return int(float64(p.Height)*(1-(p.Lat+90.0)/185.0) - 5)
	return int(float64(p.Height) * (1 - (p.Lat+90.0)/180.0))
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{20, 20, 20, 255}
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func generateVoronoi(sx, sy []int, img *image.RGBA) image.Image {
	// generate a random color for each site
	sc := make([]color.NRGBA, len(sx))
	rand.Seed(41)
	for i := range sx {
		sc[i] = color.NRGBA{
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			uint8(rand.Intn(256)),
			75}
	}

	// generate diagram by coloring each pixel with color of nearest site
	r := img.Bounds()
	for x := 0; x < r.Max.X; x++ {
		for y := 0; y < r.Max.Y; y++ {
			dMin := dot(r.Max.X, r.Max.Y)
			var sMin int
			// Find the nearest site to the current x,y position.
			for s := 0; s < len(sx); s++ {
				if d := dot(sx[s]-x, sy[s]-y); d < dMin {
					sMin = s
					dMin = d
				}
			}
			// Combine the base image with the uniform voronoi cell color.
			draw.Draw(
				img, image.Rect(x, y, x+1, y+1),
				image.NewUniform(sc[sMin]), image.ZP,
				draw.Over)
		}
	}
	return img
}

func dot(x, y int) int {
	return x*x + y*y
}

func addSitesAndVoronoiToMap(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	img, t, err := image.Decode(f)
	if err != nil {
		return err
	}
	rgb := img.(*image.RGBA)
	log.Println(t)

	r := img.Bounds()
	fmt.Println("max", r.Max)

	plotClients(rgb)

	// Load all M-Lab sites.
	sites := model.LoadSitesInfo("mlab-site-stats.json")
	sx := make([]int, len(sites))
	sy := make([]int, len(sites))

	// Plot points and labels for every site.
	for i := range sites {
		s := sites[i]
		p := Point{
			Lat:    s.Lat,
			Lon:    s.Lon,
			Width:  r.Max.X,
			Height: r.Max.Y,
		}
		sx[i] = p.X()
		sy[i] = p.Y()
		draw.Draw(
			rgb, image.Rect(p.X()-2, p.Y()-2, p.X()+2, p.Y()+2),
			image.NewUniform(color.RGBA{255, 0, 0, 255}), image.ZP, draw.Src,
		)
		addLabel(rgb, p.X()-12, p.Y()-5, s.Metro[1])
	}
	generateVoronoi(sx, sy, rgb)

	writePngFile(img, "new-"+file)
	return nil

}

func plotClients(img *image.RGBA) {
	if len(clientLocations) == 0 {
		log.Println("clientLocations is empty")
		return
	}
	log.Println("plotting client locations")
	buf := bytes.NewBuffer(clientLocations)
	r := csv.NewReader(buf)
	b := img.Bounds()
	for rec, err := r.Read(); err == nil; rec, err = r.Read() {
		lat, _ := strconv.ParseFloat(rec[1], 64)
		lon, _ := strconv.ParseFloat(rec[2], 64)
		p := Point{
			Lat:    lat,
			Lon:    lon,
			Width:  b.Max.X,
			Height: b.Max.Y,
		}
		c := color.RGBA{uint8(0), uint8(0), uint8(0), 255}

		// fmt.Println(p.X(), p.Y(), c)
		img.SetRGBA(p.X(), p.Y(), c)
	}
}

func writePngFile(img image.Image, name string) {
	f, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = png.Encode(f, img); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}
}

var clientLocations flagx.FileBytes

func init() {
	flag.Var(&clientLocations, "clients", "File to read contents.")
}

func main() {
	flag.Parse()
	addSitesAndVoronoiToMap("base.png")
}
