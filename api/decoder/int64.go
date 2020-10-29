package decoder

import (
	"encoding/json"
	"strconv"
)

type Int64 struct {
	Value *int64
}

func (d *Int64) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	*d.Value, err = strconv.ParseInt(s, 10, 64)
	return
}
