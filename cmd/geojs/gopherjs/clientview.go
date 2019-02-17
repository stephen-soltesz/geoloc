package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strconv"

	// third-party
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/stephen-soltesz/geoloc/model"
)

var (
	jQuery   = jquery.NewJQuery
	document = js.Global.Get("document")
	window   = js.Global.Get("window")
	location = js.Global.Get("location")
)

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

func (img *Image) width() int {
	return img.Object.Get("width").Int()
}

func (img *Image) height() int {
	return img.Object.Get("height").Int()
}

func X(w int, lon float64) int {
	return int(math.Round(float64(w) * (lon + 180.0) / 360.0))
}
func Y(h int, lat float64) int {
	return int(math.Round(float64(h) * (1 - (lat+90.0)/180.0)))
}

func Distance(x, y int) int {
	return x*x + y*y
}

func randColor(alpha string) string {
	var c string
	c = "#ffffff"
	tinycolor := js.Global.Get("tinycolor")
	tc := tinycolor.Invoke(c)
	// Loop until we get a dark enough color.
	for tc.Call("isLight").Bool() {
		c = fmt.Sprintf("#%02x%02x%02x%s",
			rand.Intn(0x100),
			rand.Intn(0x100),
			rand.Intn(0x100),
			alpha)
		tc = tinycolor.Invoke(c)
	}
	return c
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
	sx     []int
	sy     []int
	width  int
	height int
	color  []string
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

func downloadMetroClients(metro string, context *js.Object) {
	jquery.GetJSON("/metros/"+metro+".json", func(rawjson interface{}) {
		p := rawjson.([]interface{})
		rawPoints := make([]RawPoint, len(p))

		for i := range p {
			r := p[i].(map[string]interface{})
			rawPoints[i].Lat = r["lat"].(float64)
			rawPoints[i].Lon = r["lon"].(float64)
			rawPoints[i].Metro = r["metro"].(string)
		}
		data.parseClients(rawPoints)
		drawMetro(metro, context)
	})
}

func drawMetros(context *js.Object) {
	for loc := range metros {
		fmt.Println("loc:", loc)
		drawMetro(loc, context)
	}
}

func drawMetro(metro string, context *js.Object) {
	clients := data.metroPoints[metro]
	for j := range clients {
		p := clients[j]
		context.Set("fillStyle", colorMap[metro])
		context.Call("fillRect", p.X, p.Y, 1, 1)
	}
}

func loadSites(canvas *js.Object) {
	jquery.GetJSON("/sites.json", func(rawjson interface{}) {
		rawSites := rawjson.([]interface{})

		sites := make([]*model.Site, 0, len(rawSites))
		for i := range rawSites {
			r := rawSites[i].(map[string]interface{})
			s := &model.Site{
				Lat:  r["lat"].(float64),
				Lon:  r["lon"].(float64),
				Name: r["metro"].(string),
			}
			sites = append(sites, s)
		}
		data.sites = sites
		drawSites(canvas)
	})
}

func getOption(key string, def string) string {
	fmt.Println("spring:", location.Get("search").String())
	v, err := url.ParseQuery(location.Get("search").String()[1:])
	if err != nil || v.Get(key) == "" {
		return def
	}
	return v.Get(key)
}

func drawSites(canvas *js.Object) {
	context := canvas.Call("getContext", "2d")
	seedRaw := getOption("seed", "140")
	seed, _ := strconv.ParseInt(seedRaw, 10, 64)
	fmt.Println("seed:", seed)
	rand.Seed(seed)
	context.Set("fillStyle", "black")
	for i := range data.sites {
		s := data.sites[i]
		fmt.Println(s.Name, s.Lat, s.Lon)

		data.sx = append(data.sx, X(data.width, s.Lon))
		data.sy = append(data.sy, Y(data.height, s.Lat))

		data.color = append(data.color, randColor("44"))
		colorMap[s.Name[:3]] = randColor("ff") // scolors[i%len(scolors)]
		context.Set("globalAlpha", 0.2)
		context.Call("fillRect", data.sx[i], data.sy[i], 4, 4)
		context.Set("globalAlpha", 0.8)
		context.Call("fillText", s.Name, data.sx[i]+9, data.sy[i]+10)

	}
}
func addCanvas(containerName, canvasName string, width, height int) *js.Object {
	canvas := document.Call("createElement", "canvas")
	canvas.Set("id", canvasName)
	canvas.Set("width", width)
	canvas.Set("height", height)
	canvas.Set("className", "canvas")

	go func() {
		data.metroPoints = make(map[string][]*Point)
		data.width = width
		data.height = height
		loadSites(canvas)
	}()

	points := document.Call("createElement", "canvas")
	points.Set("id", canvasName+"-points")
	points.Set("width", width)
	points.Set("height", height)
	points.Set("className", "canvas")

	button2 := document.Call("createElement", "INPUT")
	button2.Set("type", "button")
	button2.Set("value", "Draw Voronoi")
	button2.Set("onclick", func(evt *js.Object) {
		fmt.Println(evt)
		context := canvas.Call("getContext", "2d")
		context.Set("globalAlpha", 0.15)
		var w = width
		var h = height
		var w1 = w - 2
		var h1 = h - 2
		var n = len(data.sx)
		var DMIN = Distance(h1, w1)

		colorMinDist := func(y int) {
			for x := 0; x < w1; x++ {
				var j int
				var d, dmin int
				dmin = DMIN
				j = -1
				for i := 0; i < n; i++ {
					d = Distance(data.sx[i]-x, data.sy[i]-y)
					if d < dmin {
						dmin = d
						j = i
					}
				}

				tinycolor := js.Global.Get("tinycolor")
				tc := tinycolor.Invoke(data.color[j])
				comp := tc.Call("complement", data.color[j])
				c := comp.Call("setAlpha", 0.25).Call("toHexString").String()
				context.Set("fillStyle", c)
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

			downloadMetroClients(loc, context)
		}
	})

	button4 := document.Call("createElement", "INPUT")
	button4.Set("type", "button")
	button4.Set("value", "Clear Clients")
	button4.Set("onclick", func(evt *js.Object) {
		context := points.Call("getContext", "2d")
		context.Call("clearRect", 0, 0, width, height)
		for loc := range metros {
			delete(metros, loc)
		}
	})

	var mouseDown = false
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
			// The metro location is already displayed. Remove it.
			delete(metros, loc)
		} else {
			metros[loc] = true
		}

		if _, ok := data.metroPoints[loc]; !ok {
			// Data needs to be downloaded and cached.
			downloadMetroClients(loc, context)
		}
		drawMetros(context)

	})

	document.Set("onmousedown", func(evt *js.Object) {
		mouseDown = true
	})
	document.Set("onmouseup", func(evt *js.Object) {
		mouseDown = false
	})
	document.Set("onmousemove", func(evt *js.Object) {
		if mouseDown {
			div := document.Call("getElementById", "container")
			fmt.Println("before:", div.Get("scrollLeft").Int(), div.Get("scrollTop").Int())
			div.Set("scrollLeft", div.Get("scrollLeft").Int()-evt.Get("movementX").Int())
			div.Set("scrollTop", div.Get("scrollTop").Int()-evt.Get("movementY").Int())
			fmt.Println("after:", div.Get("scrollLeft").Int(), div.Get("scrollTop").Int())
		}
	})

	jQuery(containerName).Prepend(points)
	jQuery(containerName).Prepend(canvas)
	jQuery("#buttons").Append(button2)
	jQuery("#buttons").Append(button3)
	jQuery("#buttons").Append(button4)

	return canvas
}

