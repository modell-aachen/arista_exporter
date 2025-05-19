package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type BgpVrf struct {
	RouterID string             `json:"routerId"`
	Asn      string             `json:"asn"`
	Peers    map[string]BgpPeer `json:"peers"`
	Vrf      string             `json:"vrf"`
}

type BgpPeer struct {
	Description         string        `json:"description"`
	MsgSent             int           `json:"msgSent"`
	MsgReceived         int           `json:"msgReceived"`
	PrefixReceived      int           `json:"prefixReceived"`
	PrefixAccepted      int           `json:"prefixAccepted"`
	PrefixInBest        int           `json:"prefixInBest"`
	PrefixInBestEcmp    int           `json:"prefixInBestEcmp"`
	InMsgQueue          int           `json:"inMsgQueue"`
	OutMsgQueue         int           `json:"outMsgQueue"`
	PeerState           string        `json:"peerState"`
	UpDownTime          float64       `json:"upDownTime"`
	Asn                 string        `json:"asn"`
	UnderMaintenance    bool          `json:"underMaintenance"`
	Version             int           `json:"version"`
	LldpNeighbors       []interface{} `json:"lldpNeighbors"`
	PeerStateIdleReason string        `json:"peerStateIdleReason,omitempty"`
}

type BgpCollector struct {
	Vrfs map[string]BgpVrf `json:"vrfs"`

	prefixReceivedGauge   *prometheus.GaugeVec
	prefixAcceptedGauge   *prometheus.GaugeVec
	prefixInBestGauge     *prometheus.GaugeVec
	prefixInBestEcmpGauge *prometheus.GaugeVec
	msgSentGauge          *prometheus.GaugeVec
	msgReceivedGauge      *prometheus.GaugeVec
	peerStateGauge        *prometheus.GaugeVec

	inMsgQueueGauge       *prometheus.GaugeVec
	outMsgQueueGauge      *prometheus.GaugeVec
	underMaintenanceGauge *prometheus.GaugeVec
}

var ifLabels = []string{"peer", "description", "asn", "vrf", "router_id"}
var bgpOpts = MakeSubsystemOptsFactory("bgp")

func (c *BgpCollector) GetCmd() string {
	return "show ipv6 bgp summary"
}

func (c *BgpCollector) Register(registry *prometheus.Registry) {
	c.prefixReceivedGauge = prometheus.NewGaugeVec(bgpOpts("prefix_received", "Number of prefixes received from BGP peer"), ifLabels)
	registry.MustRegister(c.prefixReceivedGauge)

	c.prefixAcceptedGauge = prometheus.NewGaugeVec(bgpOpts("prefix_accepted", "Number of prefixes accepted from BGP peer"), ifLabels)
	registry.MustRegister(c.prefixAcceptedGauge)

	c.prefixInBestGauge = prometheus.NewGaugeVec(bgpOpts("prefix_in_best", "Number of prefixes in best path from BGP peer"), ifLabels)
	registry.MustRegister(c.prefixInBestGauge)

	c.prefixInBestEcmpGauge = prometheus.NewGaugeVec(bgpOpts("prefix_in_best_ecmp", "Number of prefixes in best ECMP path from BGP peer"), ifLabels)
	registry.MustRegister(c.prefixInBestEcmpGauge)

	c.msgSentGauge = prometheus.NewGaugeVec(bgpOpts("msg_sent", "Number of BGP messages sent to peer"), ifLabels)
	registry.MustRegister(c.msgSentGauge)

	c.msgReceivedGauge = prometheus.NewGaugeVec(bgpOpts("msg_received", "Number of BGP messages received from peer"), ifLabels)
	registry.MustRegister(c.msgReceivedGauge)

	c.peerStateGauge = prometheus.NewGaugeVec(bgpOpts("peer_state", "BGP peer state: 1=Idle, 2=Connect, 3=Active, 4=OpenSent, 5=OpenConfirm, 6=Established"), ifLabels)
	registry.MustRegister(c.peerStateGauge)

	c.inMsgQueueGauge = prometheus.NewGaugeVec(bgpOpts("in_msg_queue", "Number of BGP messages in input queue"), ifLabels)
	registry.MustRegister(c.inMsgQueueGauge)

	c.outMsgQueueGauge = prometheus.NewGaugeVec(bgpOpts("out_msg_queue", "Number of BGP messages in output queue"), ifLabels)
	registry.MustRegister(c.outMsgQueueGauge)

	c.underMaintenanceGauge = prometheus.NewGaugeVec(bgpOpts("under_maintenance", "Whether the peer is under maintenance (1 if true, 0 if false)"), ifLabels)
	registry.MustRegister(c.underMaintenanceGauge)
}
func (c *BgpCollector) UpdateMetrics() {
	for vrfName, vrf := range c.Vrfs {
		for addr, peer := range vrf.Peers {
			var state float64
			switch peer.PeerState {
			case "Idle":
				state = 1
			case "Connect":
				state = 2
			case "Active":
				state = 3
			case "OpenSent":
				state = 4
			case "OpenConfirm":
				state = 5
			case "Established":
				state = 6
			default:
				state = 0
			}
			c.prefixReceivedGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(float64(peer.PrefixReceived))
			c.prefixAcceptedGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(float64(peer.PrefixAccepted))
			c.prefixInBestGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(float64(peer.PrefixInBest))
			c.prefixInBestEcmpGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(float64(peer.PrefixInBestEcmp))
			c.msgSentGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(float64(peer.MsgSent))
			c.msgReceivedGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(float64(peer.MsgReceived))
			c.peerStateGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(state)

			c.inMsgQueueGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(float64(peer.InMsgQueue))
			c.outMsgQueueGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(float64(peer.OutMsgQueue))
			maintenance := 0.0
			if peer.UnderMaintenance {
				maintenance = 1.0
			}
			c.underMaintenanceGauge.WithLabelValues(addr, peer.Description, peer.Asn, vrfName, vrf.RouterID).Set(maintenance)
		}
	}

}
