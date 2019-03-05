package jsx

import (
	"github.com/gopherjs/gopherjs/js"
)

var (
	Math     = js.Global.Get("Math")
	document = js.Global.Get("document")
	PI       = Math.Get("PI").Float()
)

func ParseFloat(v string) float64 {
	return js.Global.Call("parseFloat", v).Float()
}

func ParseInt(v string) int64 {
	return js.Global.Call("parseInt", v).Int64()
}

func GetElementById(id string) *js.Object {
	return document.Call("getElementById", id)
}

func GetCanvasById(id string) *Canvas {
	return &Canvas{Object: document.Call("getElementById", id)}
}

func Random() float64 {
	return Math.Call("random").Float()
}

func RandInt(max int64) int64 {
	return int64(Math.Call("random").Float() * float64(max))
}

func Round(v float64) int {
	return Math.Call("round", v).Int()
}

type Image struct {
	*js.Object
}

func NewImage() *Image {
	img := document.Call("createElement", "img")
	return &Image{img}
}

func (img *Image) AddEventListener(event string, capture bool, callback func()) {
	img.Call("addEventListener", event, callback, capture)
}

func (img *Image) Width() int {
	return img.Object.Get("width").Int()
}

func (img *Image) Height() int {
	return img.Object.Get("height").Int()
}

type Canvas struct {
	*js.Object
}

func (c *Canvas) GetContext() *js.Object {
	return c.Call("getContext", "2d")
}

func (c *Canvas) Width() int {
	return c.Object.Get("width").Int()
}

func (c *Canvas) Height() int {
	return c.Object.Get("height").Int()
}
