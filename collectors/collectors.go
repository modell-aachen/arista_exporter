package collectors

import "github.com/prometheus/client_golang/prometheus"

const Namespace = "arista"

func MakeSubsystemOptsFactory(Subsystem string) func(Name string, Help string) prometheus.GaugeOpts {
	return func(Name string, Help string) prometheus.GaugeOpts {
		return prometheus.GaugeOpts{
			Namespace: Namespace,
			Subsystem: Subsystem,
			Name:      Name,
			Help:      Help,
		}
	}
}
