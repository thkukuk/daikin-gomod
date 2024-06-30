// Package daikin provides functionality to interact with Daikin split
// system air conditioners equipped with a Wifi module. It is tested to work
// with the BRP072A42 Wifi interface.
package daikin

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	uriGetBasicInfo    = "/common/basic_info"
	uriGetRemoteMethod = "/common/get_remote_method"
	uriGetModelInfo    = "/aircon/get_model_info"
	uriGetControlInfo  = "/aircon/get_control_info"
	uriGetSensorInfo   = "/aircon/get_sensor_info"
	uriGetTimer        = "/aircon/get_timer"
	uriGetPrice        = "/aircon/get_price"
	uriGetTarget       = "/aircon/get_target"
	uriGetDayPowerEx   = "/aircon/get_day_power_ex"
	uriGetWeekPower    = "/aircon/get_week_power"
	uriGetYearPower    = "/aircon/get_year_power"
	uriGetProgram      = "/aircon/get_program"
	uriGetScdlTimer    = "/aircon/get_scdltimer"
	uriGetNotify       = "/aircon/get_notify"
	uriSetControlInfo  = "/aircon/set_control_info"
)

/*
type Parameter interface {
	// String with name and value for URL
	setUrlValues() string
	// Set sets this parameter's value.
	Set(string) error
	// String returns the human readable value.
	String() string
        // Float64 returns the value as float64 for prometheus.
        Float64() float64
}
*/

const (
	returnOk  = "OK"
	returnBad = "PARAM NG"
)

// Daikin represents the settings of the Daikin unit.
type Daikin struct {
	// Address is the IP address of the unit.
	Address string
	// BasicInfo contains the environment basic info.
	BasicInfo *BasicInfo
	// ControlInfo contains the environment control info.
	ControlInfo *ControlInfo
	// SensorInfo contains the environment sensor info.
	SensorInfo *SensorInfo
	// Power consumption heating and cooling
	PowerInfo *PowerInfo
}

// BasicInfo represents basic informations about the device
type BasicInfo struct {
	// Name is the human-readable name of the unit.
	Name Name
	// Version is the firmware version
	Version Version
	// Revision
	Revision String
	// Type: aircon
	Type String
}

