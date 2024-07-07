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

const (
        CmdDevStatus int = 1
	CmdPowerOn int = 2
	CmdPowerOff int = 3
)

var (
        Version = "unreleased"
        Quiet   = false
        Verbose = false
	configFile = "config.yaml"
	address string
	// Power On
	newTemperature string
	newMode string
	newFan string
	
	// daikinAcCtrlCmd represents the daikin-ac-ctrl command
	daikinAcCtrlCmd = &cobra.Command {
		Use:   "daikin-ac-ctrl",
		Short: "Control Daikin AC devices",
		Long: `Searches for Daikin Air Conditioners and manage them.`,
		Args:  cobra.ExactArgs(0),
	}
)

func init() {
        daikinAcCtrlCmd.Version = Version

	daikinAcCtrlCmd.PersistentFlags().StringVarP(&address, "address", "a", "", "Daikin aircon address")
	daikinAcCtrlCmd.PersistentFlags().StringVarP(&configFile, "config", "c", configFile, "configuration file")

	daikinAcCtrlCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", Quiet, "don't print any informative messages")
	daikinAcCtrlCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", Verbose, "become really verbose in printing messages")

	daikinAcCtrlCmd.AddCommand(
	        DevStatusCmd(),
		PowerOnCmd(),
		PowerOffCmd(),
	)
}

func DevStatusCmd() *cobra.Command {
        var subCmd = &cobra.Command {
                Use:   "status",
                Short: "Current status of daikin aircon",
                Run:   devStatus,
                Args:  cobra.ExactArgs(0),
        }

        return subCmd
}

func PowerOnCmd() *cobra.Command {
        var subCmd = &cobra.Command {
                Use:   "on",
                Short: "Power on daikin aircon",
                Run:   powerOn,
                Args:  cobra.ExactArgs(0),
        }

	subCmd.PersistentFlags().StringVarP(&newTemperature, "temperature", "t", "", "Target temperature")
	subCmd.PersistentFlags().StringVarP(&newMode, "mode", "m", "", "Operating mode (0=Auto, 2=Dehumidify, 3=Cool, 4=Heat, 6=Fan)")
	subCmd.PersistentFlags().StringVarP(&newFan, "fan", "f", "", "Fan speed (A=Auto, B=Silent, 3=Fan1, 4=Fan2, 5=Fan3, 6=Fan4, 7=Fan5)")

        return subCmd
}

func PowerOffCmd() *cobra.Command {
        var subCmd = &cobra.Command {
                Use:   "off",
                Short: "Power off daikin aircon",
                Run:   powerOff,
                Args:  cobra.ExactArgs(0),
        }

        return subCmd
}

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
	if err := daikinAcCtrlCmd.Execute(); err != nil {
                os.Exit(1)
        }
}

func devStatus(cmd *cobra.Command, args []string) {
        runDaikinAcCtrlCmd(CmdDevStatus)
}

func powerOn(cmd *cobra.Command, args []string) {
        runDaikinAcCtrlCmd(CmdPowerOn)
}

func powerOff(cmd *cobra.Command, args []string) {
        runDaikinAcCtrlCmd(CmdPowerOff)
}

func runDaikinAcCtrlCmd(cmd int) {

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
                log.Infof("Daikin AC Ctrl %s\n", Version)
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
		
		switch cmd {
    		case CmdDevStatus:
			fmt.Printf("Current %s:\n%s\n", target, d)
    		case CmdPowerOn:
			fmt.Printf("Switching %s on\n", target)
	             	d.ControlInfo.Power = daikin.PowerOn
			if len(newTemperature) > 0 {
			        if err := d.ControlInfo.Temperature.Set(newTemperature); err != nil {
			       	      log.Error(err)
				      os.Exit(1)
				}
			}
			if len(newMode) > 0 {
			   	if err := d.ControlInfo.Mode.Decode(newMode); err != nil {
			       	      log.Error(err)
				      os.Exit(1)
				}
			}
			if len(newFan) > 0 {
			   	if err := d.ControlInfo.Fan.Decode(newFan); err != nil {
			       	      log.Error(err)
				      os.Exit(1)
				}
			}
     	 	     	if err := d.SetControlInfo(); err != nil {
	       	     	       	log.Error(err)
               			os.Exit(1)
         		}
		case CmdPowerOff:
			fmt.Printf("Switching %s off\n", target)
	             	d.ControlInfo.Power = daikin.PowerOff
     	 	     	if err := d.SetControlInfo(); err != nil {
	       	     	       	log.Error(err)
               		        os.Exit(1)
         		}
    		}
	}
}
