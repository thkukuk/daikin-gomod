package daikin

import (
        "fmt"
        "strconv"
)

// Humidity
type Humidity struct {
        value int32
        param string
}

func (h *Humidity) setUrlValues() string {
	return h.param + "=" + h.String()
}

func (h *Humidity) decode(param string, v string) error {
	if v == "--" || v == "-" {
		v = "-1"
	}
	val, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("Humidity: error parsing %s=%s: %v", param, v, err)
	}
	*h = Humidity{value: int32(val), param: param}
	return nil
}

func (h *Humidity) String() string {
        if h.value == -1 {
                return "N/A"
        }
        return strconv.Itoa(int(h.value))
}

func (h *Humidity) Float64() float64 {
	return float64(h.value)
}