func setupCanvas(containerName, uri string) {
	fmt.Println("create img")
	imgRaw := document.Call("createElement", "img")
	imgRaw.Set("id", "worldmap")
	img := &Image{imgRaw}
	img.addEventListener("load", false, func() {
		fmt.Println("WxH:", img.width(), img.height())
		canvas := addCanvas("#"+containerName, "mycanvas", img.width(), img.height())
		if canvas == nil {
			fmt.Println("NULL CANVAS")
		}
		context := canvas.Call("getContext", "2d")
		context.Set("globalAlpha", 0.2)
		window.Call("setTimeout", func() {
			fmt.Println("drawing image")
			context.Set("globalAlpha", 0.2)
			context.Call("drawImage", img.Object, 0, 0)
		}, 1)
	})
	fmt.Println("set img src", uri)
	imgRaw.Set("src", uri)
}

func setupGeoJS() {
	// Download image. Onload add canvas with image.
	fmt.Println("setup canvas")
	size := getOption("size", "medium")
	setupCanvas("container", "/"+size+"-base.png")
}

var firstRun = true

// converts the canvas to an octet-stream downloadable image.
func saveImage(canvasName, linkName string) {
	url := jQuery(canvasName).Get(0).Call("toDataURL", "image/png")
	url = url.Call("replace", "image/png", "image/octet-stream")
	jQuery(linkName).Get(0).Set("href", url)
}

func main() {
	// Export function names globally.
	js.Global.Set("getUpdateInterval", getUpdateInterval)
	js.Global.Set("setupGeoJS", setupGeoJS)
	js.Global.Set("addCanvas", addCanvas)
	js.Global.Set("setupCanvas", setupCanvas)
	js.Global.Set("appendLog", appendLog)
	js.Global.Set("saveImage", saveImage)
}
