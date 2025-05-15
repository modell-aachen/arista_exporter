package collectors

import "github.com/prometheus/client_golang/prometheus"

type TemperatureCollector struct {
	ShutdownOnOverheat bool                `json:"shutdownOnOverheat"`
	SystemStatus       string              `json:"systemStatus"`
	TemperatureSensors []TemperatureSensor `json:"tempSensors"`
	PowerSupplySlots   []PSUSlot           `json:"powerSupplySlots"`

	maxTempGauge               *prometheus.GaugeVec
	tempAlertCount             *prometheus.GaugeVec
	overheatThresholdTempGauge *prometheus.GaugeVec
	criticalThresholdTempGauge *prometheus.GaugeVec
	targetTempGauge            *prometheus.GaugeVec
	currentTempGauge           *prometheus.GaugeVec
}

type PSUSlot struct {
	ENTPhysicalClass   string              `json:"entPhysicalClass"`
	RelativePosition   string              `json:"relPos"`
	TemperatureSensors []TemperatureSensor `json:"tempSensors"`
}

type TemperatureSensor struct {
	MaxTemperature           float64 `json:"maxTemperature"`
	MaxTemperatureLastChange float64 `json:"maxTemperatureLastChange"`
	HwStatus                 string  `json:"hwStatus"`
	AlertCount               int     `json:"alertCount"`
	Description              string  `json:"description"`
	OverheatThreshold        int     `json:"overheatThreshold"`
	CriticalThreshold        int     `json:"criticalThreshold"`
	InAlertState             bool    `json:"inAlertState"`
	TargetTemperature        int     `json:"targetTemperature"`
	RelPos                   string  `json:"relPos"`
	CurrentTemperature       float64 `json:"currentTemperature"`
	PidDriverCount           int     `json:"pidDriverCount"`
	IsPidDriver              bool    `json:"isPidDriver"`
	Name                     string  `json:"name"`
}

func (c *TemperatureCollector) GetCmd() string {
	return "show environment temperature"
}

var tempOpts = MakeSubsystemOptsFactory("temperature")

func (c *TemperatureCollector) Register(registry *prometheus.Registry) {
	labels := []string{"sensorName", "sensorDescription"}
	c.maxTempGauge = prometheus.NewGaugeVec(tempOpts("max_temp", "The highest temperature that this sensor has hit"), labels)
	c.tempAlertCount = prometheus.NewGaugeVec(tempOpts("temp_alert", "Temperature Sensor Alert Count"), labels)
	c.overheatThresholdTempGauge = prometheus.NewGaugeVec(tempOpts("overheat_threshold", "Overheat Temperature Threshold"), labels)
	c.criticalThresholdTempGauge = prometheus.NewGaugeVec(tempOpts("critical_threshold", "Critical Temperature Threshold"), labels)
	c.targetTempGauge = prometheus.NewGaugeVec(tempOpts("target_temp", "Target Temperature"), labels)
	c.currentTempGauge = prometheus.NewGaugeVec(tempOpts("current_temp", "Current Temperature"), labels)

	registry.MustRegister(c.maxTempGauge, c.tempAlertCount, c.overheatThresholdTempGauge, c.criticalThresholdTempGauge, c.targetTempGauge, c.currentTempGauge)
}

func (c *TemperatureCollector) UpdateMetrics() {
	for _, sensor := range c.TemperatureSensors {
		labels := []string{sensor.Name, sensor.Description}
		c.maxTempGauge.WithLabelValues(labels...).Set(sensor.MaxTemperature)
		c.tempAlertCount.WithLabelValues(labels...).Set(float64(sensor.AlertCount))
		c.overheatThresholdTempGauge.WithLabelValues(labels...).Set(float64(sensor.OverheatThreshold))
		c.criticalThresholdTempGauge.WithLabelValues(labels...).Set(float64(sensor.CriticalThreshold))
		c.targetTempGauge.WithLabelValues(labels...).Set(float64(sensor.TargetTemperature))
		c.currentTempGauge.WithLabelValues(labels...).Set(float64(sensor.CurrentTemperature))
	}

	for _, psu := range c.PowerSupplySlots {
		for _, sensor := range psu.TemperatureSensors {
			labels := []string{sensor.Name, sensor.Description}
			c.maxTempGauge.WithLabelValues(labels...).Set(sensor.MaxTemperature)
			c.tempAlertCount.WithLabelValues(labels...).Set(float64(sensor.AlertCount))
			c.overheatThresholdTempGauge.WithLabelValues(labels...).Set(float64(sensor.OverheatThreshold))
			c.criticalThresholdTempGauge.WithLabelValues(labels...).Set(float64(sensor.CriticalThreshold))
			c.targetTempGauge.WithLabelValues(labels...).Set(float64(sensor.TargetTemperature))
			c.currentTempGauge.WithLabelValues(labels...).Set(float64(sensor.CurrentTemperature))
		}
	}
}
