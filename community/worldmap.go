package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"github.com/stephen-soltesz/geoloc/community/jsx"
)

var (
	jQuery   = jquery.NewJQuery
	document = js.Global.Get("document")
	window   = js.Global.Get("window")
	location = js.Global.Get("location")
)

func randColor(alpha string) string {
	var c string
	c = "#000000"
	tinycolor := js.Global.Get("tinycolor")
	tc := tinycolor.Invoke(c)
	// Loop until we get a dark enough color.
	for tc.Call("isDark").Bool() {
		c = fmt.Sprintf("#%02x%02x%02x%s",
			jsx.RandInt(0x100),
			jsx.RandInt(0x100),
			jsx.RandInt(0x100),
			alpha)
		tc = tinycolor.Invoke(c)
	}
	return c
}
func randColor2(alpha string, i int) string {
	tinycolor := js.Global.Get("tinycolor")
	c := "#" + colors[i%len(colors)]
	tc := tinycolor.Invoke(c)
	c = tc.Call("lighten").Call("toString").String()
	tc = tinycolor.Invoke(c)
	// Loop until we get a dark enough color.
	for tc.Call("isDark").Bool() {
		c = tc.Call("lighten").Call("toString").String()
		// tinycolor("#f00").lighten().toString();
		tc = tinycolor.Invoke(c)

	}
	return c + alpha
}

func X(w int, lon float64) int {
	// return int(math.Round(float64(w) * (lon + 180.0) / 360.0))
	return int(jsx.Round(float64(w) * (lon + 180.0) / 360.0))
}
func Y(h int, lat float64) int {
	return int(jsx.Round(float64(h) * (1 - (lat+90.0)/180.0)))
}

type Site struct {
	Lat  float64
	Lon  float64
	Name string
}

type Point struct {
	X int
	Y int
	R int
}

type SiteData struct {
	sites  []*Site
	sx     []int
	sy     []int
	width  int
	height int
	color  []string
	// metros map[string][]*model.Site
	metroPoints map[string][]*Point
}

var data = &SiteData{}
var metros = map[string]bool{}
var colorMap = map[string]string{}
var colors []string

func Distance(x, y int) int {
	return x*x + y*y
}

