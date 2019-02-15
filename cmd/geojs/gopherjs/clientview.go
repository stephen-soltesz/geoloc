package main

import (
	"fmt"
	"math"
	"math/rand"

	// third-party
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/stephen-soltesz/geoloc/model"
)

var jQuery = jquery.NewJQuery
var document = js.Global.Get("document")
var window = js.Global.Get("window")
var xOffset float64 = 0
var plotSamples int64 = 240
var updateInterval float64 = 1000

func getUpdateInterval() float64 {
	return updateInterval
}

func appendLog(msg jquery.JQuery) {
	var log = jQuery("#log")
	d := log.Underlying().Index(0)
	msg.AppendTo(log)
	scrollTop := d.Get("scrollTop").Int()
	scrollHeight := d.Get("scrollHeight").Int()
	clientHeight := d.Get("clientHeight").Int()
	doScroll := (scrollTop < scrollHeight-clientHeight)
	if doScroll {
		d.Set("scrollTop", scrollHeight-clientHeight)
	}
}

type Image struct {
	*js.Object
}

func newImage(src string) *Image {
	img := document.Call("createElement", "img")
	img.Set("src", src)
	return &Image{img}
}

func (img *Image) addEventListener(event string, capture bool, callback func()) {
	img.Call("addEventListener", event, callback, capture)
}

func getXOffset(inc bool) int64 {
	if inc && xOffset < 0 {
		xOffset += (updateInterval / 1000.0)
	}
	return int64(math.Floor(xOffset))
}

func getPlotSamples() int64 {
	return plotSamples
}

func X(w int, lon float64) int {
	return int(math.Round(float64(w) * (lon + 180.0) / 360.0))
}
func Y(h int, lat float64) int {
	return int(math.Round(float64(h) * (1 - (lat+90.0)/180.0)))
}

func Distance(x, y int) int {
	return x*x + y*y
	// return (Math.pow(Math.pow(Math.abs(x), 3)+Math.pow(Math.abs(y), 3), 0.33333))
}

func randColor(alpha string) string {
	return fmt.Sprintf("#%02x%02x%02x%s",
		rand.Intn(0x100),
		rand.Intn(0x100),
		rand.Intn(0x100),
		alpha)
}

type RawPoint struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Metro string  `json:"metro"`
}
type Point struct {
	X int
	Y int
}

type Data struct {
	sites  []*model.Site
	width  int
	height int
	// metros map[string][]*model.Site
	metroPoints map[string][]*Point
}

func (d *Data) findMinSiteIndex(x, y int) int {
	var j int
	var cur, dmin int
	dmin = Distance(d.height, d.width)
	j = -1
	for i := 0; i < len(d.sites); i++ {
		sx := X(d.width, d.sites[i].Lon)
		sy := Y(d.height, d.sites[i].Lat)
		cur = Distance(sx-x, sy-y)
		if cur < dmin {
			dmin = cur
			j = i
		}
	}
	return j
}

func (d *Data) parseClients(data []RawPoint) {
	for i := range data {
		r := data[i]
		p := Point{
			X: X(d.width, r.Lon),
			Y: Y(d.height, r.Lat),
		}
		d.metroPoints[r.Metro] = append(d.metroPoints[r.Metro], &p)
	}
}

var data Data
var metros = map[string]bool{}
var colorMap = map[string]string{}
var scolors = []string{
	"#ff0000",
	"#ff8800",
	"#00ff00",
	"#00ffff",
	"#0088ff",
	"#0000ff",
	"#8800ff",
	"#ff00ff",
}

