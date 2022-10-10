package message

import (
	"bytes"
	"encoding/json"
)

type Encoder interface {
	Encode(source interface{}) ([]byte, error)
}

type JsonEncoder struct {
}

func (e *JsonEncoder) Encode(source interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(source); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type Decoder interface {
	Decode(source []byte, target interface{}) error
}

type JsonDecoder struct {
}

func (d *JsonDecoder) Decode(source []byte, target interface{}) error {
	return json.NewDecoder(bytes.NewReader(source)).Decode(target)
}
