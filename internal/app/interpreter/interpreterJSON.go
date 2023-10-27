package interpreter

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
)

type InterpreterJSON struct{}

func (i *InterpreterJSON) Decode(r io.Reader, v interface{}) error {

	var buf bytes.Buffer
	// читаем тело запроса
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}

	log.Println(buf.String())

	// десериализуем JSON в Visitor
	return json.Unmarshal(buf.Bytes(), &v)
}
