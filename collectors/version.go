package collectors

import "github.com/prometheus/client_golang/prometheus"

type VersionCollector struct {
	Uptime           float64 `json:"uptime"`
	ModelName        string  `json:"modelName"`
	InternalVersion  string  `json:"internalVersion"`
	SystemMacAddress string  `json:"systemMacAddress"`
	SerialNumber     string  `json:"serialNumber"`
	BootupTimestamp  float64 `json:"bootupTimestamp"`
	MemoryTotal      int     `json:"memTotal"`
	MemoryFree       int     `json:"memFree"`
	Version          string  `json:"version"`
	Architecture     string  `json:"architecture"`
	IsIntlVersion    bool    `json:"isIntlVersion"`
	InternalBuildId  string  `json:"internalBuildId"`
	HardwareRevision string  `json:"hardwareRevision"`

	metaInfo    *prometheus.GaugeVec
	uptime      *prometheus.GaugeVec
	memoryTotal *prometheus.GaugeVec
	memoryFree  *prometheus.GaugeVec
}

func (c *VersionCollector) GetCmd() string {
	return "show version"
}

var versionOpts = MakeSubsystemOptsFactory("meta")

func (c *VersionCollector) Register(registry *prometheus.Registry) {
	// Metadata about the switch
	c.metaInfo = prometheus.NewGaugeVec(versionOpts("version", "Meta-info about this target"),
		[]string{"modelName", "systemMacAddress", "eosVersion", "serialNumber", "architecture", "hardwareRevision"})

	// Switch Uptime
	c.uptime = prometheus.NewGaugeVec(versionOpts("uptime", "Uptime"), []string{})

	// Switch Memory Consumption
	c.memoryTotal = prometheus.NewGaugeVec(versionOpts("memory_total", "Memory Total"), []string{})
	c.memoryFree = prometheus.NewGaugeVec(versionOpts("memory_free", "Memory Free"), []string{})

	registry.MustRegister(c.metaInfo, c.uptime, c.memoryTotal, c.memoryFree)
}

func (c *VersionCollector) UpdateMetrics() {
	// Record metadata
	c.metaInfo.WithLabelValues(c.ModelName, c.SystemMacAddress, c.Version, c.SerialNumber, c.Architecture, c.HardwareRevision).Set(1)
	// Record Uptime & Memory Consumption
	c.uptime.WithLabelValues().Set(c.Uptime)
	c.memoryTotal.WithLabelValues().Set(float64(c.MemoryTotal))
	c.memoryFree.WithLabelValues().Set(float64(c.MemoryFree))
}
