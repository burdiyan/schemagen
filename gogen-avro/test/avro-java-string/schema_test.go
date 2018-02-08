package avro

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"gopkg.in/linkedin/goavro.v1"
	"testing"
)

/* Round-trip some primitive values through our serializer and goavro to verify */
var fixtures = []Event{
	{
		Id: "id1",
	},
	{
		Id: "differentid",
	},
}

func compareFixtureGoAvro(t *testing.T, actual interface{}, expected interface{}) {
	record := actual.(*goavro.Record)
	fixture := expected.(Event)
	id, err := record.Get("id")
	assert.Nil(t, err)
	assert.Equal(t, id, fixture.Id)
}

func TestRootUnionFixture(t *testing.T) {
	codec, err := goavro.NewCodec(fixtures[0].Schema())
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = writeEvent(&f, &buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, err := codec.Decode(&buf)
		if err != nil {
			t.Fatal(err)
		}
		compareFixtureGoAvro(t, datum, f)
	}
}

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err := writeEvent(&f, &buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, err := readEvent(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, datum, &f)
	}
}
