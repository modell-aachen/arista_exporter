package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

type CoolingCollector struct {
	OverrideFanSpeed           int       `json:"overrideFanSpeed"`
	CoolingMode                string    `json:"coolingMode"`
	ShutdownOnInsufficientFans bool      `json:"shutdownOnInsufficientFans"`
	AmbientTemperature         float64   `json:"ambientTemperature"`
	SystemStatus               string    `json:"systemStatus"`
	AirflowDirection           string    `json:"airflowDirection"`
	PowerSupplySlots           []FanSlot `json:"powerSupplySlots"`
	FanTraySlots               []FanSlot `json:"fanTraySlots"`

	coolingMeta             *prometheus.GaugeVec
	ambientTemperature      *prometheus.GaugeVec
	overrideFanSpeedGauge   *prometheus.GaugeVec
	fanMaxSpeedGauge        *prometheus.GaugeVec
	fanConfiguredSpeedGauge *prometheus.GaugeVec
	fanActualSpeedGauge     *prometheus.GaugeVec
}

type FanSlot struct {
	Status string      `json:"status"`
	Speed  int         `json:"speed"`
	Label  string      `json:"label"`
	Fans   []FanStatus `json:"fans"`
}

type FanStatus struct {
	Status                    string  `json:"status"`
	Uptime                    float64 `json:"uptime"`
	MaxSpeed                  int     `json:"maxSpeed"`
	ConfiguredSpeed           int     `json:"configuredSpeed"`
	ActualSpeed               int     `json:"actualSpeed"`
	SpeedStable               bool    `json:"speedStable"`
	LastSpeedStableChangeTime float64 `json:"lastSpeedStableChangeTime"`
	Label                     string  `json:"label"`
}

func (c *CoolingCollector) GetCmd() string {
	return "show environment cooling"
}

var coolingOpts = MakeSubsystemOptsFactory("cooling")

func (c *CoolingCollector) Register(registry *prometheus.Registry) {
	c.coolingMeta = prometheus.NewGaugeVec(coolingOpts("meta", "Metric containing meta information about the device's cooling system/config"),
		[]string{"coolingMode", "shutdownOnInsufficientFans", "systemStatus"})
	c.ambientTemperature = prometheus.NewGaugeVec(coolingOpts("ambient_temperature", "Ambient Temperature in Celsius"), []string{})
	c.overrideFanSpeedGauge = prometheus.NewGaugeVec(coolingOpts("override_fan_speed", "Fan speed override. If 0, the fan speed is not currently overridden"), []string{})

	labels := []string{"tray", "fanName"}
	c.fanMaxSpeedGauge = prometheus.NewGaugeVec(coolingOpts("fan_max_speed", "Maximum capable fan speed"), labels)
	c.fanConfiguredSpeedGauge = prometheus.NewGaugeVec(coolingOpts("fan_configured_speed", "Configured speed"), labels)
	c.fanActualSpeedGauge = prometheus.NewGaugeVec(coolingOpts("fan_actual_speed", "Actual speed"), labels)
	registry.MustRegister(c.fanMaxSpeedGauge)
}

func (c *CoolingCollector) UpdateMetrics() {
	c.coolingMeta.WithLabelValues(c.CoolingMode, strconv.FormatBool(c.ShutdownOnInsufficientFans), c.SystemStatus).Set(1)
	c.ambientTemperature.WithLabelValues().Set(c.AmbientTemperature)
	c.overrideFanSpeedGauge.WithLabelValues().Set(float64(c.OverrideFanSpeed))

	for _, slot := range c.PowerSupplySlots {
		for _, fan := range slot.Fans {
			c.fanMaxSpeedGauge.WithLabelValues("psu", fan.Label).Set(float64(fan.MaxSpeed))
			c.fanConfiguredSpeedGauge.WithLabelValues("psu", fan.Label).Set(float64(fan.ConfiguredSpeed))
			c.fanActualSpeedGauge.WithLabelValues("psu", fan.Label).Set(float64(fan.ActualSpeed))
		}
	}

	for _, slot := range c.FanTraySlots {
		for _, fan := range slot.Fans {
			c.fanMaxSpeedGauge.WithLabelValues("fanslot", fan.Label).Set(float64(fan.MaxSpeed))
			c.fanConfiguredSpeedGauge.WithLabelValues("fanslot", fan.Label).Set(float64(fan.ConfiguredSpeed))
			c.fanActualSpeedGauge.WithLabelValues("fanslot", fan.Label).Set(float64(fan.ActualSpeed))
		}
	}
}
