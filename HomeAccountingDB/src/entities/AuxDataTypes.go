package entities

import (
	"encoding/json"
)

type Date int

func (n *Date) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	var parts []int
	err := json.Unmarshal(b, &parts)
	if err != nil {
		return err
	}
	*n = Date(parts[0]*10000 + parts[1]*100 + parts[2])
	return nil
}

func toInt(b []byte) int {
	result := 0
	minus := false
	for _, c := range b {
		if c >= '0' && c <= '9' {
			result *= 10
			result += int(c - '0')
		} else if c == '-' {
			minus = true
		}
	}
	if minus {
		return -result
	}
	return result
}

type Decimal int

func (n *Decimal) UnmarshalJSON(b []byte) (err error) {
	*n = Decimal(toInt(b))
	return nil
}