func addCanvas(containerName, canvasName string, width, height int) {
	var sx, sy []int
	var color []string

	go func() {
		data.metroPoints = make(map[string][]*Point)
		data.sites = model.LoadSites("http://localhost:8080/sites2.json")
		data.width = width
		data.height = height
	}()

	go func() {
		jquery.GetJSON("/clients2.json", func(rawjson interface{}) {
			p := rawjson.([]interface{})
			rawPoints := make([]RawPoint, len(p))

			for i := range p {
				r := p[i].(map[string]interface{})
				rawPoints[i].Lat = r["lat"].(float64)
				rawPoints[i].Lon = r["lon"].(float64)
				rawPoints[i].Metro = r["metro"].(string)
			}
			data.parseClients(rawPoints)
		})
	}()

	canvas := document.Call("createElement", "canvas")
	canvas.Set("id", canvasName)
	canvas.Set("width", width)
	canvas.Set("height", height)
	canvas.Set("className", "canvas")

	points := document.Call("createElement", "canvas")
	points.Set("id", canvasName+"-points")
	points.Set("width", width)
	points.Set("height", height)
	points.Set("className", "canvas")

	button := document.Call("createElement", "INPUT")
	button.Set("type", "button")
	button.Set("value", "Load Sites")

	button.Set("onclick", func(evt *js.Object) {
		fmt.Println(evt)
		//jquery.GetJSON("/sites2.json", func(data interface{}) {
		//fmt.Println(data)
		context := canvas.Call("getContext", "2d")
		context.Set("fillStyle", "black")
		rand.Seed(42)

		for i := range data.sites {
			s := data.sites[i]
			fmt.Println(s.Name, s.Lat, s.Lon)

			sx = append(sx, X(width, s.Lon))
			sy = append(sy, Y(height, s.Lat))

			// color = append(color, randColor("44"))
			color = append(color, randColor("44"))
			colorMap[s.Name[:3]] = randColor("ff") // scolors[i%len(scolors)]
			context.Call(
				"fillRect",
				sx[i],
				sy[i],
				8, 8)
		}
		/*
			sites2, ok := data.([]interface{})
			if !ok {
				fmt.Println("ERROR: parsing sites2.json")
				return
			}
			context := canvas.Call("getContext", "2d")
			context.Set("fillStyle", "black")
			rand.Seed(42)

			for i := range sites2 {
				s := sites2[i].(map[string]interface{})
				fmt.Println(s["site"], s["latitude"], s["longitude"])
				sx = append(sx, X(1260, s["longitude"].(float64)))
				sy = append(sy, Y(630, s["latitude"].(float64)))
				color = append(color, randColor("44"))
				context.Call(
					"fillRect",
					sx[i],
					sy[i],
					3, 3)
			}
		*/
		//})
	})

	button2 := document.Call("createElement", "INPUT")
	button2.Set("type", "button")
	button2.Set("value", "Draw Voronoi")
	button2.Set("onclick", func(evt *js.Object) {
		fmt.Println(evt)
		context := canvas.Call("getContext", "2d")
		context.Set("globalAlpha", 0.5)
		var w = width
		var h = height
		var w1 = w - 2
		var h1 = h - 2
		var n = len(sx)
		var DMIN = Distance(h1, w1)

		colorMinDist := func(y int) {
			for x := 0; x < w1; x++ {
				var j int
				var d, dmin int
				dmin = DMIN
				j = -1
				for i := 0; i < n; i++ {
					d = Distance(sx[i]-x, sy[i]-y)
					if d < dmin {
						dmin = d
						j = i
					}
				}
				context.Set("fillStyle", color[j])
				context.Call("fillRect", x, y, 1, 1)
			}
		}

		for y := 0; y < h1; y++ {
			window.Call("setTimeout", colorMinDist, 1, y)
		}
	})

	button3 := document.Call("createElement", "INPUT")
	button3.Set("type", "button")
	button3.Set("value", "Draw All Clients")
	button3.Set("onclick", func(evt *js.Object) {
		fmt.Println(evt)
		context := points.Call("getContext", "2d")

		for i := range data.sites {
			loc := data.sites[i].Name[:3]
			if _, ok := metros[loc]; ok {
				continue
			}
			metros[loc] = true
			fmt.Println("loc:", loc)
			clients := data.metroPoints[loc]
			for j := range clients {
				p := clients[j]
				// context.Set("fillStyle", "#ff0000")
				context.Set("fillStyle", colorMap[loc])
				context.Call("fillRect", p.X, p.Y, 1, 1)
			}
		}
	})

	var mouseDown = false
	var mousePrevPos *js.Object
	document.Set("onmousedown", func(evt *js.Object) {
		mouseDown = true
		mousePrevPos = nil
		updateInterval = 200
	})
	document.Set("onmouseup", func(evt *js.Object) {
		mouseDown = false
		updateInterval = 1000
	})
	document.Set("ondblclick", func(evt *js.Object) {
		mouseDown = false
		updateInterval = 1000

		x := evt.Get("offsetX")
		y := evt.Get("offsetY")

		context := points.Call("getContext", "2d")
		context.Set("fillStyle", "#ff0000")
		context.Call("clearRect", 0, 0, width, height)

		i := data.findMinSiteIndex(x.Int(), y.Int())
		fmt.Println(data.sites[i].Name)
		loc := data.sites[i].Name[:3]

		if _, ok := metros[loc]; ok {
			// the metro location is already displayed. Remove it.
			delete(metros, loc)
		} else {
			metros[loc] = true
		}

		for loc, _ := range metros {
			fmt.Println("loc:", loc)
			clients := data.metroPoints[loc]
			for j := range clients {
				p := clients[j]
				// context.Set("fillStyle", "#ff0000")
				context.Set("fillStyle", colorMap[loc])
				context.Call("fillRect", p.X, p.Y, 1, 1)
			}
		}
	})
	document.Set("onmousemove", func(evt *js.Object) {
		if mouseDown {
			// movement attributes are not supported universally.
			if mousePrevPos != nil {
				curr := evt.Get("clientX").Int()
				prev := mousePrevPos.Get("clientX").Int()
				xOffset += float64(curr - prev)
			}
			mousePrevPos = evt
		}
	})
	document.Set("onmousewheel", func(evt *js.Object) {
		scroll := evt.Get("wheelDeltaY").Int()
		if scroll > 0 && plotSamples < 1000 {
			plotSamples += 10
		} else if scroll < 0 && plotSamples > 60 {
			plotSamples -= 10
		}
	})
	jQuery(containerName).Prepend(points)
	jQuery(containerName).Prepend(canvas)
	jQuery("#buttons").Append(button)
	jQuery("#buttons").Append(button2)
	jQuery("#buttons").Append(button3)
}

