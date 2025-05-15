package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

type InterfacesCollector struct {
	Interfaces map[string]Interface `json:"interfaces"`

	broadcastInGauge  *prometheus.GaugeVec
	unicastInGauge    *prometheus.GaugeVec
	multicastInGauge  *prometheus.GaugeVec
	discardsInGauge   *prometheus.GaugeVec
	octetsInGauge     *prometheus.GaugeVec
	broadcastOutGauge *prometheus.GaugeVec
	unicastOutGauge   *prometheus.GaugeVec
	multicastOutGauge *prometheus.GaugeVec
	discardsOutGauge  *prometheus.GaugeVec
	octetsOutGauge    *prometheus.GaugeVec
}

type Interface struct {
	OutBroadcastPackets int     `json:"outBroadcastPkts"`
	OutUnicastPackets   int     `json:"outUcastPkts"`
	OutMulticastPackets int     `json:"outMulticastPkts"`
	OutDiscards         int     `json:"outDiscards"`
	OutOctets           int     `json:"outOctets"`
	InBroadcastPackets  int     `json:"inBroadcastPkts"`
	InUnicastPackets    int     `json:"inUcastPkts"`
	InMulticastPackets  int     `json:"inMulticastPkts"`
	InDiscards          int     `json:"inDiscards"`
	InOctets            int     `json:"inOctets"`
	LastUpdateTimestamp float64 `json:"lastUpdateTimestamp"`
}

func (c *InterfacesCollector) GetCmd() string {
	return "show interfaces counters"
}

var interfacesOpts = MakeSubsystemOptsFactory("interface")

func (c *InterfacesCollector) Register(registry *prometheus.Registry) {
	// the reason this is two labels is to handle how EOS exports interface information for things like QSFP+ ports (that are just 4x SFP+ lanes).
	// Where a regular 10G port will show up like "Ethernet6", a part of a QSFP+ port will show up as "Ethernet49/1"
	ifLabels := []string{"interface", "part"}

	// Inbound gauges
	c.broadcastInGauge = prometheus.NewGaugeVec(
		interfacesOpts("broadcast_in", "Inbound Broadcast Packets on the interface"), ifLabels)
	c.unicastInGauge = prometheus.NewGaugeVec(
		interfacesOpts("unicast_in", "Inbound Unicast Packets on the interface"), ifLabels)
	c.multicastInGauge = prometheus.NewGaugeVec(
		interfacesOpts("multicast_in", "Inbound Multicast Packets on the interface"), ifLabels)
	c.discardsInGauge = prometheus.NewGaugeVec(
		interfacesOpts("discards_in", "Inbound Discards on the interface"), ifLabels)
	c.octetsInGauge = prometheus.NewGaugeVec(
		interfacesOpts("octets_in", "Inbound Octets on the interface"), ifLabels)
	// Outbound gauges
	c.broadcastOutGauge = prometheus.NewGaugeVec(
		interfacesOpts("broadcast_out", "Outbound Broadcast Packets on the interface"), ifLabels)
	c.unicastOutGauge = prometheus.NewGaugeVec(
		interfacesOpts("unicast_out", "Outbound Unicast Packets on the interface"), ifLabels)
	c.multicastOutGauge = prometheus.NewGaugeVec(
		interfacesOpts("multicast_out", "Outbound Multicast Packets on the interface"), ifLabels)
	c.discardsOutGauge = prometheus.NewGaugeVec(
		interfacesOpts("discards_out", "Outbound Discards on the interface"), ifLabels)
	c.octetsOutGauge = prometheus.NewGaugeVec(
		interfacesOpts("octets_out", "Outbound Octets on the interface"), ifLabels)

	// Register gauges
	registry.MustRegister(c.broadcastInGauge, c.unicastInGauge, c.multicastInGauge, c.discardsInGauge, c.octetsInGauge, c.broadcastOutGauge, c.unicastOutGauge, c.multicastOutGauge, c.discardsOutGauge, c.octetsOutGauge)
}

func (c *InterfacesCollector) UpdateMetrics() {
	for name, iface := range c.Interfaces {
		nameParts := strings.Split(name, "/")
		ifName := nameParts[0]
		ifPart := "1"
		if len(nameParts) > 1 {
			ifPart = nameParts[1]
		}

		// Inbound gauges
		c.broadcastInGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.InBroadcastPackets))
		c.unicastInGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.InUnicastPackets))
		c.multicastInGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.InMulticastPackets))
		c.discardsInGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.InDiscards))
		c.octetsInGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.InOctets))
		// Outbound gauges
		c.broadcastOutGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.OutBroadcastPackets))
		c.unicastOutGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.OutUnicastPackets))
		c.multicastOutGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.OutMulticastPackets))
		c.discardsOutGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.OutDiscards))
		c.octetsOutGauge.WithLabelValues(ifName, ifPart).Set(float64(iface.OutOctets))

	}
}
