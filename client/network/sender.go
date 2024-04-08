package network

import (
	"bytes"
	"encoding/gob"
)

type NetworkData map[string]interface{}

func Encode(data NetworkData) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func Decode(data []byte) (NetworkData, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	var result NetworkData
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}
