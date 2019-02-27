package main

import (
	"fmt"
	"math"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var (
	jQuery   = jquery.NewJQuery
	document = js.Global.Get("document")
	window   = js.Global.Get("window")
	location = js.Global.Get("location")
)

type Image struct {
	*js.Object
}

func newImage() *Image {
	img := document.Call("createElement", "img")
	return &Image{img}
}

func (img *Image) addEventListener(event string, capture bool, callback func()) {
	img.Call("addEventListener", event, callback, capture)
}

func (img *Image) width() int {
	return img.Object.Get("width").Int()
}

func (img *Image) height() int {
	return img.Object.Get("height").Int()
}

type Canvas struct {
	*js.Object
}

func (c *Canvas) width() int {
	return c.Object.Get("width").Int()
}

func (c *Canvas) height() int {
	return c.Object.Get("height").Int()
}

func drawImage(canvas *js.Object, uri string) {
	fmt.Println("create img")
	img := newImage()
	img.Set("id", "worldmap")
	img.addEventListener("load", false, func() {
		fmt.Println("WxH:", img.width(), img.height())
		context := canvas.Call("getContext", "2d")
		context.Set("globalAlpha", 0.2)
		window.Call("setTimeout", func() {
			fmt.Println("drawing image")
			context.Set("globalAlpha", 0.2)
			context.Call("drawImage", img.Object, 0, 0)
			context.Call("drawImage", img.Object,
				0, 0, img.width(), img.height(),
				0, 0, canvas.Get("width").Int(), canvas.Get("height").Int())
		}, 1)
	})
	fmt.Println("set img src", uri)
	// img.Set("src", uri)
}

func X(w int, lon float64) int {
	return int(math.Round(float64(w) * (lon + 180.0) / 360.0))
}
func Y(h int, lat float64) int {
	return int(math.Round(float64(h) * (1 - (lat+90.0)/180.0)))
}

type Site struct {
	Lat  float64
	Lon  float64
	Name string
}

func loadSites(canvas *js.Object) {
	rawSites := js.Global.Get("sites").Interface().([]interface{})
	sites := make([]*Site, 0, len(rawSites))
	for i := range rawSites {
		r := rawSites[i].(map[string]interface{})
		s := &Site{
			Lat:  r["lat"].(float64),
			Lon:  r["lon"].(float64),
			Name: r["metro"].(string),
		}
		sites = append(sites, s)
	}
	drawSites(canvas, sites)
}

func drawSites(canvas *js.Object, sites []*Site) {
	context := canvas.Call("getContext", "2d")
	c := &Canvas{canvas}
	sx := []int{}
	sy := []int{}
	width := c.width()
	height := c.height()
	context.Set("fillStyle", "black")
	for i := range sites {
		s := sites[i]
		fmt.Println(s.Name, s.Lat, s.Lon)

		sx = append(sx, X(width, s.Lon))
		sy = append(sy, Y(height, s.Lat))

		// data.color = append(data.color, randColor("44"))
		// colorMap[s.Name[:3]] = randColor("ff") // scolors[i%len(scolors)]

		context.Set("globalAlpha", 0.2)
		context.Call("fillRect", sx[i], sy[i], 4, 4)
		context.Set("globalAlpha", 0.8)
		context.Call("fillText", s.Name, sx[i]+9, sy[i]+10)
	}
}

func loadData(canvas, data *js.Object) {
	fmt.Println("THIS IS A TEST:", data.String())
	// loadImage(canvas, "https://storage.googleapis.com/soltesz-mlab-sandbox/v2/small-base.png")
	drawImage(canvas, "") // js.Global.Get("worldmapSmall").String())
	loadSites(canvas)
	fmt.Println("Setting ondblclick event handler")
	document.Set("ondblclick", func(evt *js.Object) {
		x := evt.Get("offsetX").Int()
		y := evt.Get("offsetY").Int()
		fmt.Printf("CLICK: x: %d, y: %d\n", x, y)
	})
	/*
		document.Set("onmousedown", func(evt *js.Object) {
			x := evt.Get("offsetX")
			y := evt.Get("offsetY")
			fmt.Printf("down: x: %d, y: %d\n", x, y)
		})
	*/
	/*document.Set("onmousemove", func(evt *js.Object) {
		fmt.Println("dX:", evt.Get("movementX").Int())
		fmt.Println("dY:", evt.Get("movementY").Int())
	})
	*/
	/*
		canvas.Set("onmousemove", func(evt *js.Object) {
			fmt.Println("cX:", evt.Get("movementX").Int())
			fmt.Println("cY:", evt.Get("movementY").Int())
		})
	*/
	return
}

func main() {
	js.Global.Set("loadData", loadData)
}
