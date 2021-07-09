package decoder

import (
	"encoding/json"
	"strconv"
	"time"
)

type EpochMs struct {
	Value *time.Time
}

func (d *EpochMs) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	epochMs, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*d.Value = time.Unix(0, epochMs*1000)
	return nil
}
