package asdf

import (
	"bytes"
	"encoding/gob"
)

func BinUnmarshal(obj ISlice, data []byte) error {
	copy(obj.Slice(), data)

	return nil
}

func BinMarshal(obj ISlice) ([]byte, error) {
	bin := obj.Slice()

	buf := make([]byte, len(bin))
	copy(buf, bin)

	return buf, nil
}

func GobMarshal(obj interface{}) ([]byte, error) {
	var w bytes.Buffer

	enc := gob.NewEncoder(&w)

	err := enc.Encode(obj)
	if nil != err {
		return nil, err
	}

	return w.Bytes(), nil
}

func GobUnmarshal(data []byte, obj interface{}) error {
	r := bytes.NewBuffer(data)
	dec := gob.NewDecoder(r)

	err := dec.Decode(obj)
	if nil != err {
		return err
	}

	return nil
}
