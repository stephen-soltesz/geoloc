package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/m-lab/go/rtx"
)

// Site is a structured representation of site info.
type Site struct {
	Name     string   `json:"site"`
	City     string   `json:"city"`
	Metro    []string `json:"metro"`
	Country  string   `json:"country"`
	Lat      float64  `json:"latitude"`
	Lon      float64  `json:"longitude"`
	Distance float64
}

func filterProd(sites []*Site) []*Site {
	filtered := make([]*Site, 0, len(sites))
	// var ord2 *Site
	// var den2 *Site
	for i := range sites {
		if sites[i].Name[4] == 'c' || sites[i].Name[4] == 't' {
			continue
		}
		if sites[i].Name == "lba01" || sites[i].Name == "trn01" || sites[i].Name == "acc02" {
			continue
		}
		/*
			if sites[i].Name == "ord02" {
				ord2 = sites[i]
			}
			if sites[i].Name == "den02" {
				den2 = sites[i]
			}
			if sites[i].Name == "yul02" {
				sites[i].Lat = sites[i].Lat + 4
				// sites[i].Lon = sites[i].Lon - 2
			}
			if sites[i].Name == "yyz02" {
				sites[i].Lat = sites[i].Lat + 2
				sites[i].Lon = sites[i].Lon - 1
			}
			if sites[i].Name == "ywg01" {
				ywg := sites[i]
				ord2 = &Site{
					Name:  ord2.Name,
					Metro: ord2.Metro,
					Lat:   ywg.Lat - 3,
					Lon:   ywg.Lon + .1,
				}
				filtered = append(filtered, ord2)
				den2 = &Site{
					Name:  den2.Name,
					Metro: den2.Metro,
					Lat:   ywg.Lat - 3,
					Lon:   ywg.Lon - .1,
				}
				filtered = append(filtered, den2)
			}
		*/
		filtered = append(filtered, sites[i])
	}
	return filtered
}

func LoadSitesInfo(name string) []*Site {
	var data []byte
	var err error
	_, err = os.Stat(name)
	if err != nil {
		url := "https://storage.googleapis.com/operator-mlab-sandbox/metadata/v0/current/mlab-site-stats.json"
		fmt.Println(url)
		resp, err := http.Get(url)
		rtx.Must(err, "Failed to get site info")
		data, err = ioutil.ReadAll(resp.Body)
		rtx.Must(err, "Failed to read data")
		err = ioutil.WriteFile(name, data, 0640)
		rtx.Must(err, "Failed to write data to %s", name)
	} else {
		data, err = ioutil.ReadFile(name)
	}
	sites := []*Site{}
	err = json.Unmarshal(data, &sites)
	rtx.Must(err, "Failed to unmarshal site json data")
	return filterProd(sites)
}

func LoadSites(url string) []*Site {
	var data []byte
	var err error

	fmt.Println(url)
	resp, err := http.Get(url)
	rtx.Must(err, "Failed to get site info")
	data, err = ioutil.ReadAll(resp.Body)
	rtx.Must(err, "Failed to read data")

	sites := []*Site{}
	err = json.Unmarshal(data, &sites)
	rtx.Must(err, "Failed to unmarshal site json data")
	return filterProd(sites)
}
