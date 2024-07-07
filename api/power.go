package daikin

import (
	"fmt"
	"strconv"
)

// Power represents the power status of the unit (off/on).
type Power int

// The power status of the unit.
const (
	PowerOff Power = 0
	PowerOn  Power = 1
)

var powerMap = map[Power]string{
	PowerOff: "Off",
	PowerOn:  "On",
}

func (p *Power) setUrlValues() string {
	return "pow=" + strconv.Itoa(int(*p))
}

func (p *Power) decode(s string) error {
	switch s {
	case "0":
		*p = Power(PowerOff)
	case "1":
		*p = Power(PowerOn)
	default:
		return fmt.Errorf("unknown pwr value: %s", s)
	}
	return nil
}

func (p *Power) String() string {
	v, ok := powerMap[*p]
	if !ok {
		return fmt.Sprintf("Unknown Power [%d]", int(*p))
	}
	return v
}

func (p *Power) Float64() float64 {
	return float64(*p)
}