func (d *SiteData) findMinSiteIndex(x, y int) int {
	var j int
	var cur, dmin int
	dmin = Distance(d.height, d.width)
	j = 0
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

func loadSites(canvas *jsx.Canvas) {
	rawSites := js.Global.Get("sites").Interface().([]interface{})
	data.sites = make([]*Site, 0, len(rawSites))
	for i := range rawSites {
		r := rawSites[i].(map[string]interface{})
		s := &Site{
			Lat:  r["lat"].(float64),
			Lon:  r["lon"].(float64),
			Name: r["metro"].(string),
		}
		metros[s.Name] = true
		data.color = append(data.color, randColor2("44", i+3))
		// colorMap[s.Name] = randColor2("ff", i+11)
		colorMap[s.Name] = randColor2("ff", i)
		data.sites = append(data.sites, s)
	}
	drawSites(canvas, data.sites)
}

func drawSites(c *jsx.Canvas, sites []*Site) {
	context := c.GetContext()
	data.sx = make([]int, 0, len(sites))
	data.sy = make([]int, 0, len(sites))
	data.width = c.Width()
	data.height = c.Height()
	context.Set("fillStyle", "black")
	for i := range sites {
		s := sites[i]

		data.sx = append(data.sx, X(data.width, s.Lon))
		data.sy = append(data.sy, Y(data.height, s.Lat))

		context.Set("globalAlpha", 0.8)
		context.Call("fillRect", data.sx[i], data.sy[i], 4, 4)
		context.Set("globalAlpha", 0.8)
		context.Call("fillText", s.Name, data.sx[i]+9, data.sy[i]+10)
	}
	fmt.Println("All sites drawn:", len(sites))
}

func loadDSData(dsData *js.Object, context *js.Object) map[string][]*Point {
	d := dsData.Interface().(map[string]interface{})
	b, _ := json.MarshalIndent(d["themeStyle"], "", "  ")
	fmt.Println(string(b))
	results := map[string][]*Point{}
	tables := d["tables"].(map[string]interface{})
	rows := tables["DEFAULT"].([]interface{})

	for i := range rows {
		r := rows[i].(map[string]interface{})
		metro := r["metroDimension"].([]interface{})[0].(string)

		latlon := r["geoDimension"].([]interface{})[0].(string)
		ll := strings.Split(latlon, ",")
		lat := jsx.ParseFloat(ll[0])
		lon := jsx.ParseFloat(ll[1])

		if i%(len(rows)/100) == 0 {
			window.Call("setTimeout", func() {
				context.Call("fillRect", i, 0, i+1, 10)
			}, 1)
		}
		p := &Point{
			X: X(data.width, lon),
			Y: Y(data.height, lat),
			R: jsx.Round(jsx.Random()*0.75) + 1, // 0+1 or 1+1
		}
		results[metro] = append(results[metro], p)
	}

	return results
}

func drawMetros(context *js.Object) {
	fmt.Println("drawing metros:", len(metros))
	for loc := range metros {
		// fmt.Println("loc:", loc)
		drawMetro(loc, context)
	}
}

func drawMetro(metro string, context *js.Object) {
	clients := data.metroPoints[metro]
	for j := range clients {
		p := clients[j]
		context.Set("fillStyle", colorMap[metro])

		context.Call("fillRect", p.X, p.Y, p.R, p.R)
		// context.Call("fillRect", p.X, p.Y, 1, 1)

		// Flexible shape and size, but ~5x slower UI when drawing ~50k points.
		//context.Call("beginPath")
		//context.Call("arc", p.X, p.Y, 0.8, 0, 2*jsx.PI, false)
		//context.Call("fill")
	}
}

func loadData(canvas, dsData *js.Object) {
	c := &jsx.Canvas{Object: canvas}
	loadSites(c)
	context := c.GetContext()
	data.metroPoints = loadDSData(dsData, context)

	fmt.Println("Setting ondblclick event handler")
	// document.Set("ondblclick", func(evt *js.Object) {
	canvas.Set("onclick", func(evt *js.Object) {
		x := evt.Get("offsetX").Int()
		y := evt.Get("offsetY").Int()
		fmt.Printf("CLICK: x: %d, y: %d\n", x, y)

		i := data.findMinSiteIndex(x, y)
		fmt.Println("nearest:", data.sites[i].Name)
		loc := data.sites[i].Name[:3]

		if _, ok := metros[loc]; ok {
			// The metro location is already displayed. Remove it.
			delete(metros, loc)
		} else {
			metros[loc] = true
		}
		if _, ok := data.metroPoints[loc]; ok {
			fmt.Println("drawing:", loc)
		}
		context := c.GetContext()
		context.Set("fillStyle", "#ff0000")
		context.Call("clearRect", 0, 0, c.Width(), c.Height())
		drawSites(c, data.sites)
		drawMetros(context)
	})
	btnClear := jsx.GetElementById("buttonClear")
	style := btnClear.Get("style")
	style.Set("position", "absolute")
	style.Set("left", "1px")
	style.Set("top", "0px")
	btnClear.Set("onclick", func(evt *js.Object) {
		context := c.GetContext()
		context.Call("clearRect", 0, 0, data.width, data.height)
		for loc := range metros {
			delete(metros, loc)
		}
		drawSites(c, data.sites)
		drawMetros(context)
	})
	btnDraw := jsx.GetElementById("buttonDraw")
	style = btnDraw.Get("style")
	style.Set("position", "absolute")
	style.Set("left", "1px")
	style.Set("top", "20px")
	btnDraw.Set("onclick", func(evt *js.Object) {
		context := c.GetContext()
		context.Call("clearRect", 0, 0, data.width, data.height)
		for i := range data.sites {
			metros[data.sites[i].Name] = true
		}
		drawSites(c, data.sites)
		drawMetros(context)
	})
	window.Call("setTimeout", func() {
		context := c.GetContext()
		context.Set("fillStyle", "#ff0000")
		context.Call("clearRect", 0, 0, c.Width(), c.Height())
		drawSites(c, data.sites)
		drawMetros(context)
	}, 1)
	return
}

func main() {
	icolors := js.Global.Call("palette", "mpn65", 29).Interface().([]interface{})
	// icolors := js.Global.Call("palette", "cb-Set3", 10).Interface().([]interface{})
	//icolors := js.Global.Call("palette", "cb-Paired", 9).Interface().([]interface{})
	// icolors := js.Global.Call("palette", "tol-rainbow", 10).Interface().([]interface{})
	for i := range icolors {
		colors = append(colors, icolors[i].(string))
	}
	fmt.Println(colors)
	js.Global.Set("loadData", loadData)
}
