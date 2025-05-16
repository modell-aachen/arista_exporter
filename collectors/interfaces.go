package collectors

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type InterfacesCollector struct {
	Interfaces map[string]Interface `json:"interfaces"`

	bandwidthGauge       *prometheus.GaugeVec
	interfaceStatusGauge *prometheus.GaugeVec

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

	// Additional counters
	totalOutErrorsGauge    *prometheus.GaugeVec
	totalInErrorsGauge     *prometheus.GaugeVec
	inTotalPacketsGauge    *prometheus.GaugeVec
	outTotalPacketsGauge   *prometheus.GaugeVec
	inDiscardsGauge        *prometheus.GaugeVec
	linkStatusChangesGauge *prometheus.GaugeVec

	// Interface statistics gauges
	inBitsRateGauge  *prometheus.GaugeVec
	inPktsRateGauge  *prometheus.GaugeVec
	outBitsRateGauge *prometheus.GaugeVec
	outPktsRateGauge *prometheus.GaugeVec

	inputRuntFramesGauge      *prometheus.GaugeVec
	inputRxPauseGauge         *prometheus.GaugeVec
	inputFcsErrorsGauge       *prometheus.GaugeVec
	inputAlignmentErrorsGauge *prometheus.GaugeVec
	inputGiantFramesGauge     *prometheus.GaugeVec
	inputSymbolErrorsGauge    *prometheus.GaugeVec
}

type Interface struct {
	Name                string  `json:"name"`
	InterfaceStatus     string  `json:"interfaceStatus"`
	LineProtocolStatus  string  `json:"lineProtocolStatus"`
	Hardware            string  `json:"hardware"`
	Bandwidth           int     `json:"bandwidth"`
	Description         string  `json:"description"`
	LastStatusChange    float64 `json:"lastStatusChangeTimestamp"`
	L2Mru               int     `json:"l2Mru"`
	Mtu                 int     `json:"mtu"`
	ForwardingModel     string  `json:"forwardingModel"`
	FallbackMode        string  `json:"fallbackMode"`
	FallbackTimeout     int     `json:"fallbackTimeout"`
	FallbackEnabled     bool    `json:"fallbackEnabled"`
	FallbackEnabledType string  `json:"fallbackEnabledType"`
	PhysicalAddress     string  `json:"physicalAddress"`

	InterfaceCounters struct {
		OutBroadcastPackets int `json:"outBroadcastPkts"`
		OutTotalPackets     int `json:"outTotalPkts"`
		OutUnicastPackets   int `json:"outUcastPkts"`
		TotalOutErrors      int `json:"totalOutErrors"`
		OutMulticastPackets int `json:"outMulticastPkts"`
		OutDiscards         int `json:"outDiscards"`
		OutOctets           int `json:"outOctets"`
		InBroadcastPackets  int `json:"inBroadcastPkts"`
		InMulticastPackets  int `json:"inMulticastPkts"`
		InUnicastPackets    int `json:"inUcastPkts"`
		InTotalPackets      int `json:"inTotalPkts"`
		InDiscards          int `json:"inDiscards"`
		InOctets            int `json:"inOctets"`
		TotalInErrors       int `json:"totalInErrors"`
		LinkStatusChanges   int `json:"linkStatusChanges"`

		InputErrorsDetail struct {
			RuntFrames      int `json:"runtFrames"`
			RxPause         int `json:"rxPause"`
			FcsErrors       int `json:"fcsErrors"`
			AlignmentErrors int `json:"alignmentErrors"`
			GiantFrames     int `json:"giantFrames"`
			SymbolErrors    int `json:"symbolErrors"`
		} `json:"inputErrorsDetail"`
	} `json:"interfaceCounters"`

	InterfaceStatistics struct {
		InBitsRate  float64 `json:"inBitsRate"`
		InPktsRate  float64 `json:"inPktsRate"`
		OutBitsRate float64 `json:"outBitsRate"`
		OutPktsRate float64 `json:"outPktsRate"`
	} `json:"interfaceStatistics"`
}

func (c *InterfacesCollector) GetCmd() string {
	return "show interfaces"
}

var interfacesOpts = MakeSubsystemOptsFactory("interface")

