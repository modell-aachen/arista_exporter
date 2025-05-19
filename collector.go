package main

import (
	"strings"

	"github.com/modell-aachen/arista_exporter/collectors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type Collector interface {
	GetCmd() string
	Register(*prometheus.Registry)
	UpdateMetrics()
}

func getCollectorMap(enabled string) map[string]Collector {
	allCollectors := map[string]Collector{
		"version":     &collectors.VersionCollector{},
		"power":       &collectors.PowerCollector{},
		"interfaces":  &collectors.InterfacesCollector{},
		"cooling":     &collectors.CoolingCollector{},
		"temperature": &collectors.TemperatureCollector{},
		"bgp":         &collectors.BgpCollector{},
	}

	collectorMap := make(map[string]Collector)

	if enabled == "" {
		for name, coll := range allCollectors {
			collectorMap[name] = coll
		}
		return collectorMap
	}

	for _, name := range strings.Split(enabled, ",") {
		name = strings.TrimSpace(name)
		if coll, ok := allCollectors[name]; ok {
			collectorMap[name] = coll
		} else {
			log.Warnf("Unknown collector: %s", name)
		}
	}

	return collectorMap
}
