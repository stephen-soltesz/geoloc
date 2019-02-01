package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/m-lab/go/flagx"
)

var (
	url   string
	input flagx.FileBytes
)

func init() {
	flag.StringVar(&url, "url", "https://mlab-ns.appspot.com/ndt_ssl?ip_address=", "URL prefix used to request")
	flag.Var(&input, "file", "File to read contents.")

	log.SetOutput(os.Stderr)
}

func main() {
	flag.Parse()

	buf := bytes.NewBuffer(input)
	r := csv.NewReader(buf)

	for rec, err := r.Read(); err == nil; rec, err = r.Read() {
		// ts := rec[0]
		ip := rec[1]
		site := rec[2]

		resp, err := http.Get(url + ip)
		if err != nil {
			log.Printf("Failed to get url for: %s %s", url, err)
			continue
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %s", err)
			continue
		}

		m := map[string]interface{}{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			log.Printf("Failed to unmarshal response: %s : %s", string(b), err)
			continue
		}

		fmt.Printf("%s,%s,%s,%s\n", time.Now().Format(time.RFC3339), site[:3], m["site"].(string)[:3], ip)
	}
}
