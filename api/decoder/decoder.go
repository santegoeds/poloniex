package decoder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/santegoeds/poloniex/errors"
)

func Unmarshal(msgs []json.RawMessage, receivers ...interface{}) error {
	if len(msgs) != len(receivers) {
		panic(
			fmt.Sprintf(
				"number of messages %d does not match number of receivers %d",
				len(msgs), len(receivers),
			),
		)
	}
	for idx := 0; idx < len(msgs); idx++ {
		if err := json.Unmarshal(msgs[idx], receivers[idx]); err != nil {
			return err
		}
	}
	return nil
}

func DecodeObject(r io.Reader) (map[string]json.RawMessage, error) {
	objData := make(map[string]json.RawMessage)
	dec := json.NewDecoder(r)
	if err := dec.Decode(&objData); err != nil {
		return nil, err
	}
	errData, ok := objData["error"]
	if !ok {
		return objData, nil
	}

	var errMsg string
	if err := json.Unmarshal(errData, &errMsg); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("%s: %w", errMsg, errors.ErrBadRequest)
}

func DecodeMessage(data []byte, out interface{}) error {
	if err := json.Unmarshal(data, out); err != nil {
		if _, rspErr := DecodeObject(bytes.NewBuffer(data)); rspErr != nil {
			return rspErr
		}
		return err
	}
	return nil
}