func (b *BasicInfo) populate(values map[string]string) error {
	for k, v := range values {
		var err error
		switch k {
		case "name":
			err = b.Name.decode("name", v)
		case "ver":
			err = b.Version.decode("ver", v)
		case "rev":
			err = b.Revision.decode("rev", v)
		case "type":
			err = b.Type.decode("type", v)
		case "ret":
			if v != returnOk {
				err = fmt.Errorf("device returned error ret=%s", v)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BasicInfo) String() string {
	return fmt.Sprintf("Name: %s\nType: %s\nFirmware Version: %s\nRevision: %s",
		b.Name.String(), b.Type.String(), b.Version.String(), b.Revision.String())
}


// SensorInfo represents current sensor values.
type SensorInfo struct {
	// HomeTemperature is the home (interior) temperature.
	HomeTemperature Temperature
	// OutsideTemperature is the external temperature.
	OutsideTemperature Temperature
	// Humidity is the current interior humidity.
	Humidity Humidity
}

func (s *SensorInfo) populate(values map[string]string) error {
	for k, v := range values {
		var err error
		switch k {
		case "htemp":
			err = s.HomeTemperature.decode("htemp", v)
		case "otemp":
			err = s.OutsideTemperature.decode("otemp", v)
		case "hhum":
			err = s.Humidity.decode("hhum", v)
		case "ret":
			if v != returnOk {
				err = fmt.Errorf("device returned error ret=%s", v)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SensorInfo) String() string {
	return fmt.Sprintf("Inside temperature: %s\nInside humidity: %s\nOutside temperature: %s", s.HomeTemperature.String(), s.Humidity.String(), s.OutsideTemperature.String())
}

// ControlInfo represents the control status of the unit.
type ControlInfo struct {
	// Power is the current power status of the unit.
	Power Power
	// Mode is the operating mode of the unit.
	Mode Mode
	// Fan is the fan speed of the unit.
	Fan Fan
	// FanDir is the fan louvre setting of the unit.
	FanDir FanDir
	// Temperature is the current set temperature of the unit.
	Temperature Temperature
	// Humidity is the set humidity of the unit.
	Humidity Humidity
}

func (c *ControlInfo) urlValues() string {
	values := c.Power.setUrlValues();
	values = values + "&" + c.Mode.setUrlValues()
	values = values + "&" + c.Fan.setUrlValues()
	values = values + "&" + c.FanDir.setUrlValues()
	values = values + "&" + c.Temperature.setUrlValues()
	values = values + "&" + c.Humidity.setUrlValues()
	return values
}

func (c *ControlInfo) populate(values map[string]string) error {
	for k, v := range values {
		var err error
		switch k {
		case "pow":
			err = c.Power.decode(v)
		case "mode":
			err = c.Mode.decode(v)
		case "stemp":
			err = c.Temperature.decode("stemp", v)
		case "shum":
			err = c.Humidity.decode("shum", v)
		case "f_rate":
			err = c.Fan.decode(v)
		case "f_dir":
			err = c.FanDir.decode(v)
		case "ret":
			if v != returnOk {
				err = fmt.Errorf("device returned error ret=%s", v)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ControlInfo) String() string {
	return fmt.Sprintf("Power: %s\nMode: %s\nSet temperature: %s\nSet humidity: %s\nFan speed: %s\nFan louvre: %s",
		c.Power.String(), c.Mode.String(), c.Temperature.String(), c.Humidity.String(), c.Fan.String(), c.FanDir.String())
}

// PowerInfo represents power usage over the current day
type PowerInfo struct {
	DayHeat KWattHours
	DayCool KWattHours
}

// ret=OK,curr_day_heat=0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0,prev_1day_heat=0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0,curr_day_cool=0/1/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0,prev_1day_cool=0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0/0
func (w *PowerInfo) populate(values map[string]string) error {
	for k, v := range values {
		var err error
		switch k {
		case "curr_day_heat":
		case "curr_day_cool":
			elems := strings.Split(v, "/")
			if len(elems) != 24 {
				return fmt.Errorf("expected 24 elements in day power data, got %d", len(elems))
			}
			var total int64;
			for i := 0; i < 24; i++ {
				n, err := strconv.Atoi(elems[i]);
				if err != nil {
					return fmt.Errorf("error parsing day power data[i]=%s: %v", i, elems[i], err)
				}
				total = total + int64(n);
			}
			// data is 0.1 kWh
			kWh := float64(total) / 10.0;
			if k == "curr_day_heat" {
				w.DayHeat.decode(k, strconv.FormatFloat(kWh, 'f', 1, 64))
			} else {
				w.DayCool.decode(k, strconv.FormatFloat(kWh, 'f', 1, 64))
			}
		case "ret":
			if v != returnOk {
				err = fmt.Errorf("device returned error ret=%s", v)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *PowerInfo) String() string {
	return fmt.Sprintf("Power consumption cooling: %s kWh\nPower consumption heating: %s kWh",
		c.DayCool.String(), c.DayHeat.String())
}


func (d *Daikin) parseResponse(resp *http.Response) (map[string]string, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(strings.NewReader(string(body)))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) != 1 {
		return nil, fmt.Errorf("Have %d rows of records, want just one", len(records))
	}

	values := map[string]string{}
	for _, rec := range records[0] {
		parts := strings.SplitN(rec, "=", 2)
		values[parts[0]] = parts[1]
	}
	return values, nil

}

// GetBasicInfo gets the basic information for the unit.
func (d *Daikin) GetBasicInfo() error {
	resp, err := http.Get(fmt.Sprintf("http://%s%s", d.Address, uriGetBasicInfo))
	if err != nil {
		return err
	}
	d.BasicInfo = &BasicInfo{}
	vals, err := d.parseResponse(resp)
	if err != nil {
		return err
	}
	return d.BasicInfo.populate(vals)
}

// Set configures the current setting to the unit.
func (d *Daikin) SetControlInfo() error {
      	resp, err := http.Get(fmt.Sprintf("http://%s%s?%s", d.Address,
	                      uriSetControlInfo, d.ControlInfo.urlValues()))
	if err != nil {
		return err
	}
	vals, err := d.parseResponse(resp)
	if err != nil {
		return err
	}
	if v := vals["ret"]; v != "OK" {
		return fmt.Errorf("device returned error ret=%s", v)
	}
	return nil
}

// GetControlInfo gets the current control settings for the unit.
func (d *Daikin) GetControlInfo() error {
	resp, err := http.Get(fmt.Sprintf("http://%s%s", d.Address, uriGetControlInfo))
	if err != nil {
		return err
	}
	d.ControlInfo = &ControlInfo{}
	vals, err := d.parseResponse(resp)
	if err != nil {
		return err
	}
	return d.ControlInfo.populate(vals)
}

// GetSensorInfo gets the current sensor values for the unit.
func (d *Daikin) GetSensorInfo() error {
	resp, err := http.Get(fmt.Sprintf("http://%s%s", d.Address, uriGetSensorInfo))
	if err != nil {
		return err
	}
	d.SensorInfo = &SensorInfo{}
	vals, err := d.parseResponse(resp)
	if err != nil {
		return err
	}
	return d.SensorInfo.populate(vals)
}

// GetPowerInfo gets the current power consumption for the unit.
func (d *Daikin) GetPowerInfo() error {
	resp, err := http.Get(fmt.Sprintf("http://%s%s", d.Address, uriGetDayPowerEx))
	if err != nil {
		return err
	}
	d.PowerInfo = &PowerInfo{}
	vals, err := d.parseResponse(resp)
	if err != nil {
		return err
	}
	return d.PowerInfo.populate(vals)
}

func (d *Daikin) String() string {
	var ret string
	if d.BasicInfo != nil {
		ret = ret + d.BasicInfo.String() + "\n"
	}
	if d.ControlInfo != nil {
		ret = ret + d.ControlInfo.String() + "\n"
	}
	if d.SensorInfo != nil {
		ret = ret + d.SensorInfo.String() + "\n"
	}
	if d.PowerInfo != nil {
		ret = ret + d.PowerInfo.String() + "\n"
	}
	return ret
}
