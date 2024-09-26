package collector

import (
	"strings"

	"github.com/EncoreTechnologies/netapp-api-exporter/pkg/netapp"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type SystemCollector struct {
	filerName   string
	versionDesc *prometheus.Desc
	client      *netapp.Client
}

func NewSystemCollector(client *netapp.Client, filerName string) *SystemCollector {
	return &SystemCollector{
		filerName: filerName,
		client:    client,
		versionDesc: prometheus.NewDesc(
			"netapp_filer_system_version",
			"Info about ontap version in labels `version` and `full_version`",
			[]string{"full_version", "version"},
			nil,
		),
	}
}

func (c *SystemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.versionDesc
}

func (c *SystemCollector) Collect(ch chan<- prometheus.Metric) {
	fullVersion, err := c.client.GetSystemVersion()
	if err != nil {
		log.WithError(err).Error("get system version failed")
		return
	}
	idx := strings.Index(fullVersion, ":")
	if idx == -1 {
		log.Warnf("[%s] Failed to extract version from string %q", c.filerName, fullVersion)
		return
	}
	version := fullVersion[:idx]
	ch <- prometheus.MustNewConstMetric(
		c.versionDesc,
		prometheus.GaugeValue,
		0.0,
		fullVersion,
		version,
	)
}
