package daikin

import (
        "fmt"
	"math"
        "strconv"
)

// Fan is the fan speed of the Daikin unit.
type Fan string

// Fan values. Not all may be valid on all models.
const (
	FanAuto   Fan = "A"
	FanSilent Fan = "B"
	Fan1      Fan = "3"
	Fan2      Fan = "4"
	Fan3      Fan = "5"
	Fan4      Fan = "6"
	Fan5      Fan = "7"
)

var fanMap = map[Fan]string{
	FanAuto:   "Auto",
	FanSilent: "Silent",
	Fan1:      "1",
	Fan2:      "2",
	Fan3:      "3",
	Fan4:      "4",
	Fan5:      "5",
}

func (f *Fan) setUrlValues() string {
	return "f_rate=" + string(*f)
}

func (f *Fan) decode(s string) error {
	// check if value is supported
	if _, ok := fanMap[Fan(s)]; !ok {
		return fmt.Errorf("unknown pwr value: %s", s)
	}
	*f = Fan(s)

	return nil
}

func (f *Fan) String() string {
	v, ok := fanMap[*f]
	if !ok {
		return fmt.Sprintf("Unknown Fan [%v]", *f)
	}
	return v
}

func (f *Fan) Float64() float64 {
	switch *f {
	case FanAuto:
		return -1.0
	case FanSilent:
		return 0.0
	default:
		val, err := strconv.ParseFloat(f.String(), 64)
		if err != nil {
			return math.NaN()
		}
		return val
	}

}

