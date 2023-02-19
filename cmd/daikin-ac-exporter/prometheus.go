// Copyright 2023 Thorsten Kukuk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	log "github.com/thkukuk/mqtt-exporter/pkg/logger"
	"github.com/thkukuk/daikin-gomod/api"
	"github.com/prometheus/client_golang/prometheus"
)

var (
        namespace = "daikin_ac"

        device_info = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "device_info"),
                "device info, name and firmware",
                []string{"target", "type", "name", "version", "revision"}, nil,
        )

        htemp = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "htemp"),
                "sensor info, inside temperature (htemp)",
                []string{"target"}, nil,
        )

        hhum = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "hhum"),
                "sensor info, inside humidity (hhum)",
                []string{"target"}, nil,
        )

        otemp = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "otemp"),
                "sensor info, outside temperature (otemp)",
                []string{"target"}, nil,
        )

        er = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "err"),
                "sensor info, err",
                []string{"target"}, nil,
        )

        cmpfreq = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "cmpfreq"),
                "sensor info, cmpfreq",
                []string{"target"}, nil,
        )

        mompow = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "mompow"),
                "sensor info, mompow",
                []string{"target"}, nil,
        )

        filter_sign = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "filter_sign"),
                "sensor info, filter_sign",
                []string{"target"}, nil,
        )

        pow = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "pow"),
                "control info, power (pow)",
                []string{"target"}, nil,
        )

        mode = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "mode"),
                "control info, mode",
                []string{"target"}, nil,
        )

        adv = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "adv"),
                "control info, adv",
                []string{"target"}, nil,
        )

        stemp = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "stemp"),
                "control info, target temperature (stemp)",
                []string{"target"}, nil,
        )

        shum = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "shum"),
                "control info, target humidity (shum)",
                []string{"target"}, nil,
        )

        dt1 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dt1"),
                "control info, dt1",
                []string{"target"}, nil,
        )

        dt2 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dt2"),
                "control info, dt2",
                []string{"target"}, nil,
        )

        dt3 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dt3"),
                "control info, dt3",
                []string{"target"}, nil,
        )

        dt4 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dt4"),
                "control info, dt4",
                []string{"target"}, nil,
        )

        dt5 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dt5"),
                "control info, dt5",
                []string{"target"}, nil,
        )

        dt7 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dt7"),
                "control info, dt7",
                []string{"target"}, nil,
        )

        dh1 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dh1"),
                "control info, dh1",
                []string{"target"}, nil,
        )

	dh2 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dh2"),
                "control info, dh2",
                []string{"target"}, nil,
        )

        dh3 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dh3"),
                "control info, dh3",
                []string{"target"}, nil,
        )

        dh4 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dh4"),
                "control info, dh4",
                []string{"target"}, nil,
        )

        dh5 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dh5"),
                "control info, dh5",
                []string{"target"}, nil,
        )

        dh7 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dh7"),
                "control info, dh7",
                []string{"target"}, nil,
        )

        dhh = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dhh"),
                "control info, dhh",
                []string{"target"}, nil,
        )

        b_mode = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "b_mode"),
                "control info, b_mode",
                []string{"target"}, nil,
        )

        b_stemp = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "b_stemp"),
                "control info, b_stemp",
                []string{"target"}, nil,
        )

        b_shum = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "b_shum"),
                "control info, b_shum",
                []string{"target"}, nil,
        )

        alert = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "alert"),
                "control info, alert",
                []string{"target"}, nil,
        )

        f_rate = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "f_rate"),
                "control info, fan rate mode (f_rate)",
                []string{"target"}, nil,
        )

        b_f_rate = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "b_f_rate"),
                "control info, b_f_rate",
                []string{"target"}, nil,
        )

        dfr1 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfr1"),
                "control info, dfr1",
                []string{"target"}, nil,
        )

        dfr2 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfr2"),
                "control info, dfr2",
                []string{"target"}, nil,
        )

        dfr3 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfr3"),
                "control info, dfr3",
                []string{"target"}, nil,
        )

        dfr4 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfr4"),
                "control info, dfr4",
                []string{"target"}, nil,
        )

        dfr5 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfr5"),
                "control info, dfr5",
                []string{"target"}, nil,
        )

        dfr6 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfr6"),
                "control info, dfr6",
                []string{"target"}, nil,
        )

        dfr7 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfr7"),
                "control info, dfr7",
                []string{"target"}, nil,
        )

        dfrh = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfrh"),
                "control info, dfrh",
                []string{"target"}, nil,
        )

        f_dir = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "f_dir"),
                "control info, fan direction (f_dir)",
                []string{"target"}, nil,
        )

        b_f_dir = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "b_f_dir"),
                "control info, b_f_dir",
                []string{"target"}, nil,
        )

        dfd1 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfd1"),
                "control info, dfd1",
                []string{"target"}, nil,
        )

        dfd2 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfd2"),
                "control info, dfd2",
                []string{"target"}, nil,
        )

        dfd3 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfd3"),
                "control info, dfd3",
                []string{"target"}, nil,
        )

        dfd4 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfd4"),
                "control info, dfd4",
                []string{"target"}, nil,
        )

        dfd5 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfd5"),
                "control info, dfd5",
                []string{"target"}, nil,
        )

        dfd6 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfd6"),
                "control info, dfd6",
                []string{"target"}, nil,
        )

        dfd7 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfd7"),
                "control info, dfd7",
                []string{"target"}, nil,
        )

        dfdh = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dfdh"),
                "control info, dfdh",
                []string{"target"}, nil,
        )

        stemp_a = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "stemp_a"),
                "control info, stemp_a",
                []string{"target"}, nil,
        )

        dt1_a = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dt1_a"),
                "control info, dt1_a",
                []string{"target"}, nil,
        )

        dt7_a = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "dt7_a"),
                "control info, dt7_a",
                []string{"target"}, nil,
        )

        b_stemp_a = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "b_stemp_a"),
                "control info, b_stemp_a",
                []string{"target"}, nil,
        )

        f_dir_ud = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "f_dir_ud"),
                "control info, f_dir_ud",
                []string{"target"}, nil,
        )

        f_dir_lr = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "f_dir_lr"),
                "control info, f_dir_lr",
                []string{"target"}, nil,
        )

        b_f_dir_ud = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "b_f_dir_ud"),
                "control info, b_f_dir_ud",
                []string{"target"}, nil,
        )

        b_f_dir_lr = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "b_f_dir_lr"),
                "control info, b_f_dir_lr",
                []string{"target"}, nil,
        )

        ndfd1 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "ndfd1"),
                "control info, ndfd1",
                []string{"target"}, nil,
        )

        ndfd2 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "ndfd2"),
                "control info, ndfd2",
                []string{"target"}, nil,
        )

        ndfd3 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "ndfd3"),
                "control info, ndfd3",
                []string{"target"}, nil,
        )

        ndfd4 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "ndfd4"),
                "control info, ndfd4",
                []string{"target"}, nil,
        )

        ndfd5 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "ndfd5"),
                "control info, ndfd5",
                []string{"target"}, nil,
        )

        ndfd6 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "ndfd6"),
                "control info, ndfd6",
                []string{"target"}, nil,
        )

        ndfd7 = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "ndfd7"),
                "control info, ndfd7",
                []string{"target"}, nil,
        )

        ndfdh = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "ndfdh"),
                "control info, ndfdh",
                []string{"target"}, nil,
        )

        curr_day_heat = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "curr_day_heat"),
                "power info, power consumption heating",
                []string{"target"}, nil,
        )

        curr_day_cool = prometheus.NewDesc(
                prometheus.BuildFQName(namespace, "", "curr_day_cool"),
                "power info, power consumption cooling",
                []string{"target"}, nil,
        )
)

