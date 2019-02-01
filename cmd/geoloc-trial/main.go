package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"
)

var (
	url   string
	input flagx.FileBytes
)

func init() {
	flag.StringVar(&url, "url", "https://mlab-ns.appspot.com/ndt_ssl?ip_address=", "URL prefix used to request")
	flag.Var(&input, "file", "File to read contents.")
}

func main() {
	flag.Parse()

	buf := bytes.NewBuffer(input)
	r := csv.NewReader(buf)

	for rec, err := r.Read(); err == nil; rec, err = r.Read() {
		ip := rec[0]
		site := rec[1]

		resp, err := http.Get(url + ip)
		rtx.Must(err, "Failed to get url for: %s", url)

		b, err := ioutil.ReadAll(resp.Body)
		rtx.Must(err, "Failed to read response body")

		m := map[string]interface{}{}
		err = json.Unmarshal(b, &m)
		rtx.Must(err, "Failed to unmarshal response: %s", string(b))

		fmt.Printf("%s appengine:%s maxmind:%s\n", ip, site[:3], m["site"].(string)[:3])
	}
}
