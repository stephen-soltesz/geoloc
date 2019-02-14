package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strconv"
	"strings"
	"text/template"

	"github.com/m-lab/go/rtx"
)

var (
	bindAddress   = flag.String("hostname", "0.0.0.0", "Interface to bind to.")
	viewerPort    = flag.Int("viewer_port", 8080, "Port for web viewer.")
	collectorPort = flag.Int("collector_port", 3131, "Port for data collector.")
	timestamp     = flag.Bool("timestamp", false, "Use timestamps as x-axis.")
	plotWidth     = flag.Int("plot_width", 1260, "Plot width in pixels.")
	plotHeight    = flag.Int("plot_height", 630, "Plot height in pixels.")
	plotSamples   = flag.Int("samples", 240, "Number of samples wide to make plots.")
	debug         = flag.Bool("debug", false, "Enable debug messages on stderr.")
	debugLogger   *log.Logger
	profile       = flag.Bool("profile", false, "Enable profiling.")
)

func checkFlags() {
	if *viewerPort == *collectorPort {
		fmt.Println("Error: viewer and collector cannot use the same port.")
		os.Exit(1)
	}
	if *debug {
		debugLogger = log.New(os.Stderr, "", 0)
	} else {
		debugLogger = log.New(ioutil.Discard, "", 0)
	}
}

func main() {
	flag.Parse()
	checkFlags()

	if *profile {
		f, err := os.Create("lineviewer.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// collector := collection.Default()
	// collector.Usetime = *timestamp
	startViewServer(*bindAddress, *viewerPort)
	// startCollectorServer(*bindAddress, *collectorPort)
	//select {}
}

func startViewServer(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	debugLogger.Printf("HTTP: listen on: %s\n", addr)

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/svg", serveSvg)
	http.HandleFunc("/config", serveConfig)
	http.HandleFunc("/json", serveJSON)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

const (
	// directory of local resources
	localPrefix = "resources"
)

func getPng(offset, samples int64) []byte {
	b, err := ioutil.ReadFile("resources/base.png")
	rtx.Must(err, "failed to open base.png")
	// collector := collection.Default()
	// err := collector.Plot(&img, *plotWidth, *plotHeight, float64(offset), float64(samples))
	// if err != nil {
	// return nil
	// }
	return b
}

type Config struct {
	Success      bool `json:"success"`
	WidthPixels  int  `json:"width"`
	HeightPixels int  `json:"height"`
	PlotSamples  int  `json:"samples"`
}

func serveConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cfg := &Config{Success: true, WidthPixels: *plotWidth,
		HeightPixels: *plotHeight, PlotSamples: *plotSamples}
	msg, err := json.Marshal(cfg)
	if err != nil {
		http.Error(w, "Server Error", 500)
		return
	}
	w.Write(msg)
	return
}

func serveJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	msg := []byte(`{}`)
	// msg, err := json.Marshal(collection.Default())
	// if err != nil {
	// 	http.Error(w, "Server Error", 500)
	// 	return
	// }
	w.Write(msg)
	return
}

func getFormValue(formVals url.Values, key string, defVal int64) (int64, error) {
	var val int64
	var err error

	valStr, ok := formVals[key]
	if !ok {
		return defVal, nil
	}

	if len(valStr) > 0 {
		val, err = strconv.ParseInt(valStr[0], 10, 64)
		if err != nil {
			return defVal, err
		}
	} else {
		val = defVal
	}
	return val, nil
}

func splitAll(path string) (string, string, string) {
	dir, file := filepath.Split(path)
	if file == "" {
		file = "home.html"
	}
	ext := strings.Trim(filepath.Ext(file), ".")
	return dir, file, ext
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var err error

	fmt.Println(r.URL)
	if r.Method != "GET" {
		http.Error(w, "Method nod allowed", 405)
		return
	}
	_, file, ext := splitAll(r.URL.Path)

	resourcefile := fmt.Sprintf("%s/%s", localPrefix, file)
	fmt.Println(resourcefile)
	data, err = ioutil.ReadFile(resourcefile)

	if err != nil {
		debugLogger.Printf("Error: requested: %s\n", r.URL.Path)
		http.Error(w, "Not found", 404)
		return
	}

	// fmt.Println("data", string(data))
	tmpl := template.Must(template.New(resourcefile).Parse(string(data)))
	if ext == "html" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	} else if ext == "js" {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	} else if ext == "css" {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}
	tmpl.Execute(w, r.Host)
}

func serveSvg(w http.ResponseWriter, r *http.Request) {
	var err error
	var offset int64
	var samples int64

	if r.Method != "GET" {
		http.Error(w, "Method nod allowed", 405)
		return
	}
	formVals, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	if offset, err = getFormValue(formVals, "offset", 0); err != nil {
		http.Error(w, "Bad Request. Could not convert parameters.", 400)
		return
	}

	if samples, err = getFormValue(formVals, "samples", 240); err != nil {
		http.Error(w, "Bad Request. Could not convert parameters.", 400)
		return
	}

	png := getPng(offset, samples)
	if png != nil {
		if *debug {
			// save the current image to a file.
			fmt.Println("writing debug.svg")
			ioutil.WriteFile("debug.svg", png, 0644)
		}
		// w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
		w.Header().Set("Content-Type", "image/png")
		w.Write(png)
	}
}