type Collector struct {
	Devices *daikin.DaikinNetwork
}

func newCollector(config ConfigType) *Collector {
	if Verbose {
		log.Debug("Creating prometheus collector...")
	}

	// XXX return error, don't abort
	d, err := daikin.NewNetwork(daikin.DebugOption(Verbose),
				    daikin.AddressOption(address))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
        if err = d.Discover(); err != nil {
		log.Fatalf("Discover Error: %v", err)
        }

	return &Collector{
		Devices: d,
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- device_info
        ch <- htemp
        ch <- hhum
        ch <- otemp
        ch <- er
        ch <- cmpfreq
        ch <- mompow
        ch <- filter_sign
        ch <- pow
        ch <- mode
        ch <- adv
        ch <- stemp
        ch <- shum
        ch <- dt1
        ch <- dt2
        ch <- dt3
        ch <- dt4
        ch <- dt5
        ch <- dt7
        ch <- dh1
        ch <- dh2
        ch <- dh3
        ch <- dh4
        ch <- dh5
        ch <- dh7
        ch <- dhh
        ch <- b_mode
        ch <- b_stemp
        ch <- b_shum
        ch <- alert
        ch <- f_rate
        ch <- b_f_rate
        ch <- dfr1
        ch <- dfr2
        ch <- dfr3
        ch <- dfr4
        ch <- dfr5
        ch <- dfr6
        ch <- dfr7
        ch <- dfrh
        ch <- f_dir
        ch <- b_f_dir
        ch <- dfd1
        ch <- dfd2
        ch <- dfd3
        ch <- dfd4
        ch <- dfd5
        ch <- dfd6
        ch <- dfd7
        ch <- dfdh
        ch <- stemp_a
        ch <- dt1_a
        ch <- dt7_a
        ch <- b_stemp_a
        ch <- f_dir_ud
        ch <- f_dir_lr
        ch <- b_f_dir_ud
        ch <- b_f_dir_lr
        ch <- ndfd1
        ch <- ndfd2
        ch <- ndfd3
        ch <- ndfd4
        ch <- ndfd5
        ch <- ndfd6
        ch <- ndfd7
        ch <- ndfdh
	ch <- curr_day_heat
	ch <- curr_day_cool
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {

        for target, d := range c.Devices.Devices {

                if err := d.GetBasicInfo(); err != nil {
                        log.Error(err)
                        continue
                }
                if err := d.GetControlInfo(); err != nil {
                        log.Error(err)
                        continue
                }
                if err := d.GetSensorInfo(); err != nil {
                        log.Error(err)
                        continue
                }
                if err := d.GetPowerInfo(); err != nil {
                        log.Error(err)
                        continue
                }
		if Verbose {
			log.Debugf("Current %s:\n%s\n\n", target, d)
		}

		// Device Info
		ch <- prometheus.MustNewConstMetric(device_info, prometheus.GaugeValue, 0, target, d.BasicInfo.Type.String(), d.BasicInfo.Name.String(), d.BasicInfo.Version.String(), d.BasicInfo.Revision.String())

		// Sensor Info
		ch <- prometheus.MustNewConstMetric(htemp, prometheus.GaugeValue, d.SensorInfo.HomeTemperature.Float64(), target)
		ch <- prometheus.MustNewConstMetric(hhum, prometheus.GaugeValue, d.SensorInfo.Humidity.Float64(), target)
		ch <- prometheus.MustNewConstMetric(otemp, prometheus.GaugeValue, d.SensorInfo.OutsideTemperature.Float64(), target)

		// Control Info
		ch <- prometheus.MustNewConstMetric(pow, prometheus.GaugeValue, d.ControlInfo.Power.Float64(), target)
		ch <- prometheus.MustNewConstMetric(mode, prometheus.GaugeValue, d.ControlInfo.Mode.Float64(), target)
		ch <- prometheus.MustNewConstMetric(stemp, prometheus.GaugeValue, d.ControlInfo.Temperature.Float64(), target)
		ch <- prometheus.MustNewConstMetric(shum, prometheus.GaugeValue, d.ControlInfo.Humidity.Float64(), target)
		ch <- prometheus.MustNewConstMetric(f_rate, prometheus.GaugeValue, d.ControlInfo.Fan.Float64(), target)
		ch <- prometheus.MustNewConstMetric(f_dir, prometheus.GaugeValue, d.ControlInfo.FanDir.Float64(), target)

		// Power Info
		ch <- prometheus.MustNewConstMetric(curr_day_cool, prometheus.GaugeValue, d.PowerInfo.DayCool.Float64(), target)
		ch <- prometheus.MustNewConstMetric(curr_day_heat, prometheus.GaugeValue, d.PowerInfo.DayHeat.Float64(), target)
	}
}
