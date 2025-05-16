package collectors

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type PowerCollector struct {
	PowerSupplies    map[string]*PowerSupply `json:"powerSupplies"`
	psuMeta          *prometheus.GaugeVec
	psuUptime        *prometheus.GaugeVec
	psuCapacity      *prometheus.GaugeVec
	psuInputCurrent  *prometheus.GaugeVec
	psuOutputCurrent *prometheus.GaugeVec
	psuOutputPower   *prometheus.GaugeVec
}

type PowerSupply struct {
	OutputPower   float64                 `json:"outputPower"`
	State         string                  `json:"state"`
	ModelName     string                  `json:"modelName"`
	Capacity      int                     `json:"capacity"`
	InputCurrent  float64                 `json:"inputCurrent"`
	OutputCurrent float64                 `json:"outputCurrent"`
	Uptime        float64                 `json:"uptime"`
	Managed       bool                    `json:"managed"`
	TempSensors   map[string]TempSensor   `json:"tempSensors"`
	Fans          map[string]PSUFanStatus `json:"fans"`
}

type PSUFanStatus struct {
	Status string `json:"status"`
	Speed  int    `json:"speed"`
}

type TempSensor struct {
	Status      string `json:"status"`
	Temperature int    `json:"temperature"`
}

func (c *PowerCollector) GetCmd() string {
	return "show system environment power"
}

var psuOpts = MakeSubsystemOptsFactory("power_supply")

func (c *PowerCollector) Register(registry *prometheus.Registry) {
	// PSU Meta Gauge (serves to provide information that might be useful on dashboards
	c.psuMeta = prometheus.NewGaugeVec(psuOpts("meta", "Provides meta-info about each power supply unit. The gauge value can be 0 or 1, where 1 is set when the PSU State is OK"),
		[]string{"psuId", "model", "capacity", "managed"})
	defaultLabels := []string{"psuId"}
	// PSU Capacity Gauge
	c.psuCapacity = prometheus.NewGaugeVec(psuOpts("capacity", "Power Supply Capacity, in Watts"), defaultLabels)
	// PSU Uptime Gauge
	c.psuUptime = prometheus.NewGaugeVec(psuOpts("uptime", "PSU Uptime"),
		defaultLabels)
	// Input/Output current gauges
	c.psuInputCurrent = prometheus.NewGaugeVec(psuOpts("current_in", "Input Current from wall in AC amps"),
		defaultLabels)
	c.psuOutputCurrent = prometheus.NewGaugeVec(psuOpts("current_out", "Output Current to device in DC amps "),
		defaultLabels)
	c.psuOutputPower = prometheus.NewGaugeVec(psuOpts("power", "Power consumption in watts"),
		defaultLabels)

	// Register all metrics
	registry.MustRegister(c.psuMeta, c.psuUptime, c.psuInputCurrent, c.psuOutputCurrent, c.psuOutputPower)
}

func (c *PowerCollector) UpdateMetrics() {
	for id, psu := range c.PowerSupplies {
		// PSU Meta Info
		if psu.State == "ok" {
			c.psuMeta.WithLabelValues(id, psu.ModelName, strconv.Itoa(psu.Capacity), strconv.FormatBool(psu.Managed)).Set(1)
		} else {
			c.psuMeta.WithLabelValues(id, psu.ModelName, strconv.Itoa(psu.Capacity), strconv.FormatBool(psu.Managed)).Set(1)
		}

		// PSU Capacity
		c.psuCapacity.WithLabelValues(id).Set(float64(psu.Capacity))

		// PSU Uptime
		c.psuUptime.WithLabelValues(id).Set(psu.Uptime)

		// PSU Power Consumption
		c.psuInputCurrent.WithLabelValues(id).Set(psu.InputCurrent)
		c.psuOutputCurrent.WithLabelValues(id).Set(psu.OutputCurrent)
		c.psuOutputPower.WithLabelValues(id).Set(psu.OutputPower)
	}
}
