package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"
	"github.com/stephen-soltesz/geoloc/model"
)

var (
	metroBytes flagx.FileBytes
	output     string
)

func init() {
	flag.Var(&metroBytes, "input", "File name for raw BQ output.")
	flag.StringVar(&output, "output", "", "Directory to write metro files.")
}

type goodPoint struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Metro string  `json:"metro"`
	Color string  `json:"color"`
}

func main() {
	flag.Parse()

	sites := model.LoadSitesInfo("mlab-site-info.json")
	foundSites := map[string]goodPoint{}

	var metroMap map[string]bool
	err := json.Unmarshal(metroBytes, &metroMap)
	rtx.Must(err, "Failed to unmarshal")

	for i := range sites {
		if _, ok := metroMap[sites[i].Name[:3]]; !ok {
			continue
		}
		p := goodPoint{
			Lat:   sites[i].Lat,
			Lon:   sites[i].Lon,
			Metro: sites[i].Name[:3],
		}
		foundSites[p.Metro] = p
	}

	b, err := json.MarshalIndent(&foundSites, "", "    ")
	rtx.Must(err, "Failed to Marshal found sites")

	err = ioutil.WriteFile(output, b, 0644)
	rtx.Must(err, "Failed to write file")
}
