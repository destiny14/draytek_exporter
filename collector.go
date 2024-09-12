// Copyright 2018 Ben Kochie <superq@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	vigorv5 "github.com/SuperQ/draytek_exporter/vigor_v5"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "draytek"

// Exporter collects Vigor stats from the given server and exports them using
// the prometheus metrics package.
type Exporter struct {
	v *vigorv5.Vigor
}

// NewExporter returns an initialized Exporter.
func NewExporter(v *vigorv5.Vigor) *Exporter {
	return &Exporter{v: v}
}

var (
	draytekUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the draytek instance status successful?",
		nil, nil,
	)
	draytekInfoDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "info"),
		"Info about the draytek router",
		[]string{"dsl_version", "mode", "profile", "status", "annex"}, nil,
	)

	actualRateDownDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "actual_bps"),
		"The actual downstream bits per second rate",
		nil, nil,
	)
	actualRateUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "actual_bps"),
		"The actual upstream bits per second rate",
		nil, nil,
	)
	attainableRateDownDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "attainable_bps"),
		"The attainable downstream bits per second rate",
		nil, nil,
	)
	attainableRateUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "attainable_bps"),
		"The attainable upstream bits per second rate",
		nil, nil,
	)
	snrMarginDownDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "downstream", "snr_margin_db"),
		"The downstream SNR margin in dB",
		nil, nil,
	)
	snrMarginUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "upstream", "snr_margin_db"),
		"The downstream SNR margin in dB",
		nil, nil,
	)
	lineAttenuationNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "line_attenuation_db"),
		"The near end line attenuation in dB",
		nil, nil,
	)
	lineAttenuationFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "line_attenuation_db"),
		"The far end line attenuation in dB",
		nil, nil,
	)
	crcNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "crc_errors"),
		"Number of CRC errors on near end",
		nil, nil,
	)
	crcFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "crc_errors"),
		"Number of CRC errors on far end",
		nil, nil,
	)
	uasNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "unavailable_seconds"),
		"Number of unavailable seconds on near end",
		nil, nil,
	)
	uasFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "unavailable_seconds"),
		"Number of unavailable seconds on far end",
		nil, nil,
	)
	hecErrorsNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "hec_errors"),
		"Number of HEC errors on near end",
		nil, nil,
	)
	hecErrorsFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "hec_errors"),
		"Number of HEC errors on far end",
		nil, nil,
	)
	esNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "errored_seconds"),
		"Number of errored seconds on near end",
		nil, nil,
	)
	esFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "errored_seconds"),
		"Number of errored seconds on far end",
		nil, nil,
	)
	sesNearEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "severely_errored_seconds"),
		"Number of severely errored seconds on near end",
		nil, nil,
	)
	sesFarEndDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "severely_errored_seconds"),
		"Number of severely errored seconds on far end",
		nil, nil,
	)
	LOSFailureNearEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "los_failure"),
		"Number of LOS failures on near end",
		nil, nil,
	)
	LOSFailureFarEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "los_failure"),
		"Number of LOS failures on far end",
		nil, nil,
	)
	LOFFailureNearEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "lof_failure"),
		"Number of LOF failures on near end",
		nil, nil,
	)
	LOFFailureFarEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "lof_failure"),
		"Number of LOF failures on far end",
		nil, nil,
	)
	LPRFailureNearEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "lpr_failure"),
		"Number of LPR failures on near end",
		nil, nil,
	)
	LPRFailureFarEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "lpr_failure"),
		"Number of LPR failures on far end",
		nil, nil,
	)
	LCDFailureNearEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "lcd_failure"),
		"Number of LCD failures on near end",
		nil, nil,
	)
	LCDFailureFarEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "lcd_failure"),
		"Number of LCD failures on far end",
		nil, nil,
	)
	RFECNearEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "near_end", "rfec"),
		"Number of RFEC bytes on near end",
		nil, nil,
	)
	RFECFarEnd = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "far_end", "rfec"),
		"Number of RFEC bytes on far end",
		nil, nil,
	)
)

