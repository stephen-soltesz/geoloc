package main

import (
	"fmt"
	"math"
	"math/rand"

	// third-party
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
)

var jQuery = jquery.NewJQuery
var document = js.Global.Get("document")
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

func addCanvas(containerName, canvasName string, width, height int) {
	var sx, sy []int
	var color []string
	canvas := document.Call("createElement", "canvas")
	canvas.Set("id", canvasName)
	canvas.Set("width", width)
	canvas.Set("height", height)

	button := document.Call("createElement", "INPUT")
	button.Set("type", "button")
	button.Set("value", "Load Sites")
	button.Set("onclick", func(evt *js.Object) {
		fmt.Println(evt)
		jquery.GetJSON("/sites2.json", func(data interface{}) {
			fmt.Println(data)
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
		})
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
		var d, dmin int
		var j int
		var w1 = w - 2
		var h1 = h - 2
		var n = len(sx)
		for y := 0; y < h1; y++ {
			for x := 0; x < w1; x++ {
				dmin = Distance(h1, w1)
				j = -1
				for i := 0; i < n; i++ {
					d = Distance(sx[i]-x, sy[i]-y)
					if d < dmin {
						dmin = d
						j = i
					}
				}

				// var p = c.getImageData(x, y, 1, 1).data;
				// var hex = "#" + ("000000" + rgbToHex(p[0], p[1], p[2])).slice(-6);
				// var color = getColorStr(color1[0], color1[1], color1[2], color2[0], color2[1], color2[2], cyclePct);

				context.Set("fillStyle", color[j])
				context.Call("fillRect", x, y, 1, 1)
				// ctx.fillRect(x, y, 1, 1)

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
		// fmt.Println("X:", evt.Get("clientX"))
		// fmt.Println("Y:", evt.Get("clientY"))

		context := canvas.Call("getContext", "2d")
		context.Set("fillStyle", "#ff0000")
		context.Call("fillRect", x, y, 4, 4)

		// fmt.Println("X:", evt.Get("offsetX"))
		// fmt.Println("X:", evt.Get("pageX"))
		// fmt.Println("X:", evt.Get("screenX"))
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
	jQuery(containerName).Prepend(canvas)
	jQuery(containerName).Append(button)
	jQuery(containerName).Append(button2)
}

func updateCanvas(name, uri string) {
	canvas := jQuery(name).Underlying().Index(0)
	context := canvas.Call("getContext", "2d")

	img := newImage(uri)
	img.addEventListener("load", false, func() {
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
