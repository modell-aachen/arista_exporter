package collectors

import (
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
)

type BgpPeer struct {
	Description      string  `json:"description"`
	MsgSent          int     `json:"msgSent"`
	MsgReceived      int     `json:"msgReceived"`
	PrefixReceived   int     `json:"prefixReceived"`
	PrefixAccepted   int     `json:"prefixAccepted"`
	PrefixInBest     int     `json:"prefixInBest"`
	PrefixInBestEcmp int     `json:"prefixInBestEcmp"`
	InMsgQueue       int     `json:"inMsgQueue"`
	OutMsgQueue      int     `json:"outMsgQueue"`
	PeerState        string  `json:"peerState"`
	UpDownTime       float64 `json:"upDownTime"`
	Asn              string  `json:"asn"`
}

type BgpCollector struct {
	Peers map[string]BgpPeer `json:"peers"`

	prefixReceivedGauge   *prometheus.GaugeVec
	prefixAcceptedGauge   *prometheus.GaugeVec
	prefixInBestGauge     *prometheus.GaugeVec
	prefixInBestEcmpGauge *prometheus.GaugeVec
	msgSentGauge          *prometheus.GaugeVec
	msgReceivedGauge      *prometheus.GaugeVec
	peerStateGauge        *prometheus.GaugeVec
}

var ifLabels = []string{"peer", "description", "asn"}
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

	c.peerStateGauge = prometheus.NewGaugeVec(bgpOpts("peer_state", "BGP peer state (1 if Established, 0 otherwise)"), ifLabels)
	registry.MustRegister(c.peerStateGauge)
}

func (c *BgpCollector) UpdateMetrics() {
	log.Infof("Updating BGP metrics: %s", c)
	for addr, peer := range c.Peers {
		state := 0.0
		if peer.PeerState == "Established" {
			state = 1.0
		}
		c.prefixReceivedGauge.WithLabelValues(addr, peer.Description, peer.Asn).Set(float64(peer.PrefixReceived))
		c.prefixAcceptedGauge.WithLabelValues(addr, peer.Description, peer.Asn).Set(float64(peer.PrefixAccepted))
		c.prefixInBestGauge.WithLabelValues(addr, peer.Description, peer.Asn).Set(float64(peer.PrefixInBest))
		c.prefixInBestEcmpGauge.WithLabelValues(addr, peer.Description, peer.Asn).Set(float64(peer.PrefixInBestEcmp))
		c.msgSentGauge.WithLabelValues(addr, peer.Description, peer.Asn).Set(float64(peer.MsgSent))
		c.msgReceivedGauge.WithLabelValues(addr, peer.Description, peer.Asn).Set(float64(peer.MsgReceived))
		c.peerStateGauge.WithLabelValues(addr, peer.Description, peer.Asn).Set(state)
	}
}
