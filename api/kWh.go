package daikin

import (
	"fmt"
	"strconv"
)

type KWattHours struct {
	value float64
	param string
}

func (w *KWattHours) setUrlValues() string {
	return ""
}

func (w *KWattHours) decode(param string, v string) error {
	if v == "-" {
		v = "-1"
	}
	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("error parsing watt hours=%s: %v", v, err)
	}
	*w = KWattHours{value: val, param: param}
	return nil
}

func (w *KWattHours) String() string {
	return strconv.FormatFloat(w.value, 'f', 1, 64)
}

func (w *KWattHours) Float64() float64 {
	return w.value
}
