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

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
	"github.com/thkukuk/daikin-gomod/api"
	log "github.com/thkukuk/mqtt-exporter/pkg/logger"
)

type ConfigType struct {
	Address string `yaml:"address,omitempty"`
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
	// daikinAcInfoCmd represents the daikin-info command
	daikinAcInfoCmd := &cobra.Command{
		Use:   "daikin-ac-info",
		Short: "Prints all Daikin AC values",
		Long: `Searches for Daikin Air Conditioners and dumps all available infos.`,
		Run: runDaikinAcInfoCmd,
		Args:  cobra.ExactArgs(0),
	}

        daikinAcInfoCmd.Version = Version

	daikinAcInfoCmd.Flags().StringVarP(&address, "address", "a", "", "Daikin aircon address")
	daikinAcInfoCmd.Flags().StringVarP(&configFile, "config", "c", configFile, "configuration file")

	daikinAcInfoCmd.Flags().BoolVarP(&Quiet, "quiet", "q", Quiet, "don't print any informative messages")
	daikinAcInfoCmd.Flags().BoolVarP(&Verbose, "verbose", "v", Verbose, "become really verbose in printing messages")

	if err := daikinAcInfoCmd.Execute(); err != nil {
                os.Exit(1)
        }
}

func runDaikinAcInfoCmd(cmd *cobra.Command, args []string) {

	if !Quiet {
		log.Infof("Read yaml config %q\n", configFile)
	}
	config, err := read_yaml_config(configFile)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

        if len(address) == 0 && len(config.Address) > 0 {
                address = config.Address
        }

	if !Quiet {
                log.Infof("Daikin AC Info %s\n", Version)
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

	d, err := daikin.NewNetwork(daikin.DebugOption(Verbose),
		daikin.AddressOption(address))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
        if err = d.Discover(); err != nil {
		log.Fatalf("Discover Error: %v", err)
        }

	for target, d := range d.Devices {

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
		fmt.Printf("Current %s:\n%s\n", target, d)
	}
}