// Describe describes all the metrics ever exported by the draytek_exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- draytekUpDesc
	ch <- draytekInfoDesc
	ch <- actualRateDownDesc
	ch <- actualRateUpDesc
	ch <- attainableRateDownDesc
	ch <- attainableRateUpDesc
	ch <- snrMarginDownDesc
	ch <- snrMarginUpDesc
	ch <- lineAttenuationNearEndDesc
	ch <- lineAttenuationFarEndDesc
	ch <- crcNearEndDesc
	ch <- crcFarEndDesc
	ch <- uasNearEndDesc
	ch <- uasFarEndDesc
	ch <- hecErrorsNearEndDesc
	ch <- hecErrorsFarEndDesc
	ch <- esNearEndDesc
	ch <- esFarEndDesc
	ch <- sesNearEndDesc
	ch <- sesFarEndDesc
	ch <- LOSFailureNearEnd
	ch <- LOSFailureFarEnd
	ch <- LOFFailureNearEnd
	ch <- LOFFailureFarEnd
	ch <- LPRFailureNearEnd
	ch <- LPRFailureFarEnd
	ch <- LCDFailureNearEnd
	ch <- LCDFailureFarEnd
	ch <- LCDFailureNearEnd
	ch <- LCDFailureFarEnd
	ch <- RFECNearEnd
	ch <- RFECFarEnd
}

// Collect fetches the stats from the draytek router and delivers them as
// Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	status, err := e.v.FetchStatus()
	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			draytekUpDesc, prometheus.GaugeValue, 0.0,
		)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		draytekUpDesc, prometheus.GaugeValue, 1.0,
	)

	ch <- prometheus.MustNewConstMetric(
		draytekInfoDesc, prometheus.GaugeValue, 1.0,
		status.DSLVersion, status.Mode, status.Profile, status.Status, status.Annex,
	)
	ch <- prometheus.MustNewConstMetric(
		actualRateDownDesc, prometheus.GaugeValue, float64(status.ActualRateDownstream),
	)
	ch <- prometheus.MustNewConstMetric(
		actualRateUpDesc, prometheus.GaugeValue, float64(status.ActualRateUpstream),
	)
	ch <- prometheus.MustNewConstMetric(
		attainableRateDownDesc, prometheus.GaugeValue, float64(status.AttainableRateDownstream),
	)
	ch <- prometheus.MustNewConstMetric(
		attainableRateUpDesc, prometheus.GaugeValue, float64(status.AttainableRateUpstream),
	)
	ch <- prometheus.MustNewConstMetric(
		snrMarginDownDesc, prometheus.GaugeValue, status.SNRMarginDownstream,
	)
	ch <- prometheus.MustNewConstMetric(
		snrMarginUpDesc, prometheus.GaugeValue, status.SNRMarginUpstream,
	)
	ch <- prometheus.MustNewConstMetric(
		lineAttenuationNearEndDesc, prometheus.GaugeValue, status.LineAttenuationNearEnd,
	)
	ch <- prometheus.MustNewConstMetric(
		lineAttenuationFarEndDesc, prometheus.GaugeValue, status.LineAttenuationFarEnd,
	)
	ch <- prometheus.MustNewConstMetric(
		crcNearEndDesc, prometheus.GaugeValue, float64(status.CRCNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		crcFarEndDesc, prometheus.GaugeValue, float64(status.CRCFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		uasNearEndDesc, prometheus.GaugeValue, float64(status.UASNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		uasFarEndDesc, prometheus.GaugeValue, float64(status.UASFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		hecErrorsNearEndDesc, prometheus.GaugeValue, float64(status.HECErrorsNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		hecErrorsFarEndDesc, prometheus.GaugeValue, float64(status.HECErrorsFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		esNearEndDesc, prometheus.GaugeValue, float64(status.ESNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		esFarEndDesc, prometheus.GaugeValue, float64(status.ESFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		sesNearEndDesc, prometheus.GaugeValue, float64(status.SESNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		sesFarEndDesc, prometheus.GaugeValue, float64(status.SESFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		LOSFailureNearEnd, prometheus.GaugeValue, float64(status.LOSFailureNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		LOSFailureFarEnd, prometheus.GaugeValue, float64(status.LOSFailureFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		LOFFailureNearEnd, prometheus.GaugeValue, float64(status.LOFFailureNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		LOFFailureFarEnd, prometheus.GaugeValue, float64(status.LOFFailureFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		LPRFailureNearEnd, prometheus.GaugeValue, float64(status.LPRFailureNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		LPRFailureFarEnd, prometheus.GaugeValue, float64(status.LPRFailureFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		LCDFailureNearEnd, prometheus.GaugeValue, float64(status.LCDFailureNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		LCDFailureFarEnd, prometheus.GaugeValue, float64(status.LCDFailureFarEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		RFECNearEnd, prometheus.GaugeValue, float64(status.RFECNearEnd),
	)
	ch <- prometheus.MustNewConstMetric(
		RFECFarEnd, prometheus.GaugeValue, float64(status.RFECFarEnd),
	)
}
