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
	for i := range sites {
		if sites[i].Name[4] == 'c' || sites[i].Name[4] == 't' {
			continue
		}
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
