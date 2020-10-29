package decoder

import (
	"encoding/json"
	"strconv"
)

type Float64 struct {
	Value *float64
}

func (d *Float64) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}

	*d.Value, err = strconv.ParseFloat(s, 64)
	return
}
