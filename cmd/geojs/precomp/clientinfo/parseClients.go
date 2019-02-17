package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"
)

var (
	input  flagx.FileBytes
	output string
	outdir string
)

func init() {
	flag.Var(&input, "input", "File name for raw BQ output.")
	flag.StringVar(&outdir, "outdir", "", "Directory to write metro files.")
}

type rawPoint struct {
	Lat   string `json:"lat"`
	Lon   string `json:"lon"`
	Metro string `json:"metro"`
}

type goodPoint struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Metro string  `json:"metro"`
}

var (
	metros = map[string][]goodPoint{}
)

func main() {
	flag.Parse()

	var raw []rawPoint
	json.Unmarshal(input, &raw)

	for i := range raw {
		lat, _ := strconv.ParseFloat(raw[i].Lat, 64)
		lon, _ := strconv.ParseFloat(raw[i].Lon, 64)
		name := raw[i].Metro
		p := goodPoint{
			Lat:   lat,
			Lon:   lon,
			Metro: name,
		}
		metros[name] = append(metros[name], p)
	}

	found := map[string]bool{}

	for name, points := range metros {
		fmt.Print(name + " ")
		b, err := json.MarshalIndent(&points, "", "    ")
		rtx.Must(err, "Failed to Marshal good points")

		err = ioutil.WriteFile(path.Join(outdir, name+".json"), b, 0644)
		rtx.Must(err, "Failed to write file")

		found[name] = true
	}
	fmt.Println()

	b, err := json.MarshalIndent(&found, "", "    ")
	rtx.Must(err, "Failed to Marshal good points")

	err = ioutil.WriteFile("metros.json", b, 0644)
	rtx.Must(err, "Failed to write file")
}
