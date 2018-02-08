package avro

import (
	"bytes"
	"encoding/json"
	"github.com/alanctgardner/gogen-avro/container"
	"gopkg.in/linkedin/goavro.v1"
	"testing"
)

/* Round-trip some primitive values through our container file writer and goavro to verify */

func TestNullEncoding(t *testing.T) {
	roundTripWithCodec(container.Null, t)
}

func TestSnappyEncoding(t *testing.T) {
	roundTripWithCodec(container.Deflate, t)
}

func TestDeflateEncoding(t *testing.T) {
	roundTripWithCodec(container.Snappy, t)
}

func roundTripWithCodec(codec container.Codec, t *testing.T) {
	fixtures := make([]PrimitiveTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	// Write the container file contents to the buffer
	var containerWriter *container.Writer
	containerWriter, err = NewPrimitiveTestRecordWriter(&buf, codec, 2)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fixtures {
		// Write the record to the container file
		err = containerWriter.WriteRecord(&f)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Flush the buffers to ensure the last block has been written
	err = containerWriter.Flush()
	if err != nil {
		t.Fatal(err)
	}

	reader, err := goavro.NewReader(goavro.FromReader(&buf))
	if err != nil {
		t.Fatal(err)
	}

	var i int
	for reader.Scan() {
		datum, err := reader.Read()
		if err != nil {
			t.Fatal(err)
		}
		compareFixtureGoAvro(t, datum, fixtures[i])
		i = i + 1
	}
}
