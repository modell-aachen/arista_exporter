package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"

	"github.com/alecthomas/kingpin"
	"github.com/aristanetworks/goeapi"
)

const (
	metricsPath = "/metrics"
)

var (
	configFile        = kingpin.Flag("config.file", "Arista exporter config file").Default(".eapi.conf").String()
	listenAddress     = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9465").String()
	enabledCollectors = kingpin.Flag("enabled-collectors", "Comma-separated list of collectors to enable. If empty, all are enabled.").Default("").String()
)

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}

	if !(slices.Contains(goeapi.Connections(), target)) {
		http.Error(w, "Target does not exist in config", http.StatusNotFound)
		return
	}

	log.Infof("Inbound request for target: %s", target)

	// Initialize eAPI Handle
	node, err := goeapi.ConnectTo(target)
	if err != nil {
		log.Errorf("Failed to connect to %q: %v", target, err)
		http.Error(w, fmt.Sprintf("Failed to connect to %q", target), http.StatusInternalServerError)
		return
	}

	eapiHandle, err := node.GetHandle("json")
	if err != nil {
		log.Errorf("Failed to get handle for %q: %v", target, err)
		http.Error(w, fmt.Sprintf("Failed to connect to %q", target), http.StatusInternalServerError)
		return
	}

	// Collectors
	collectorMap := getCollectorMap(*enabledCollectors)

	// Specific metrics registry to handle this request
	reg := prometheus.NewRegistry()

	// Register prometheus metrics and eAPI commands
	for name, coll := range collectorMap {
		coll.Register(reg)
		if aErr := eapiHandle.AddCommand(coll); aErr != nil {
			log.Fatalf("Failed to add command for collector %s", name)
			http.Error(w, fmt.Sprintf("Failed to add command for collector %s", name), http.StatusInternalServerError)
			return
		}
	}

	// Execute commands
	if cErr := eapiHandle.Call(); cErr != nil {
		http.Error(w, "Failed to run Arista eAPI Command", http.StatusInternalServerError)
		log.Errorf("Failed to run Arista eAPI Command: %v", cErr)
		return
	}
	log.Infof("Arista eAPI Command(s) ran successfully")

	// Update metrics
	for _, coll := range collectorMap {
		coll.UpdateMetrics()
	}
	log.Infof("Prometheus Metrics Updated")

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func startServer() {
	log.Infof("Starting arista exporter (Version: %s)", version.Print("arista_exporter"))
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Arista Exporter (Version ` + version.Print("arista_exporter") + `)</title></head>
			<body>
			<h1>Arista Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	http.HandleFunc(metricsPath, handleMetricsRequest)

	log.Infof("Listening for %s on %s", metricsPath, *listenAddress)

	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}

func main() {
	kingpin.Version(version.Print("arista_exporter"))
	kingpin.Parse()

	eapiAbsoluteConfigPath, err := filepath.Abs(*configFile)
	if err != nil {
		log.Fatalf("Invalid eapiConfig %s", *configFile)
	}

	log.Infoln("Loading configuration from", "path", eapiAbsoluteConfigPath)
	goeapi.LoadConfig(eapiAbsoluteConfigPath)

	log.Infoln("Valid Targets:", "targets", strings.Join(goeapi.Connections(), " "))

	startServer()
}
