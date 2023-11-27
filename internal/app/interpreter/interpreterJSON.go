package interpreter

import (
	"bytes"
	"encoding/json"
	"io"
)

type InterpreterJSON struct{}

func (i *InterpreterJSON) Unmarshal(r io.Reader, v interface{}) error {

	var buf bytes.Buffer
	// читаем тело запроса
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	// десериализуем JSON в Visitor
	return json.Unmarshal(buf.Bytes(), &v)
}

func (i *InterpreterJSON) Decode(r io.Reader, v interface{}) error {

	return json.NewDecoder(r).Decode(v)
}