func (c *InterfacesCollector) Register(registry *prometheus.Registry) {
	ifLabels := []string{"interface", "part", "description", "physical_address"}

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

	// Additional counters
	c.totalOutErrorsGauge = prometheus.NewGaugeVec(
		interfacesOpts("errors_out", "Total outbound errors on the interface"), ifLabels)
	c.totalInErrorsGauge = prometheus.NewGaugeVec(
		interfacesOpts("errors_in", "Total inbound errors on the interface"), ifLabels)
	c.inTotalPacketsGauge = prometheus.NewGaugeVec(
		interfacesOpts("packets_in_total", "Total inbound packets on the interface"), ifLabels)
	c.outTotalPacketsGauge = prometheus.NewGaugeVec(
		interfacesOpts("packets_out_total", "Total outbound packets on the interface"), ifLabels)
	c.inDiscardsGauge = prometheus.NewGaugeVec(
		interfacesOpts("discards_in_total", "Total inbound discards on the interface"), ifLabels)
	c.linkStatusChangesGauge = prometheus.NewGaugeVec(
		interfacesOpts("link_changes", "Link status changes on the interface"), ifLabels)

	// Interface statistics gauges
	c.inBitsRateGauge = prometheus.NewGaugeVec(
		interfacesOpts("in_bits_rate", "Inbound bits rate on the interface"), ifLabels)
	c.inPktsRateGauge = prometheus.NewGaugeVec(
		interfacesOpts("in_pkts_rate", "Inbound packets rate on the interface"), ifLabels)
	c.outBitsRateGauge = prometheus.NewGaugeVec(
		interfacesOpts("out_bits_rate", "Outbound bits rate on the interface"), ifLabels)
	c.outPktsRateGauge = prometheus.NewGaugeVec(
		interfacesOpts("out_pkts_rate", "Outbound packets rate on the interface"), ifLabels)

	c.bandwidthGauge = prometheus.NewGaugeVec(
		interfacesOpts("bandwidth", "Interface bandwidth in bits per second"), ifLabels)

	c.interfaceStatusGauge = prometheus.NewGaugeVec(
		interfacesOpts("status", "Interface status: 1 if connected, 0 otherwise"), ifLabels)

	c.inputRuntFramesGauge = prometheus.NewGaugeVec(
		interfacesOpts("input_runt_frames", "Input runt frames on the interface"), ifLabels)
	c.inputRxPauseGauge = prometheus.NewGaugeVec(
		interfacesOpts("input_rx_pause", "Input RX pause frames on the interface"), ifLabels)
	c.inputFcsErrorsGauge = prometheus.NewGaugeVec(
		interfacesOpts("input_fcs_errors", "Input FCS errors on the interface"), ifLabels)
	c.inputAlignmentErrorsGauge = prometheus.NewGaugeVec(
		interfacesOpts("input_alignment_errors", "Input alignment errors on the interface"), ifLabels)
	c.inputGiantFramesGauge = prometheus.NewGaugeVec(
		interfacesOpts("input_giant_frames", "Input giant frames on the interface"), ifLabels)
	c.inputSymbolErrorsGauge = prometheus.NewGaugeVec(
		interfacesOpts("input_symbol_errors", "Input symbol errors on the interface"), ifLabels)

	// Register gauges
	registry.MustRegister(
		c.broadcastInGauge, c.unicastInGauge, c.multicastInGauge, c.discardsInGauge, c.octetsInGauge,
		c.broadcastOutGauge, c.unicastOutGauge, c.multicastOutGauge, c.discardsOutGauge, c.octetsOutGauge,
		c.totalOutErrorsGauge, c.totalInErrorsGauge, c.inTotalPacketsGauge, c.outTotalPacketsGauge,
		c.inDiscardsGauge, c.linkStatusChangesGauge,
		c.inBitsRateGauge, c.inPktsRateGauge, c.outBitsRateGauge, c.outPktsRateGauge,
		c.inputRuntFramesGauge, c.inputRxPauseGauge, c.inputFcsErrorsGauge, c.inputAlignmentErrorsGauge,
		c.inputGiantFramesGauge, c.inputSymbolErrorsGauge,
	)
	registry.MustRegister(c.bandwidthGauge)
	registry.MustRegister(c.interfaceStatusGauge)
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
		c.broadcastInGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InBroadcastPackets))
		c.unicastInGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InUnicastPackets))
		c.multicastInGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InMulticastPackets))
		c.discardsInGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InDiscards))
		c.octetsInGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InOctets))

		// Outbound gauges
		c.broadcastOutGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.OutBroadcastPackets))
		c.unicastOutGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.OutUnicastPackets))
		c.multicastOutGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.OutMulticastPackets))
		c.discardsOutGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.OutDiscards))
		c.octetsOutGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.OutOctets))

		// Counters
		c.totalOutErrorsGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.TotalOutErrors))
		c.totalInErrorsGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.TotalInErrors))
		c.inTotalPacketsGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InTotalPackets))
		c.outTotalPacketsGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.OutTotalPackets))
		c.inDiscardsGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InDiscards))
		c.linkStatusChangesGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.LinkStatusChanges))

		// Interface statistics metrics
		c.inBitsRateGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(iface.InterfaceStatistics.InBitsRate)
		c.inPktsRateGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(iface.InterfaceStatistics.InPktsRate)
		c.outBitsRateGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(iface.InterfaceStatistics.OutBitsRate)
		c.outPktsRateGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(iface.InterfaceStatistics.OutPktsRate)

		// Input errors detail metrics
		c.inputRuntFramesGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InputErrorsDetail.RuntFrames))
		c.inputRxPauseGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InputErrorsDetail.RxPause))
		c.inputFcsErrorsGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InputErrorsDetail.FcsErrors))

		c.inputAlignmentErrorsGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InputErrorsDetail.AlignmentErrors))
		c.inputGiantFramesGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InputErrorsDetail.GiantFrames))
		c.inputSymbolErrorsGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.InterfaceCounters.InputErrorsDetail.SymbolErrors))

		// Bandwidth and status
		c.bandwidthGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(float64(iface.Bandwidth))

		var status float64
		if iface.InterfaceStatus == "connected" {
			status = 1
		} else {
			status = 0
		}
		c.interfaceStatusGauge.WithLabelValues(ifName, ifPart, iface.Description, iface.PhysicalAddress).Set(status)
	}
}