func updateCanvas(name, uri string) {
	canvas := jQuery(name).Underlying().Index(0)
	context := canvas.Call("getContext", "2d")

	img := newImage(uri)
	img.addEventListener("load", false, func() {
		context.Set("globalAlpha", 0.2)
		context.Call("drawImage", img.Object, 0, 0)
	})
}

var firstRun = true

func jsOnConfig(containerName string, data *js.Object) {
	if firstRun {
		width := data.Get("width").Int()
		height := data.Get("height").Int()
		plotSamples = data.Get("samples").Int64()
		addCanvas(containerName, "mycanvas", width, height)
		firstRun = false
	}
}

// opens a websocket to socketUrl and adds a canvas to containerName
func setupSocket(socketUrl, containerName string) {
	jquery.GetJSON("/config", func(data *js.Object) {
		jsOnConfig(containerName, data)
	})
}

// converts the canvas to an octet-stream downloadable image.
func saveImage(canvasName, linkName string) {
	url := jQuery(canvasName).Get(0).Call("toDataURL", "image/png")
	url = url.Call("replace", "image/png", "image/octet-stream")
	jQuery(linkName).Get(0).Set("href", url)
}

func main() {
	// export function names globally.
	js.Global.Set("getXOffset", getXOffset)
	js.Global.Set("getPlotSamples", getPlotSamples)
	js.Global.Set("getUpdateInterval", getUpdateInterval)
	js.Global.Set("setupSocket", setupSocket)
	js.Global.Set("addCanvas", addCanvas)
	js.Global.Set("updateCanvas", updateCanvas)
	js.Global.Set("appendLog", appendLog)
	js.Global.Set("saveImage", saveImage)
}
