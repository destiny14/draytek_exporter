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
package vigorv5

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-kit/log/level"
	"github.com/tidwall/gjson"
)

var ErrUpdateFailed = errors.New("dsl status update failed")
var ErrParseFailed = errors.New("dsl status parse failed")

const (
	dslStatusGeneral = `{"param":[],"ct":[{"0MONITORING_DSL_GENERAL":[]},{"1MON_DSL_STREAM_TABLE":[]},{"1MON_DSL_END_TABLE":[]}]}`
)

type Status struct {
	Status     string
	Mode       string
	Profile    string
	Annex      string
	DSLVersion string

	ActualRateDownstream     int
	ActualRateUpstream       int
	AttainableRateDownstream int
	AttainableRateUpstream   int
	SNRMarginDownstream      float64
	SNRMarginUpstream        float64
	LineAttenuationNearEnd   float64
	LineAttenuationFarEnd    float64
	CRCNearEnd               int
	CRCFarEnd                int
	UASNearEnd               int
	UASFarEnd                int
	HECErrorsNearEnd         int
	HECErrorsFarEnd          int
	ESNearEnd                int
	ESFarEnd                 int
	SESNearEnd               int
	SESFarEnd                int
	LOSFailureNearEnd        int
	LOSFailureFarEnd         int
	LOFFailureNearEnd        int
	LOFFailureFarEnd         int
	LPRFailureNearEnd        int
	LPRFailureFarEnd         int
	LCDFailureNearEnd        int
	LCDFailureFarEnd         int
	RFECNearEnd              int
	RFECFarEnd               int
}

func (v *Vigor) FetchStatus() (Status, error) {
	post := vigorForm{
		pid: "0MONITORING_DSL_GENERAL",
		op:  "501",
		ct:  dslStatusGeneral,
	}

	resp, err := v.postWithLogin(post)
	if err != nil {
		level.Debug(v.logger).Log("msg", "Got error from post", "err", err)
		return Status{}, err
	}

	return v.parseDSLStatusGeneralJSON(resp)
}

func (v *Vigor) parseDSLStatusGeneralJSON(respJSON string) (Status, error) {
	value := gjson.Get(respJSON, "ct.0.0MONITORING_DSL_GENERAL.#(Name==\"Setting\")")
	if !value.Exists() {
		level.Debug(v.logger).Log("msg", "Unable to get settings", "response_json", respJSON)
		return Status{}, ErrParseFailed
	}

	level.Debug(v.logger).Log("msg", "Parsed DSL Status General json", "json", value.String())

	status := Status{
		Status:     value.Get("Status").String(),
		Mode:       value.Get("Mode").String(),
		Profile:    value.Get("Profile").String(),
		Annex:      value.Get("Annex").String(),
		DSLVersion: value.Get("DSL_Version").String(),
	}

	streamTable := value.Get("Stream_Table").Array()
	for _, v := range streamTable {
		switch v.Get("Name").String() {
		case "Actual Rate":
			status.ActualRateDownstream = parseKbps(v.Get("Downstream").String())
			status.ActualRateUpstream = parseKbps(v.Get("Upstream").String())
		case "Attainable Rate":
			status.AttainableRateDownstream = parseKbps(v.Get("Downstream").String())
			status.AttainableRateUpstream = parseKbps(v.Get("Upstream").String())
		case "SNR Margin":
			status.SNRMarginDownstream = parsedB(v.Get("Downstream").String())
			status.SNRMarginUpstream = parsedB(v.Get("Upstream").String())
		}
	}

	endTable := value.Get("End_Table").Array()
	for _, v := range endTable {
		switch v.Get("Name").String() {
		case "Attenuation":
			status.LineAttenuationNearEnd = parsedB(v.Get("Near_End").String())
			status.LineAttenuationFarEnd = parsedB(v.Get("Far_End").String())
		case "CRC":
			status.CRCNearEnd = int(v.Get("Near_End").Int())
			status.CRCFarEnd = int(v.Get("Far_End").Int())
		case "ES":
			status.ESNearEnd = int(v.Get("Near_End").Int())
			status.ESFarEnd = int(v.Get("Far_End").Int())
		case "SES":
			status.SESNearEnd = int(v.Get("Near_End").Int())
			status.SESFarEnd = int(v.Get("Far_End").Int())
		case "UAS":
			status.UASNearEnd = int(v.Get("Near_End").Int())
			status.UASFarEnd = int(v.Get("Far_End").Int())
		case "HEC Errors":
			status.HECErrorsNearEnd = int(v.Get("Near_End").Int())
			status.HECErrorsFarEnd = int(v.Get("Far_End").Int())
		case "LOS Failure":
			status.LOSFailureNearEnd = int(v.Get("Near_End").Int())
			status.LOSFailureFarEnd = int(v.Get("Far_End").Int())
		case "LOF Failure":
			status.LOFFailureNearEnd = int(v.Get("Near_End").Int())
			status.LOFFailureFarEnd = int(v.Get("Far_End").Int())
		case "LPR Failure":
			status.LPRFailureNearEnd = int(v.Get("Near_End").Int())
			status.LPRFailureFarEnd = int(v.Get("Far_End").Int())
		case "LCD Failure":
			status.LCDFailureNearEnd = int(v.Get("Near_End").Int())
			status.LCDFailureFarEnd = int(v.Get("Far_End").Int())
		case "RFEC":
			status.RFECNearEnd = int(v.Get("Near_End").Int())
			status.RFECFarEnd = int(v.Get("Far_End").Int())
		}
	}

	return status, nil
}

func parseKbps(s string) int {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0
	}
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return x * 1000
}

func parsedB(s string) float64 {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0
	}
	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}
	return x
}
