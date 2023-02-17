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
	"errors"
	"fmt"
	"io/fs"
        "io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"net/http"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
	log "github.com/thkukuk/mqtt-exporter/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defListen = ":9071"
)

type ConfigType struct {
        Address string `yaml:"address,omitempty"`
        Listen  string `yaml:"listen"`
	Verbose bool   `yaml:"verbose"`
}

var (
        Version = "unreleased"
        Quiet   = false
        Verbose = false
	configFile = "config.yaml"
        address string
)

func read_yaml_config(conffile string) (ConfigType, error) {

        var config ConfigType

        file, err := ioutil.ReadFile(conffile)
        if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return config, nil
		} else {
			return config, fmt.Errorf("Cannot read %q: %v", conffile, err)
		}
        }
        err = yaml.Unmarshal(file, &config)
        if err != nil {
                return config, fmt.Errorf("Unmarshal error: %v", err)
        }

        return config, nil
}


func main() {
	// daikinAcExporterCmd represents the daikin-ac-exporter command
	daikinAcExporterCmd := &cobra.Command{
		Use:   "daikin-ac-exporter",
		Short: "Exports Daikin AC values as prometheus metrics",
		Long: `Starts a daemon which exports the values of the Daikin Air Conditioner as metrics for Proemtheus.`,
		Run: runDaikinAcExporterCmd,
		Args:  cobra.ExactArgs(0),
	}

        daikinAcExporterCmd.Version = Version

        daikinAcExporterCmd.Flags().StringVarP(&address, "address", "a", "", "Daikin aircon address")
	daikinAcExporterCmd.Flags().StringVarP(&configFile, "config", "c", configFile, "configuration file")

	daikinAcExporterCmd.Flags().BoolVarP(&Quiet, "quiet", "q", Quiet, "don't print any informative messages")
	daikinAcExporterCmd.Flags().BoolVarP(&Verbose, "verbose", "v", Verbose, "become really verbose in printing messages")

	if err := daikinAcExporterCmd.Execute(); err != nil {
                os.Exit(1)
        }
}

func runDaikinAcExporterCmd(cmd *cobra.Command, args []string) {
	var err error

	if !Quiet {
		log.Infof("Read yaml config %q\n", configFile)
	}
	config, err := read_yaml_config(configFile)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	if !Quiet {
                log.Infof("Daikin AC Exporter %s is starting...\n", Version)
        }

        quit := make(chan os.Signal, 1)
        signal.Notify(quit, os.Interrupt)
        signal.Notify(quit, syscall.SIGTERM)

        go func() {
                <-quit
                if !Quiet {
                        log.Info("Terminated via Signal. Shutting down...")
                }
                os.Exit(0)
        }()

	if len(config.Listen) == 0 {
		config.Listen = defListen
        }

	if len(address) == 0 && len(config.Address) > 0 {
		address = config.Address
	}

	if config.Verbose {
		Verbose = true
	}

	collector := newCollector(config)
        prometheus.MustRegister(collector)
        http.Handle("/metrics", promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
                // XXX ErrorLog: log,
        }))

	if !Quiet {
                log.Infof("Starting http server on %s", config.Listen)
        }
        log.Fatal(http.ListenAndServe(config.Listen, nil))
}
