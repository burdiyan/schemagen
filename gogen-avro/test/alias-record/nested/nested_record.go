// Code generated by github.com/alanctgardner/gogen-avro. DO NOT EDIT.
/*
 * SOURCE:
 *     nested.avsc
 */
package avro

import (
	"io"
)

type NestedRecord struct {
	StringField string
	BoolField   bool
	BytesField  []byte
}

func DeserializeNestedRecord(r io.Reader) (*NestedRecord, error) {
	return readNestedRecord(r)
}

func (r *NestedRecord) Schema() string {
	return "{\"fields\":[{\"name\":\"StringField\",\"type\":\"string\"},{\"name\":\"BoolField\",\"type\":\"boolean\"},{\"name\":\"BytesField\",\"type\":\"bytes\"}],\"name\":\"NestedRecord\",\"type\":\"record\"}"
}

func (r *NestedRecord) Serialize(w io.Writer) error {
	return writeNestedRecord(r, w)
}